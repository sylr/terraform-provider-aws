package iot_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccIoTTopicRule_basic(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
					resource.TestCheckResourceAttr("aws_iot_topic_rule.rule", "name", fmt.Sprintf("test_rule_%s", rName)),
					resource.TestCheckResourceAttr("aws_iot_topic_rule.rule", "description", "Example rule"),
					resource.TestCheckResourceAttr("aws_iot_topic_rule.rule", "enabled", "true"),
					resource.TestCheckResourceAttr("aws_iot_topic_rule.rule", "sql", "SELECT * FROM 'topic/test'"),
					resource.TestCheckResourceAttr("aws_iot_topic_rule.rule", "sql_version", "2015-10-08"),
					resource.TestCheckResourceAttr("aws_iot_topic_rule.rule", "tags.%", "0"),
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

func TestAccIoTTopicRule_cloudWatchAlarm(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_cloudWatchalarm(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_cloudWatchLogs(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIoTTopicRule_cloudWatchLogs(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_cloudWatchMetric(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_cloudWatchmetric(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_dynamoDB(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_dynamoDB(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTopicRule_dynamoDB_rangekeys(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
				),
			},
		},
	})
}

func TestAccIoTTopicRule_dynamoDBv2(t *testing.T) {
	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_dynamoDBv2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
				),
			},
		},
	})
}

func TestAccIoTTopicRule_elasticSearch(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_elasticSearch(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_firehose(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_firehose(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_Firehose_separator(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_firehose_separator(rName, "\n"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTopicRule_firehose_separator(rName, ","),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
				),
			},
		},
	})
}

func TestAccIoTTopicRule_kinesis(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_kinesis(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_lambda(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_lambda(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_republish(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_republish(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_republishWithQos(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_republish_with_qos(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_s3(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_s3(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_sns(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_sns(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_sqs(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_sqs(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_Step_functions(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_step_functions(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

func TestAccIoTTopicRule_IoT_analytics(t *testing.T) {
	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_iot_analytics(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
				),
			},
		},
	})
}

func TestAccIoTTopicRule_IoT_events(t *testing.T) {
	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_iot_events(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
				),
			},
		},
	})
}

func TestAccIoTTopicRule_tags(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRuleTags1(rName, "key1", "user@example"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "user@example"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTopicRuleTags2(rName, "key1", "user@example", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "user@example"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccTopicRuleTags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccIoTTopicRule_errorAction(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_errorAction(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
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

// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/16115
func TestAccIoTTopicRule_updateKinesisErrorAction(t *testing.T) {
	rName := sdkacctest.RandString(5)
	resourceName := "aws_iot_topic_rule.rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, iot.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckTopicRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTopicRule_kinesis(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
					resource.TestCheckResourceAttr(resourceName, "error_action.#", "0"),
				),
			},

			{
				Config: testAccTopicRule_errorAction(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTopicRuleExists("aws_iot_topic_rule.rule"),
					resource.TestCheckResourceAttr(resourceName, "error_action.#", "1"),
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

func testAccCheckTopicRuleDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).IoTConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_iot_topic_rule" {
			continue
		}

		input := &iot.ListTopicRulesInput{}

		out, err := conn.ListTopicRules(input)

		if err != nil {
			return err
		}

		for _, r := range out.Rules {
			if *r.RuleName == rs.Primary.ID {
				return fmt.Errorf("IoT topic rule still exists:\n%s", r)
			}
		}

	}

	return nil
}

func testAccCheckTopicRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).IoTConn
		input := &iot.ListTopicRulesInput{}

		output, err := conn.ListTopicRules(input)

		if err != nil {
			return err
		}

		for _, rule := range output.Rules {
			if aws.StringValue(rule.RuleName) == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("IoT Topic Rule (%s) not found", rs.Primary.ID)
	}
}

func testAccTopicRuleRole(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "iot_role" {
  name = "test_role_%[1]s"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "iot.${data.aws_partition.current.dns_suffix}"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF

}

resource "aws_iam_policy" "policy" {
  name        = "test_policy_%[1]s"
  path        = "/"
  description = "My test policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "*",
      "Resource": "*"
    }
  ]
}
EOF

}

resource "aws_iam_policy_attachment" "attach_policy" {
  name       = "test_policy_attachment_%[1]s"
  roles      = [aws_iam_role.iot_role.name]
  policy_arn = aws_iam_policy.policy.arn
}
`, rName)
}

func testAccTopicRule_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"
}
`, rName)
}

func testAccTopicRule_cloudWatchalarm(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  cloudwatch_alarm {
    alarm_name   = "myalarm"
    role_arn     = aws_iam_role.iot_role.arn
    state_reason = "test"
    state_value  = "OK"
  }
}
`, rName))
}

func testAccIoTTopicRule_cloudWatchLogs(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  cloudwatch_logs {
    log_group_name = "mylogs"
    role_arn       = aws_iam_role.iot_role.arn
  }
}
`, rName))
}

func testAccTopicRule_cloudWatchmetric(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  cloudwatch_metric {
    metric_name      = "FakeData"
    metric_namespace = "FakeData"
    metric_value     = "FakeData"
    metric_unit      = "FakeData"
    role_arn         = aws_iam_role.iot_role.arn
  }
}
`, rName))
}

func testAccTopicRule_dynamoDBv2(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT field as column_name FROM 'topic/test'"
  sql_version = "2015-10-08"

  dynamodbv2 {
    put_item {
      table_name = "test"
    }

    role_arn = aws_iam_role.iot_role.arn
  }
}
`, rName))
}

