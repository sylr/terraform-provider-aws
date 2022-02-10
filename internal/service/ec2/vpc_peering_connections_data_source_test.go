package ec2_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccEC2VPCPeeringConnectionsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCPeeringConnectionsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_vpc_peering_connections.test_by_filters", "ids.#", "2"),
				),
			},
		},
	})
}

func TestAccEC2VPCPeeringConnectionsDataSource_NoMatches(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCPeeringConnectionsDataSourceConfig_NoMatches,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_vpc_peering_connections.test", "ids.#", "0"),
				),
			},
		},
	})
}

const testAccVPCPeeringConnectionsDataSourceConfig = `
resource "aws_vpc" "foo" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = "terraform-testacc-vpc-peering-connection-data-source-foo"
    Type = "primary"
  }
}

resource "aws_vpc" "bar" {
  cidr_block = "10.2.0.0/16"

  tags = {
    Name = "terraform-testacc-vpc-peering-connection-data-source-bar"
    Type = "secondary"
  }
}

resource "aws_vpc" "baz" {
  cidr_block = "10.3.0.0/16"

  tags = {
    Name = "terraform-testacc-vpc-peering-connection-data-source-baz"
    Type = "secondary"
  }
}

resource "aws_vpc_peering_connection" "conn1" {
  vpc_id      = aws_vpc.foo.id
  peer_vpc_id = aws_vpc.bar.id
  auto_accept = true

  tags = {
    Name        = "terraform-testacc-vpc-peering-connection-data-source-foo-to-bar"
    Environment = "test"
  }
}

resource "aws_vpc_peering_connection" "conn2" {
  vpc_id      = aws_vpc.foo.id
  peer_vpc_id = aws_vpc.baz.id
  auto_accept = true

  tags = {
    Name        = "terraform-testacc-vpc-peering-connection-data-source-foo-to-baz"
    Environment = "test"
  }
}

data "aws_vpc_peering_connections" "test_by_filters" {
  filter {
    name   = "vpc-peering-connection-id"
    values = [aws_vpc_peering_connection.conn1.id, aws_vpc_peering_connection.conn2.id]
  }
}
`

const testAccVPCPeeringConnectionsDataSourceConfig_NoMatches = `
data "aws_vpc_peering_connections" "test" {
  tags = {
    Name = "Non-Existent"
  }
}
`
