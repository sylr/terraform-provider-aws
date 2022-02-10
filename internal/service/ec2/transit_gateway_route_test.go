package ec2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
)

func testAccTransitGatewayRoute_basic(t *testing.T) {
	var transitGatewayRoute1 ec2.TransitGatewayRoute
	resourceName := "aws_ec2_transit_gateway_route.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	transitGatewayVpcAttachmentResourceName := "aws_ec2_transit_gateway_vpc_attachment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTransitGatewayRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayRouteDestinationCIDRBlockConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayRouteExists(resourceName, &transitGatewayRoute1),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "blackhole", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_attachment_id", transitGatewayVpcAttachmentResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_route_table_id", transitGatewayResourceName, "association_default_route_table_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTransitGatewayRoute_basic_ipv6(t *testing.T) {
	var transitGatewayRoute1 ec2.TransitGatewayRoute
	resourceName := "aws_ec2_transit_gateway_route.test_ipv6"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	transitGatewayVpcAttachmentResourceName := "aws_ec2_transit_gateway_vpc_attachment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTransitGatewayRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayRouteDestinationCIDRBlockConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayRouteExists(resourceName, &transitGatewayRoute1),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", "2001:db8::/56"),
					resource.TestCheckResourceAttr(resourceName, "blackhole", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_attachment_id", transitGatewayVpcAttachmentResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_route_table_id", transitGatewayResourceName, "association_default_route_table_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTransitGatewayRoute_blackhole(t *testing.T) {
	var transitGatewayRoute1 ec2.TransitGatewayRoute
	resourceName := "aws_ec2_transit_gateway_route.test_blackhole"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTransitGatewayRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayRouteDestinationCIDRBlockConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayRouteExists(resourceName, &transitGatewayRoute1),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", "10.1.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "blackhole", "true"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_attachment_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_route_table_id", transitGatewayResourceName, "association_default_route_table_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTransitGatewayRoute_disappears(t *testing.T) {
	var transitGateway1 ec2.TransitGateway
	var transitGatewayRoute1 ec2.TransitGatewayRoute
	resourceName := "aws_ec2_transit_gateway_route.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTransitGatewayRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayRouteDestinationCIDRBlockConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(transitGatewayResourceName, &transitGateway1),
					testAccCheckTransitGatewayRouteExists(resourceName, &transitGatewayRoute1),
					testAccCheckTransitGatewayRouteDisappears(&transitGateway1, &transitGatewayRoute1),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccTransitGatewayRoute_disappears_TransitGatewayAttachment(t *testing.T) {
	var transitGateway1 ec2.TransitGateway
	var transitGatewayRoute1 ec2.TransitGatewayRoute
	var transitGatewayVpcAttachment1 ec2.TransitGatewayVpcAttachment
	resourceName := "aws_ec2_transit_gateway_route.test"
	transitGatewayVpcAttachmentResourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckTransitGateway(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTransitGatewayRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayRouteDestinationCIDRBlockConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(transitGatewayResourceName, &transitGateway1),
					testAccCheckTransitGatewayRouteExists(resourceName, &transitGatewayRoute1),
					testAccCheckTransitGatewayVPCAttachmentExists(transitGatewayVpcAttachmentResourceName, &transitGatewayVpcAttachment1),
					testAccCheckTransitGatewayVPCAttachmentDisappears(&transitGatewayVpcAttachment1),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckTransitGatewayRouteExists(resourceName string, transitGatewayRoute *ec2.TransitGatewayRoute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 Transit Gateway Route ID is set")
		}

		transitGatewayRouteTableID, destination, err := tfec2.DecodeTransitGatewayRouteID(rs.Primary.ID)

		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		route, err := tfec2.DescribeTransitGatewayRoute(conn, transitGatewayRouteTableID, destination)

		if err != nil {
			return err
		}

		if route == nil {
			return fmt.Errorf("EC2 Transit Gateway Route not found")
		}

		*transitGatewayRoute = *route

		return nil
	}
}

func testAccCheckTransitGatewayRouteDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ec2_transit_gateway_route" {
			continue
		}

		transitGatewayRouteTableID, destination, err := tfec2.DecodeTransitGatewayRouteID(rs.Primary.ID)

		if err != nil {
			return err
		}

		route, err := tfec2.DescribeTransitGatewayRoute(conn, transitGatewayRouteTableID, destination)

		if tfawserr.ErrMessageContains(err, "InvalidRouteTableID.NotFound", "") {
			continue
		}

		if err != nil {
			return err
		}

		if route == nil {
			continue
		}

		return fmt.Errorf("EC2 Transit Gateway Route (%s) still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckTransitGatewayRouteDisappears(transitGateway *ec2.TransitGateway, transitGatewayRoute *ec2.TransitGatewayRoute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		input := &ec2.DeleteTransitGatewayRouteInput{
			DestinationCidrBlock:       transitGatewayRoute.DestinationCidrBlock,
			TransitGatewayRouteTableId: transitGateway.Options.AssociationDefaultRouteTableId,
		}

		_, err := conn.DeleteTransitGatewayRoute(input)

		return err
	}
}

func testAccTransitGatewayRouteDestinationCIDRBlockConfig() string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), `
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "tf-acc-test-ec2-transit-gateway-route"
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = "tf-acc-test-ec2-transit-gateway-route"
  }
}

resource "aws_ec2_transit_gateway" "test" {}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id
}

resource "aws_ec2_transit_gateway_route" "test" {
  destination_cidr_block         = "0.0.0.0/0"
  transit_gateway_attachment_id  = aws_ec2_transit_gateway_vpc_attachment.test.id
  transit_gateway_route_table_id = aws_ec2_transit_gateway.test.association_default_route_table_id
}

resource "aws_ec2_transit_gateway_route" "test_ipv6" {
  destination_cidr_block         = "2001:db8::/56"
  transit_gateway_attachment_id  = aws_ec2_transit_gateway_vpc_attachment.test.id
  transit_gateway_route_table_id = aws_ec2_transit_gateway.test.association_default_route_table_id
}

resource "aws_ec2_transit_gateway_route" "test_blackhole" {
  destination_cidr_block         = "10.1.0.0/16"
  blackhole                      = true
  transit_gateway_route_table_id = aws_ec2_transit_gateway.test.association_default_route_table_id
}
`)
}
