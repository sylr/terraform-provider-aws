package sesv2_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfsesv2 "github.com/hashicorp/terraform-provider-aws/internal/service/sesv2"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccSESV2DedicatedIPAssignment_serial(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic":      testAccSESV2DedicatedIPAssignment_basic,
		"disappears": testAccSESV2DedicatedIPAssignment_disappears,
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccSESV2DedicatedIPAssignment_basic(t *testing.T) { // nosemgrep:ci.sesv2-in-func-name
	if os.Getenv("SES_DEDICATED_IP") == "" {
		t.Skip("Environment variable SES_DEDICATED_IP is not set")
	}

	ip := os.Getenv("SES_DEDICATED_IP")
	poolName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sesv2_dedicated_ip_assignment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDedicatedIPAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedIPAssignmentConfig_basic(ip, poolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedIPAssignmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "destination_pool_name", poolName),
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

func testAccSESV2DedicatedIPAssignment_disappears(t *testing.T) { // nosemgrep:ci.sesv2-in-func-name
	if os.Getenv("SES_DEDICATED_IP") == "" {
		t.Skip("Environment variable SES_DEDICATED_IP is not set")
	}

	ip := os.Getenv("SES_DEDICATED_IP")
	poolName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sesv2_dedicated_ip_assignment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESV2EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDedicatedIPAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedIPAssignmentConfig_basic(ip, poolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedIPAssignmentExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfsesv2.ResourceDedicatedIPAssignment(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckDedicatedIPAssignmentDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SESV2Client
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_sesv2_dedicated_ip_assignment" {
			continue
		}

		_, err := tfsesv2.FindDedicatedIPAssignmentByID(ctx, conn, rs.Primary.ID)
		if err != nil {
			var nfe *types.NotFoundException
			if errors.As(err, &nfe) {
				return nil
			}
			if errors.Is(err, tfsesv2.ErrIncorrectPoolAssignment) {
				return nil
			}
			return err
		}

		return create.Error(names.SESV2, create.ErrActionCheckingDestroyed, tfsesv2.ResNameDedicatedIPAssignment, rs.Primary.ID, errors.New("not destroyed"))
	}

	return nil
}

func testAccCheckDedicatedIPAssignmentExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.SESV2, create.ErrActionCheckingExistence, tfsesv2.ResNameDedicatedIPAssignment, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.SESV2, create.ErrActionCheckingExistence, tfsesv2.ResNameDedicatedIPAssignment, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SESV2Client
		ctx := context.Background()
		_, err := tfsesv2.FindDedicatedIPAssignmentByID(ctx, conn, rs.Primary.ID)
		if err != nil {
			return create.Error(names.SESV2, create.ErrActionCheckingExistence, tfsesv2.ResNameDedicatedIPAssignment, rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccDedicatedIPAssignmentConfig_basic(ip, poolName string) string {
	return fmt.Sprintf(`
resource "aws_sesv2_dedicated_ip_pool" "test" {
  pool_name = %[2]q
}

resource "aws_sesv2_dedicated_ip_assignment" "test" {
  ip                    = %[1]q
  destination_pool_name = aws_sesv2_dedicated_ip_pool.test.pool_name
}
`, ip, poolName)
}