func testAccTopicRule_dynamoDB(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  dynamodb {
    hash_key_field = "hash_key_field"
    hash_key_value = "hash_key_value"
    payload_field  = "payload_field"
    role_arn       = aws_iam_role.iot_role.arn
    table_name     = "table_name"
  }
}
`, rName))
}

func testAccTopicRule_dynamoDB_rangekeys(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  dynamodb {
    hash_key_field  = "hash_key_field"
    hash_key_value  = "hash_key_value"
    payload_field   = "payload_field"
    range_key_field = "range_key_field"
    range_key_value = "range_key_value"
    range_key_type  = "STRING"
    role_arn        = aws_iam_role.iot_role.arn
    table_name      = "table_name"
    operation       = "INSERT"
  }
}
`, rName))
}

func testAccTopicRule_elasticSearch(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
data "aws_region" "current" {}

resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  elasticsearch {
    endpoint = "https://domain.${data.aws_region.current.name}.es.${data.aws_partition.current.dns_suffix}"
    id       = "myIdentifier"
    index    = "myindex"
    type     = "mydocument"
    role_arn = aws_iam_role.iot_role.arn
  }
}
`, rName))
}

func testAccTopicRule_firehose(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  firehose {
    delivery_stream_name = "mystream"
    role_arn             = aws_iam_role.iot_role.arn
  }
}
`, rName))
}

func testAccTopicRule_firehose_separator(rName, separator string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  firehose {
    delivery_stream_name = "mystream"
    role_arn             = aws_iam_role.iot_role.arn
    separator            = %q
  }
}
`, rName, separator))
}

func testAccTopicRule_kinesis(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  kinesis {
    stream_name = "mystream"
    role_arn    = aws_iam_role.iot_role.arn
  }
}
`, rName))
}

func testAccTopicRule_lambda(rName string) string {
	return fmt.Sprintf(`
data "aws_region" "current" {}

data "aws_partition" "current" {}

resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  lambda {
    function_arn = "arn:${data.aws_partition.current.partition}:lambda:${data.aws_region.current.name}:123456789012:function:ProcessKinesisRecords"
  }
}
`, rName)
}

func testAccTopicRule_republish(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  republish {
    role_arn = aws_iam_role.iot_role.arn
    topic    = "mytopic"
  }
}
`, rName))
}

func testAccTopicRule_republish_with_qos(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  republish {
    role_arn = aws_iam_role.iot_role.arn
    topic    = "mytopic"
    qos      = 1
  }
}
`, rName))
}

func testAccTopicRule_s3(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  s3 {
    bucket_name = "mybucket"
    key         = "mykey"
    role_arn    = aws_iam_role.iot_role.arn
  }
}
`, rName))
}

func testAccTopicRule_sns(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
data "aws_region" "current" {}

resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  sns {
    role_arn   = aws_iam_role.iot_role.arn
    target_arn = "arn:${data.aws_partition.current.partition}:sns:${data.aws_region.current.name}:123456789012:my_corporate_topic"
  }
}
`, rName))
}

func testAccTopicRule_sqs(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  sqs {
    queue_url  = "fakedata"
    role_arn   = aws_iam_role.iot_role.arn
    use_base64 = false
  }
}
`, rName))
}

func testAccTopicRule_step_functions(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  step_functions {
    execution_name_prefix = "myprefix"
    state_machine_name    = "mystatemachine"
    role_arn              = aws_iam_role.iot_role.arn
  }
}
`, rName))
}

func testAccTopicRule_iot_analytics(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  iot_analytics {
    channel_name = "fakedata"
    role_arn     = aws_iam_role.iot_role.arn
  }
}
`, rName))
}

func testAccTopicRule_iot_events(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  iot_events {
    input_name = "fake_input_name"
    role_arn   = aws_iam_role.iot_role.arn
    message_id = "fake_message_id"
  }
}
`, rName))
}

func testAccTopicRuleTags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_iot_topic_rule" "test" {
  name        = "test_rule_%[1]s"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccTopicRuleTags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_iot_topic_rule" "test" {
  name        = "test_rule_%[1]s"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccTopicRule_errorAction(rName string) string {
	return acctest.ConfigCompose(
		testAccTopicRuleRole(rName),
		fmt.Sprintf(`
resource "aws_iot_topic_rule" "rule" {
  name        = "test_rule_%[1]s"
  description = "Example rule"
  enabled     = true
  sql         = "SELECT * FROM 'topic/test'"
  sql_version = "2015-10-08"

  kinesis {
    stream_name = "mystream"
    role_arn    = aws_iam_role.iot_role.arn
  }

  error_action {
    kinesis {
      stream_name = "mystream"
      role_arn    = aws_iam_role.iot_role.arn
    }
  }
}
`, rName))
}
