package ec2

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceVPCPeeringConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPCPeeringCreate,
		Read:   resourceVPCPeeringRead,
		Update: resourceVPCPeeringUpdate,
		Delete: resourceVPCPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"peer_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auto_accept": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"accept_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"accepter":  vpcPeeringConnectionOptionsSchema(),
			"requester": vpcPeeringConnectionOptionsSchema(),
			"tags":      tftags.TagsSchema(),
			"tags_all":  tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceVPCPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	// Create the vpc peering connection
	createOpts := &ec2.CreateVpcPeeringConnectionInput{
		PeerVpcId:         aws.String(d.Get("peer_vpc_id").(string)),
		VpcId:             aws.String(d.Get("vpc_id").(string)),
		TagSpecifications: ec2TagSpecificationsFromKeyValueTags(tags, ec2.ResourceTypeVpcPeeringConnection),
	}

	if v, ok := d.GetOk("peer_owner_id"); ok {
		createOpts.PeerOwnerId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("peer_region"); ok {
		if _, ok := d.GetOk("auto_accept"); ok {
			return fmt.Errorf("peer_region cannot be set whilst auto_accept is true when creating a vpc peering connection")
		}
		createOpts.PeerRegion = aws.String(v.(string))
	}

	log.Printf("[DEBUG] VPC Peering Create options: %#v", createOpts)

	resp, err := conn.CreateVpcPeeringConnection(createOpts)
	if err != nil {
		return fmt.Errorf("Error creating VPC Peering Connection: %s", err)
	}

	// Get the ID and store it
	rt := resp.VpcPeeringConnection
	d.SetId(aws.StringValue(rt.VpcPeeringConnectionId))
	log.Printf("[INFO] VPC Peering Connection ID: %s", d.Id())

	err = vpcPeeringConnectionWaitUntilAvailable(conn, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error waiting for VPC Peering Connection to become available: %s", err)
	}

	return resourceVPCPeeringUpdate(d, meta)
}

func resourceVPCPeeringRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*conns.AWSClient)
	defaultTagsConfig := client.DefaultTagsConfig
	ignoreTagsConfig := client.IgnoreTagsConfig

	pcRaw, statusCode, err := vpcPeeringConnectionRefreshState(client.EC2Conn, d.Id())()
	// Allow a failed VPC Peering Connection to fallthrough,
	// to allow rest of the logic below to do its work.
	if err != nil && statusCode != ec2.VpcPeeringConnectionStateReasonCodeFailed {
		return fmt.Errorf("Error reading VPC Peering Connection: %s", err)
	}

	// The failed status is a status that we can assume just means the
	// connection is gone. Destruction isn't allowed, and it eventually
	// just "falls off" the console. See GH-2322
	status := map[string]bool{
		ec2.VpcPeeringConnectionStateReasonCodeDeleted:  true,
		ec2.VpcPeeringConnectionStateReasonCodeDeleting: true,
		ec2.VpcPeeringConnectionStateReasonCodeExpired:  true,
		ec2.VpcPeeringConnectionStateReasonCodeFailed:   true,
		ec2.VpcPeeringConnectionStateReasonCodeRejected: true,
		"": true, // AWS consistency issue, see vpcPeeringConnectionRefreshState
	}
	if _, ok := status[statusCode]; ok {
		log.Printf("[WARN] VPC Peering Connection (%s) has status code %s, removing from state", d.Id(), statusCode)
		d.SetId("")
		return nil
	}

	pc := pcRaw.(*ec2.VpcPeeringConnection)
	log.Printf("[DEBUG] VPC Peering Connection response: %#v", pc)

	log.Printf("[DEBUG] Account ID %s, VPC PeerConn Requester %s, Accepter %s",
		client.AccountID, *pc.RequesterVpcInfo.OwnerId, *pc.AccepterVpcInfo.OwnerId)

	if (client.AccountID == aws.StringValue(pc.AccepterVpcInfo.OwnerId)) && (client.AccountID != aws.StringValue(pc.RequesterVpcInfo.OwnerId)) {
		// We're the accepter
		d.Set("peer_owner_id", pc.RequesterVpcInfo.OwnerId)
		d.Set("peer_vpc_id", pc.RequesterVpcInfo.VpcId)
		d.Set("vpc_id", pc.AccepterVpcInfo.VpcId)
	} else {
		// We're the requester
		d.Set("peer_owner_id", pc.AccepterVpcInfo.OwnerId)
		d.Set("peer_vpc_id", pc.AccepterVpcInfo.VpcId)
		d.Set("vpc_id", pc.RequesterVpcInfo.VpcId)
	}

	d.Set("peer_region", pc.AccepterVpcInfo.Region)
	d.Set("accept_status", pc.Status.Code)

	if err := d.Set("accepter", flattenVPCPeeringConnectionOptions(pc.AccepterVpcInfo.PeeringOptions)); err != nil {
		return fmt.Errorf("Error setting VPC Peering Connection accepter information: %s", err)
	}
	if err := d.Set("requester", flattenVPCPeeringConnectionOptions(pc.RequesterVpcInfo.PeeringOptions)); err != nil {
		return fmt.Errorf("Error setting VPC Peering Connection requester information: %s", err)
	}

	tags := KeyValueTags(pc.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}
	if err != nil {
		return fmt.Errorf("Error setting VPC Peering Connection tags: %s", err)
	}

	return nil
}

