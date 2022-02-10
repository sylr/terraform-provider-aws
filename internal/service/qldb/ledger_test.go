package qldb_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/qldb"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccQLDBLedger_basic(t *testing.T) {
	var qldbCluster qldb.DescribeLedgerOutput
	rInt := sdkacctest.RandInt()
	resourceName := "aws_qldb_ledger.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(qldb.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, qldb.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckLedgerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLedgerConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLedgerExists(resourceName, &qldbCluster),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "qldb", regexp.MustCompile(`ledger/.+`)),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile("test-ledger-[0-9]+")),
					resource.TestCheckResourceAttr(resourceName, "permissions_mode", "ALLOW_ALL"),
					resource.TestCheckResourceAttr(resourceName, "deletion_protection", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccQLDBLedger_update(t *testing.T) {
	var qldbCluster qldb.DescribeLedgerOutput
	rInt := sdkacctest.RandInt()
	resourceName := "aws_qldb_ledger.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(qldb.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, qldb.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckLedgerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLedgerConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLedgerExists(resourceName, &qldbCluster),
					resource.TestCheckResourceAttr(resourceName, "permissions_mode", "ALLOW_ALL"),
				),
			},
			{
				Config: testAccLedgerConfig_update(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLedgerExists(resourceName, &qldbCluster),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "qldb", regexp.MustCompile(`ledger/.+`)),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile("test-ledger-[0-9]+")),
					resource.TestCheckResourceAttr(resourceName, "permissions_mode", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "deletion_protection", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func testAccCheckLedgerDestroy(s *terraform.State) error {
	return testAccCheckLedgerDestroyWithProvider(s, acctest.Provider)
}

func testAccCheckLedgerDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*conns.AWSClient).QLDBConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_qldb_ledger" {
			continue
		}

		// Try to find the Group
		var err error
		resp, err := conn.DescribeLedger(
			&qldb.DescribeLedgerInput{
				Name: aws.String(rs.Primary.ID),
			})

		if err == nil {
			if len(aws.StringValue(resp.Name)) != 0 && aws.StringValue(resp.Name) == rs.Primary.ID {
				return fmt.Errorf("QLDB Ledger %s still exists", rs.Primary.ID)
			}
		}

		// Return nil if the cluster is already destroyed
		if tfawserr.ErrMessageContains(err, qldb.ErrCodeResourceNotFoundException, "") {
			continue
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckLedgerExists(n string, v *qldb.DescribeLedgerOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No QLDB Ledger ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).QLDBConn
		resp, err := conn.DescribeLedger(&qldb.DescribeLedgerInput{
			Name: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		if *resp.Name == rs.Primary.ID {
			*v = *resp
			return nil
		}

		return fmt.Errorf("QLDB Ledger (%s) not found", rs.Primary.ID)
	}
}

func testAccLedgerConfig_basic(n int) string {
	return fmt.Sprintf(`
resource "aws_qldb_ledger" "test" {
  name                = "test-ledger-%d"
  permissions_mode    = "ALLOW_ALL"
  deletion_protection = false
}
`, n)
}

func testAccLedgerConfig_update(n int) string {
	return fmt.Sprintf(`
resource "aws_qldb_ledger" "test" {
  name                = "test-ledger-%d"
  permissions_mode    = "STANDARD"
  deletion_protection = false
}
`, n)
}

func TestAccQLDBLedger_tags(t *testing.T) {
	var cluster1, cluster2, cluster3 qldb.DescribeLedgerOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_qldb_ledger.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(qldb.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, qldb.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckLedgerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLedgerTags1Config(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLedgerExists(resourceName, &cluster1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLedgerTags2Config(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLedgerExists(resourceName, &cluster2),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccLedgerTags1Config(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLedgerExists(resourceName, &cluster3),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccLedgerTags1Config(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_qldb_ledger" "test" {
  name                = %[1]q
  permissions_mode    = "ALLOW_ALL"
  deletion_protection = false

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccLedgerTags2Config(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_qldb_ledger" "test" {
  name                = %[1]q
  permissions_mode    = "ALLOW_ALL"
  deletion_protection = false

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
