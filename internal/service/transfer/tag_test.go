package transfer_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/transfer"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	tftransfer "github.com/hashicorp/terraform-provider-aws/internal/service/transfer"
)

func TestAccTransferTag_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_transfer_tag.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, transfer.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagConfig_basic(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "value", "value1"),
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

func TestAccTransferTag_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_transfer_tag.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, transfer.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagConfig_basic(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tftransfer.ResourceTag(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTransferTag_value(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_transfer_tag.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, transfer.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagConfig_basic(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "value", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTagConfig_basic(rName, "key1", "value1updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "value", "value1updated"),
				),
			},
		},
	})
}

func TestAccTransferTag_system(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_transfer_tag.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, transfer.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagConfig_basic(rName, "aws:transfer:customHostname", "abc.example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key", "aws:transfer:customHostname"),
					resource.TestCheckResourceAttr(resourceName, "value", "abc.example.com"),
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

func testAccTagConfig_basic(rName string, key string, value string) string {
	return fmt.Sprintf(`
resource "aws_transfer_server" "test" {
  identity_provider_type = "SERVICE_MANAGED"

  tags = {
    Name = %[1]q
  }

  lifecycle {
    ignore_changes = [tags]
  }
}

resource "aws_transfer_tag" "test" {
  resource_arn = aws_transfer_server.test.arn
  key          = %[2]q
  value        = %[3]q
}
`, rName, key, value)
}