func resourceVPCPeeringConnectionAccept(conn *ec2.EC2, id string) (string, error) {
	log.Printf("[INFO] Accept VPC Peering Connection with ID: %s", id)

	req := &ec2.AcceptVpcPeeringConnectionInput{
		VpcPeeringConnectionId: aws.String(id),
	}

	resp, err := conn.AcceptVpcPeeringConnection(req)
	if err != nil {
		return "", err
	}

	return aws.StringValue(resp.VpcPeeringConnection.Status.Code), nil
}

func resourceVPCPeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	if d.HasChange("tags_all") && !d.IsNewResource() {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating EC2 VPC Peering Connection (%s) tags: %s", d.Id(), err)
		}
	}

	pcRaw, statusCode, err := vpcPeeringConnectionRefreshState(conn, d.Id())()
	if err != nil {
		return fmt.Errorf("Error reading VPC Peering Connection: %s", err)
	}

	if pcRaw == nil {
		log.Printf("[WARN] VPC Peering Connection (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if _, ok := d.GetOk("auto_accept"); ok && statusCode == ec2.VpcPeeringConnectionStateReasonCodePendingAcceptance {
		statusCode, err = resourceVPCPeeringConnectionAccept(conn, d.Id())
		if err != nil {
			return fmt.Errorf("Unable to accept VPC Peering Connection: %s", err)
		}
		log.Printf("[DEBUG] VPC Peering Connection accept status: %s", statusCode)

		// "OperationNotPermitted: Peering pcx-0000000000000000 is not active. Peering options can be added only to active peerings."
		if err := vpcPeeringConnectionWaitUntilAvailable(conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("Error waiting for VPC Peering Connection to become available: %s", err)
		}
	}

	if d.HasChanges("accepter", "requester") {
		if statusCode == ec2.VpcPeeringConnectionStateReasonCodeActive || statusCode == ec2.VpcPeeringConnectionStateReasonCodeProvisioning {
			pc := pcRaw.(*ec2.VpcPeeringConnection)
			crossRegionPeering := aws.StringValue(pc.RequesterVpcInfo.Region) != aws.StringValue(pc.AccepterVpcInfo.Region)

			req := &ec2.ModifyVpcPeeringConnectionOptionsInput{
				VpcPeeringConnectionId: aws.String(d.Id()),
			}
			if d.HasChange("accepter") {
				req.AccepterPeeringConnectionOptions = expandVPCPeeringConnectionOptions(d.Get("accepter").([]interface{}), crossRegionPeering)
			}
			if d.HasChange("requester") {
				req.RequesterPeeringConnectionOptions = expandVPCPeeringConnectionOptions(d.Get("requester").([]interface{}), crossRegionPeering)
			}

			log.Printf("[DEBUG] Modifying VPC Peering Connection options: %s", req)
			if _, err := conn.ModifyVpcPeeringConnectionOptions(req); err != nil {
				return fmt.Errorf("error modifying VPC Peering Connection (%s) Options: %s", d.Id(), err)
			}
		} else {
			return fmt.Errorf("Unable to modify peering options. The VPC Peering Connection "+
				"%q is not active. Please set `auto_accept` attribute to `true`, "+
				"or activate VPC Peering Connection manually.", d.Id())
		}
	}

	return resourceVPCPeeringRead(d, meta)
}

func resourceVPCPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	req := &ec2.DeleteVpcPeeringConnectionInput{
		VpcPeeringConnectionId: aws.String(d.Id()),
	}

	_, err := conn.DeleteVpcPeeringConnection(req)

	if tfawserr.ErrMessageContains(err, "InvalidVpcPeeringConnectionID.NotFound", "") {
		return nil
	}

	// "InvalidStateTransition: Invalid state transition for pcx-0000000000000000, attempted to transition from failed to deleting"
	if tfawserr.ErrMessageContains(err, "InvalidStateTransition", "to deleting") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error deleting VPC Peering Connection (%s): %s", d.Id(), err)
	}

	if err := WaitForVPCPeeringConnectionDeletion(conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("Error waiting for VPC Peering Connection (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

// vpcPeeringConnection returns the VPC peering connection corresponding to the specified identifier.
// Returns nil if no VPC peering connection is found or the connection has reached a terminal state
// according to https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-lifecycle.
func vpcPeeringConnection(conn *ec2.EC2, vpcPeeringConnectionID string) (*ec2.VpcPeeringConnection, error) {
	outputRaw, _, err := StatusVPCPeeringConnection(conn, vpcPeeringConnectionID)()

	if output, ok := outputRaw.(*ec2.VpcPeeringConnection); ok {
		return output, err
	}

	return nil, err
}

func vpcPeeringConnectionRefreshState(conn *ec2.EC2, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.DescribeVpcPeeringConnections(&ec2.DescribeVpcPeeringConnectionsInput{
			VpcPeeringConnectionIds: aws.StringSlice([]string{id}),
		})
		if err != nil {
			if tfawserr.ErrMessageContains(err, "InvalidVpcPeeringConnectionID.NotFound", "") {
				return nil, ec2.VpcPeeringConnectionStateReasonCodeDeleted, nil
			}

			return nil, "", err
		}

		if resp == nil || resp.VpcPeeringConnections == nil ||
			len(resp.VpcPeeringConnections) == 0 || resp.VpcPeeringConnections[0] == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// our peering connection yet. Return an empty state.
			return nil, "", nil
		}
		pc := resp.VpcPeeringConnections[0]
		if pc.Status == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// our peering connection yet. Return an empty state.
			return nil, "", nil
		}
		statusCode := aws.StringValue(pc.Status.Code)

		// A VPC Peering Connection can exist in a failed state due to
		// incorrect VPC ID, account ID, or overlapping IP address range,
		// thus we short circuit before the time out would occur.
		if statusCode == ec2.VpcPeeringConnectionStateReasonCodeFailed {
			return nil, statusCode, errors.New(aws.StringValue(pc.Status.Message))
		}

		return pc, statusCode, nil
	}
}

func vpcPeeringConnectionOptionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allow_remote_vpc_dns_resolution": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"allow_classic_link_to_remote_vpc": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"allow_vpc_to_remote_classic_link": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func vpcPeeringConnectionWaitUntilAvailable(conn *ec2.EC2, id string, timeout time.Duration) error {
	// Wait for the vpc peering connection to become available
	log.Printf("[DEBUG] Waiting for VPC Peering Connection (%s) to become available.", id)
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			ec2.VpcPeeringConnectionStateReasonCodeInitiatingRequest,
			ec2.VpcPeeringConnectionStateReasonCodeProvisioning,
		},
		Target: []string{
			ec2.VpcPeeringConnectionStateReasonCodePendingAcceptance,
			ec2.VpcPeeringConnectionStateReasonCodeActive,
		},
		Refresh: vpcPeeringConnectionRefreshState(conn, id),
		Timeout: timeout,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for VPC Peering Connection (%s) to become available: %s", id, err)
	}
	return nil
}

func WaitForVPCPeeringConnectionDeletion(conn *ec2.EC2, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			ec2.VpcPeeringConnectionStateReasonCodeActive,
			ec2.VpcPeeringConnectionStateReasonCodePendingAcceptance,
			ec2.VpcPeeringConnectionStateReasonCodeDeleting,
		},
		Target: []string{
			ec2.VpcPeeringConnectionStateReasonCodeRejected,
			ec2.VpcPeeringConnectionStateReasonCodeDeleted,
		},
		Refresh: vpcPeeringConnectionRefreshState(conn, id),
		Timeout: timeout,
	}

	_, err := stateConf.WaitForState()

	return err
}
