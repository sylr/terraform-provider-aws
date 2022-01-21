package ec2_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccEC2DefaultVPC_basic(t *testing.T) {
	var vpc ec2.Vpc

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDefaultVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDefaultVPCBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckVPCExists("aws_default_vpc.foo", &vpc),
					resource.TestCheckResourceAttr("aws_default_vpc.foo", "cidr_block", "172.31.0.0/16"),
					resource.TestCheckResourceAttr(
						"aws_default_vpc.foo", "cidr_block", "172.31.0.0/16"),
					resource.TestCheckResourceAttr(
						"aws_default_vpc.foo", "tags.%", "1"),
					resource.TestCheckResourceAttr(
						"aws_default_vpc.foo", "tags.Name", "Default VPC"),
					resource.TestCheckResourceAttrSet(
						"aws_default_vpc.foo", "arn"),
					resource.TestCheckResourceAttr(
						"aws_default_vpc.foo", "assign_generated_ipv6_cidr_block", "false"),
					resource.TestCheckResourceAttr(
						"aws_default_vpc.foo", "ipv6_association_id", ""),
					resource.TestCheckResourceAttr(
						"aws_default_vpc.foo", "ipv6_cidr_block", ""),
					acctest.CheckResourceAttrAccountID("aws_default_vpc.foo", "owner_id"),
				),
			},
		},
	})
}

func testAccCheckDefaultVPCDestroy(s *terraform.State) error {
	// We expect VPC to still exist
	return nil
}

const testAccDefaultVPCBasicConfig = `
resource "aws_default_vpc" "foo" {
  tags = {
    Name = "Default VPC"
  }
}
`
