package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceVPCPeeringConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPCPeeringConnectionsRead,

		Schema: map[string]*schema.Schema{
			"filter": CustomFiltersSchema(),
			"tags":   tftags.TagsSchemaComputed(),
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceVPCPeeringConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	req := &ec2.DescribeVpcPeeringConnectionsInput{}

	req.Filters = append(req.Filters, BuildTagFilterList(
		Tags(tftags.New(d.Get("tags").(map[string]interface{}))),
	)...)
	req.Filters = append(req.Filters, BuildCustomFilterList(
		d.Get("filter").(*schema.Set),
	)...)
	if len(req.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		req.Filters = nil
	}

	resp, err := conn.DescribeVpcPeeringConnections(req)
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("error reading EC2 VPC Peering Connections: empty response")
	}

	var ids []string
	for _, pcx := range resp.VpcPeeringConnections {
		ids = append(ids, aws.StringValue(pcx.VpcPeeringConnectionId))
	}

	d.SetId(meta.(*conns.AWSClient).Region)

	err = d.Set("ids", ids)
	if err != nil {
		return err
	}

	return nil
}
