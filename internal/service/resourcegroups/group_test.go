package resourcegroups_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/resourcegroups"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfresourcegroups "github.com/hashicorp/terraform-provider-aws/internal/service/resourcegroups"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccResourceGroupsGroup_basic(t *testing.T) {
	var v resourcegroups.Group
	resourceName := "aws_resourcegroups_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	desc1 := "Hello World"
	desc2 := "Foo Bar"

	query2 := `{
  "ResourceTypeFilters": [
    "AWS::EC2::Instance"
  ],
  "TagFilters": [
    {
      "Key": "Hello",
      "Values": [
        "World"
      ]
    }
  ]
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, resourcegroups.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_basic(rName, desc1, testAccResourceGroupQueryConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGroupExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", desc1),
					resource.TestCheckResourceAttr(resourceName, "resource_query.0.query", testAccResourceGroupQueryConfig+"\n"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGroupConfig_basic(rName, desc2, query2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", desc2),
					resource.TestCheckResourceAttr(resourceName, "resource_query.0.query", query2+"\n"),
				),
			},
		},
	})
}

func TestAccResourceGroupsGroup_tags(t *testing.T) {
	var v resourcegroups.Group
	resourceName := "aws_resourcegroups_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	desc1 := "Hello World"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, resourcegroups.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_tags1(rName, desc1, testAccResourceGroupQueryConfig, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGroupExists(resourceName, &v),
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
				Config: testAccGroupConfig_tags2(rName, desc1, testAccResourceGroupQueryConfig, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGroupExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccGroupConfig_tags1(rName, desc1, testAccResourceGroupQueryConfig, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGroupExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccResourceGroupsGroup_Configuration(t *testing.T) {
	var v resourcegroups.Group
	resourceName := "aws_resourcegroups_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	desc1 := "Hello World"
	desc2 := "Foo Bar"
	configType1 := "AWS::EC2::HostManagement"
	configType2 := "AWS::ResourceGroups::Generic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, resourcegroups.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_configuration(rName, desc1, configType1, configType2, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGroupExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", desc1),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.type", configType1),
					resource.TestCheckResourceAttr(resourceName, "configuration.1.type", configType2),
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.parameters.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.parameters.0.name", "allowed-host-families"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.parameters.0.values.0", "mac1"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Check that changing the auto-allocate value is represented
			{
				Config: testAccGroupConfig_configuration(rName, desc1, configType1, configType2, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceGroupExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", desc1),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.type", configType1),
					resource.TestCheckResourceAttr(resourceName, "configuration.1.type", configType2),
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.parameters.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.parameters.2.name", "auto-allocate-host"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.parameters.2.values.0", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
				),
			},
			{
				Config: testAccGroupConfig_configuration(rName, desc2, configType1, configType2, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", desc2),
				),
			},
			// Check that trying to change the configuration group to a resource-query group fails
			{
				Config:      testAccGroupConfig_basic(rName, desc1, testAccResourceGroupQueryConfig),
				ExpectError: regexp.MustCompile(`conversion between resource-query and configuration group types is not possible`),
			},
		},
	})
}

func testAccCheckResourceGroupExists(n string, v *resourcegroups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Resource Groups Group ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ResourceGroupsConn

		output, err := tfresourcegroups.FindGroupByName(context.Background(), conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckResourceGroupDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ResourceGroupsConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_resourcegroups_group" {
			continue
		}

		_, err := tfresourcegroups.FindGroupByName(context.Background(), conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Resource Groups Group %s still exists", rs.Primary.ID)
	}

	return nil
}

const testAccResourceGroupQueryConfig = `{
  "ResourceTypeFilters": [
    "AWS::EC2::Instance"
  ],
  "TagFilters": [
    {
      "Key": "Stage",
      "Values": [
        "Test"
      ]
    }
  ]
}`

func testAccGroupConfig_basic(rName, desc, query string) string {
	return fmt.Sprintf(`
resource "aws_resourcegroups_group" "test" {
  name        = %[1]q
  description = %[2]q

  resource_query {
    query = <<JSON
%[3]s
JSON

  }
}
`, rName, desc, query)
}

func testAccGroupConfig_tags1(rName, desc, query, tag1Key, tag1Value string) string {
	return fmt.Sprintf(`
resource "aws_resourcegroups_group" "test" {
  name        = %[1]q
  description = %[2]q

  resource_query {
    query = <<JSON
%[3]s
JSON

  }

  tags = {
    %[4]q = %[5]q
  }
}
`, rName, desc, query, tag1Key, tag1Value)
}

func testAccGroupConfig_tags2(rName, desc, query, tag1Key, tag1Value, tag2Key, tag2Value string) string {
	return fmt.Sprintf(`
resource "aws_resourcegroups_group" "test" {
  name        = %[1]q
  description = %[2]q

  resource_query {
    query = <<JSON
%[3]s
JSON

  }

  tags = {
    %[4]q = %[5]q
    %[6]q = %[7]q
  }
}
`, rName, desc, query, tag1Key, tag1Value, tag2Key, tag2Value)
}

func testAccGroupConfig_configuration(rName, desc, cType1, cType2 string, autoAllocateHost bool) string {
	return fmt.Sprintf(`
resource "aws_resourcegroups_group" "test" {
  name        = %[1]q
  description = %[2]q

  configuration {
    type = %[3]q

    parameters {
      name = "allowed-host-families"
      values = [
        "mac1",
      ]
    }

    parameters {
      name = "any-host-based-license-configuration"
      values = [
        "true",
      ]
    }

    parameters {
      name = "auto-allocate-host"
      values = [
        "%[4]t",
      ]
    }

    parameters {
      name = "auto-release-host"
      values = [
        "true",
      ]
    }
  }

  configuration {
    type = %[5]q

    parameters {
      name = "allowed-resource-types"
      values = [
        "AWS::EC2::Host",
      ]
    }

    parameters {
      name = "deletion-protection"
      values = [
        "UNLESS_EMPTY"
      ]
    }
  }
}
`, rName, desc, cType1, autoAllocateHost, cType2)
}
