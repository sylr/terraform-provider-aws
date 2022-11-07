package ses_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/ses"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func testAccActiveReceiptRuleSetDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "data.aws_ses_active_receipt_rule_set.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			testAccPreCheck(t)
			testAccPreCheckReceiptRule(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, ses.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckActiveReceiptRuleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccActiveReceiptRuleSetDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckActiveReceiptRuleSetExists(resourceName),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "ses", fmt.Sprintf("receipt-rule-set/%s", rName)),
				),
			},
		},
	})
}

func testAccActiveReceiptRuleSetDataSource_noActiveRuleSet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			testAccPreCheck(t)
			testAccPreCheckUnsetActiveRuleSet(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, ses.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccActiveReceiptRuleSetDataSourceConfig_noActiveRuleSet(),
				ExpectError: regexp.MustCompile("empty result"),
			},
		},
	})
}

func testAccActiveReceiptRuleSetDataSourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "aws_ses_receipt_rule_set" "test" {
  rule_set_name = %[1]q
}

resource "aws_ses_active_receipt_rule_set" "test" {
  rule_set_name = aws_ses_receipt_rule_set.test.rule_set_name
}

data "aws_ses_active_receipt_rule_set" "test" {
  depends_on = [aws_ses_active_receipt_rule_set.test]
}
`, name)
}

func testAccActiveReceiptRuleSetDataSourceConfig_noActiveRuleSet() string {
	return `
data "aws_ses_active_receipt_rule_set" "test" {}
`
}

func testAccPreCheckUnsetActiveRuleSet(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SESConn

	output, err := conn.DescribeActiveReceiptRuleSet(&ses.DescribeActiveReceiptRuleSetInput{})
	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}
	if output == nil || output.Metadata == nil {
		return
	}
	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}

	_, err = conn.SetActiveReceiptRuleSet(&ses.SetActiveReceiptRuleSetInput{
		RuleSetName: nil,
	})
	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}
	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}
