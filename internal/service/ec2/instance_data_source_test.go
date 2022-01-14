package ec2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccEC2InstanceDataSource_basic(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "arn", resourceName, "arn"),
					resource.TestCheckNoResourceAttr(datasourceName, "user_data_base64"),
					resource.TestCheckResourceAttr(datasourceName, "outpost_arn", ""),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_tags(t *testing.T) {
	rInt := sdkacctest.RandInt()
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_Tags(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_azUserData(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_AzUserData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "availability_zone", resourceName, "availability_zone"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_data", resourceName, "user_data"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_gp2IopsDevice(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_gp2IopsDevice,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.#", resourceName, "root_block_device.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.volume_size", resourceName, "root_block_device.0.volume_size"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.volume_type", resourceName, "root_block_device.0.volume_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.device_name", resourceName, "root_block_device.0.device_name"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.iops", resourceName, "root_block_device.0.iops"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_gp3ThroughputDevice(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_gp3ThroughputDevice,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.#", resourceName, "root_block_device.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.volume_size", resourceName, "root_block_device.0.volume_size"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.volume_type", resourceName, "root_block_device.0.volume_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.device_name", resourceName, "root_block_device.0.device_name"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.throughput", resourceName, "root_block_device.0.throughput"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_blockDevices(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_blockDevices,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.#", resourceName, "root_block_device.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.volume_size", resourceName, "root_block_device.0.volume_size"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.volume_type", resourceName, "root_block_device.0.volume_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.0.device_name", resourceName, "root_block_device.0.device_name"),
					resource.TestCheckResourceAttrPair(datasourceName, "ebs_block_device.#", resourceName, "ebs_block_device.#"),
					//resource.TestCheckResourceAttrPair(datasourceName, "ephemeral_block_device.#", resourceName, "ephemeral_block_device.#"),
					// ephemeral block devices don't get saved properly due to API limitations, so this can't actually be tested right now
				),
			},
		},
	})
}

// Test to verify that ebs_block_device kms_key_id does not elicit a panic
func TestAccEC2InstanceDataSource_EBSBlockDevice_kmsKeyID(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_EbsBlockDevice_KmsKeyId,
			},
		},
	})
}

// Test to verify that root_block_device kms_key_id does not elicit a panic
func TestAccEC2InstanceDataSource_RootBlockDevice_kmsKeyID(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_RootBlockDevice_KmsKeyId,
			},
		},
	})
}

func TestAccEC2InstanceDataSource_rootInstanceStore(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_rootInstanceStore,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "ebs_block_device.#", resourceName, "ebs_block_device.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "ebs_optimized", resourceName, "ebs_optimized"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_block_device.#", resourceName, "root_block_device.#"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_privateIP(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", sdkacctest.RandString(12))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_privateIP(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "private_ip", resourceName, "private_ip"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_secondaryPrivateIPs(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", sdkacctest.RandString(12))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_secondaryPrivateIPs(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "secondary_private_ips", resourceName, "secondary_private_ips"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_ipv6Addresses(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_ipv6Addresses(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "ipv6_addresses.#", resourceName, "ipv6_address_count"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_keyPair(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-test-key-%d", sdkacctest.RandInt())

	publicKey, _, err := sdkacctest.RandSSHKeyPair(acctest.DefaultEmailAddress)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_keyPair(rName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "key_name", resourceName, "key_name"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_vpc(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", sdkacctest.RandString(12))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_VPC(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_data", resourceName, "user_data"),
					resource.TestCheckResourceAttrPair(datasourceName, "associate_public_ip_address", resourceName, "associate_public_ip_address"),
					resource.TestCheckResourceAttrPair(datasourceName, "tenancy", resourceName, "tenancy"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_placementGroup(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", sdkacctest.RandString(12))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_PlacementGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "placement_group", resourceName, "placement_group"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_securityGroups(t *testing.T) {
	rInt := sdkacctest.RandInt()
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_SecurityGroups(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_data", resourceName, "user_data"),
					resource.TestCheckResourceAttrPair(datasourceName, "vpc_security_group_ids.#", resourceName, "vpc_security_group_ids.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "security_groups.#", resourceName, "security_groups.#"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_vpcSecurityGroups(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", sdkacctest.RandString(12))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_VPCSecurityGroups(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "ami", resourceName, "ami"),
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "vpc_security_group_ids.#", resourceName, "vpc_security_group_ids.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "security_groups.#", resourceName, "security_groups.#"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_GetPasswordData_trueToFalse(t *testing.T) {
	datasourceName := "data.aws_instance.test"

	rName := fmt.Sprintf("tf-testacc-instance-%s", sdkacctest.RandString(12))
	publicKey, _, err := sdkacctest.RandSSHKeyPair(acctest.DefaultEmailAddress)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_getPasswordData(rName, publicKey, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_password_data", "true"),
					resource.TestCheckResourceAttrSet(datasourceName, "password_data"),
				),
			},
			{
				Config: testAccInstanceDataSourceConfig_getPasswordData(rName, publicKey, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_password_data", "false"),
					resource.TestCheckNoResourceAttr(datasourceName, "password_data"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_GetPasswordData_falseToTrue(t *testing.T) {
	datasourceName := "data.aws_instance.test"

	rName := fmt.Sprintf("tf-testacc-instance-%s", sdkacctest.RandString(12))
	publicKey, _, err := sdkacctest.RandSSHKeyPair(acctest.DefaultEmailAddress)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_getPasswordData(rName, publicKey, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_password_data", "false"),
					resource.TestCheckNoResourceAttr(datasourceName, "password_data"),
				),
			},
			{
				Config: testAccInstanceDataSourceConfig_getPasswordData(rName, publicKey, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_password_data", "true"),
					resource.TestCheckResourceAttrSet(datasourceName, "password_data"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_getUserData(t *testing.T) {
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", sdkacctest.RandString(12))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfigGetUserData(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_user_data", "true"),
					resource.TestCheckResourceAttr(datasourceName, "user_data_base64", "IyEvYmluL2Jhc2gKCmVjaG8gImhlbGxvIHdvcmxkIgo="),
				),
			},
			{
				Config: testAccInstanceDataSourceConfigGetUserData(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_user_data", "false"),
					resource.TestCheckNoResourceAttr(datasourceName, "user_data_base64"),
				),
			},
			{
				Config: testAccInstanceDataSourceConfigGetUserData(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_user_data", "true"),
					resource.TestCheckResourceAttr(datasourceName, "user_data_base64", "IyEvYmluL2Jhc2gKCmVjaG8gImhlbGxvIHdvcmxkIgo="),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_GetUserData_noUserData(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", sdkacctest.RandString(12))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfigGetUserDataNoUserData(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_user_data", "true"),
					resource.TestCheckNoResourceAttr(datasourceName, "user_data_base64"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_data_base64", resourceName, "user_data_base64"),
				),
			},
			{
				Config: testAccInstanceDataSourceConfigGetUserDataNoUserData(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_user_data", "false"),
					resource.TestCheckNoResourceAttr(datasourceName, "user_data_base64"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_data_base64", resourceName, "user_data_base64"),
				),
			},
			{
				Config: testAccInstanceDataSourceConfigGetUserDataNoUserData(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "get_user_data", "true"),
					resource.TestCheckNoResourceAttr(datasourceName, "user_data_base64"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_data_base64", resourceName, "user_data_base64"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_creditSpecification(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := fmt.Sprintf("tf-testacc-instance-%s", sdkacctest.RandString(12))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{

				Config: testAccInstanceDataSourceConfig_creditSpecification(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "credit_specification.#", resourceName, "credit_specification.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "credit_specification.0.cpu_credits", resourceName, "credit_specification.0.cpu_credits"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_metadataOptions(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_metadataOptions(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "metadata_options.#", resourceName, "metadata_options.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "metadata_options.0.http_endpoint", resourceName, "metadata_options.0.http_endpoint"),
					resource.TestCheckResourceAttrPair(datasourceName, "metadata_options.0.http_tokens", resourceName, "metadata_options.0.http_tokens"),
					resource.TestCheckResourceAttrPair(datasourceName, "metadata_options.0.http_put_response_hop_limit", resourceName, "metadata_options.0.http_put_response_hop_limit"),
					resource.TestCheckResourceAttrPair(datasourceName, "metadata_options.0.instance_metadata_tags", resourceName, "metadata_options.0.instance_metadata_tags"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_enclaveOptions(t *testing.T) {
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_enclaveOptions(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "enclave_options.#", resourceName, "enclave_options.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "enclave_options.0.enabled", resourceName, "enclave_options.0.enabled"),
				),
			},
		},
	})
}

func TestAccEC2InstanceDataSource_blockDeviceTags(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_instance.test"
	datasourceName := "data.aws_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_blockDeviceTags(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "instance_type", resourceName, "instance_type"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
var testAccInstanceDataSourceConfig = acctest.ConfigLatestAmazonLinuxHvmEbsAmi() + `
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.small"

  tags = {
    Name = "HelloWorld"
  }
}

data "aws_instance" "test" {
  filter {
    name   = "instance-id"
    values = [aws_instance.test.id]
  }
}
`

// Use the tags attribute to filter
func testAccInstanceDataSourceConfig_Tags(rInt int) string {
	return acctest.ConfigLatestAmazonLinuxHvmEbsAmi() + fmt.Sprintf(`
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.small"

  tags = {
    Name     = "HelloWorld"
    TestSeed = "%[1]d"
  }
}

data "aws_instance" "test" {
  instance_tags = {
    Name     = aws_instance.test.tags["Name"]
    TestSeed = "%[1]d"
  }
}
`, rInt)
}

// filter on tag, populate more attributes
var testAccInstanceDataSourceConfig_AzUserData = acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(),
	acctest.ConfigLatestAmazonLinuxHvmEbsAmi(), `
resource "aws_instance" "test" {
  ami               = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  availability_zone = data.aws_availability_zones.available.names[0]

  instance_type = "t2.micro"
  user_data     = "test:-with-character's"

  tags = {
    TFAccTest = "YesThisIsATest"
  }
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`)

// GP2IopsDevice
var testAccInstanceDataSourceConfig_gp2IopsDevice = acctest.ConfigLatestAmazonLinuxHvmEbsAmi() + `
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t3.medium"

  root_block_device {
    volume_type = "gp2"
    volume_size = 11
  }
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`

// GP3ThroughputDevice
var testAccInstanceDataSourceConfig_gp3ThroughputDevice = acctest.ConfigLatestAmazonLinuxHvmEbsAmi() + `
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t3.medium"

  root_block_device {
    volume_type = "gp3"
    volume_size = 10
    throughput  = 300
  }
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`

// Block Device
var testAccInstanceDataSourceConfig_blockDevices = acctest.ConfigLatestAmazonLinuxHvmEbsAmi() + `
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t3.medium"

  root_block_device {
    volume_type = "gp2"
    volume_size = 11
  }

  ebs_block_device {
    device_name = "/dev/sdb"
    volume_size = 9
  }

  ebs_block_device {
    device_name = "/dev/sdc"
    volume_size = 10
    volume_type = "io1"
    iops        = 100
  }

  # Encrypted ebs block device
  ebs_block_device {
    device_name = "/dev/sdd"
    volume_size = 12
    encrypted   = true
  }

  ephemeral_block_device {
    device_name  = "/dev/sde"
    virtual_name = "ephemeral0"
  }

  ebs_block_device {
    device_name = "/dev/sdf"
    volume_size = 10
    volume_type = "gp3"
    throughput  = 300
  }
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`

var testAccInstanceDataSourceConfig_EbsBlockDevice_KmsKeyId = acctest.ConfigLatestAmazonLinuxHvmEbsAmi() + `
resource "aws_kms_key" "test" {
  deletion_window_in_days = 7
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t3.medium"

  root_block_device {
    volume_type = "gp2"
    volume_size = 11
  }

  ebs_block_device {
    device_name = "/dev/sdb"
    encrypted   = true
    kms_key_id  = aws_kms_key.test.arn
    volume_size = 9
  }
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`

var testAccInstanceDataSourceConfig_RootBlockDevice_KmsKeyId = acctest.ConfigLatestAmazonLinuxHvmEbsAmi() + `
resource "aws_kms_key" "test" {
  deletion_window_in_days = 7
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t3.medium"

  root_block_device {
    encrypted   = true
    kms_key_id  = aws_kms_key.test.arn
    volume_type = "gp2"
    volume_size = 11
  }
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`

var testAccInstanceDataSourceConfig_rootInstanceStore = acctest.ConfigLatestAmazonLinuxHvmEbsAmi() + `
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t3.medium"
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`

func testAccInstanceDataSourceConfig_privateIP(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigLatestAmazonLinuxHvmEbsAmi(), testAccInstanceVPCBasicConfig(rName), `
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.micro"
  subnet_id     = aws_subnet.test.id
  private_ip    = "10.0.0.42"
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`)
}

func testAccInstanceDataSourceConfig_secondaryPrivateIPs(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigLatestAmazonLinuxHvmEbsAmi(), testAccInstanceVPCBasicConfig(rName), `
resource "aws_instance" "test" {
  ami                   = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type         = "t2.micro"
  subnet_id             = aws_subnet.test.id
  secondary_private_ips = ["10.0.0.42"]
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`)
}

func testAccInstanceDataSourceConfig_ipv6Addresses(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigLatestAmazonLinuxHvmEbsAmi(), testAccInstanceVPCIPv6Config(rName), fmt.Sprintf(`
resource "aws_instance" "test" {
  ami                = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type      = "t2.micro"
  subnet_id          = aws_subnet.test.id
  ipv6_address_count = 1

  tags = {
    Name = %[1]q
  }
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`, rName))
}

func testAccInstanceDataSourceConfig_keyPair(rName, publicKey string) string {
	return acctest.ConfigCompose(acctest.ConfigLatestAmazonLinuxHvmEbsAmi(), fmt.Sprintf(`
resource "aws_key_pair" "test" {
  key_name   = %[1]q
  public_key = %[2]q
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.micro"
  key_name      = aws_key_pair.test.key_name

  tags = {
    Name = %[1]q
  }
}

data "aws_instance" "test" {
  filter {
    name   = "tag:Name"
    values = [%[1]q]
  }

  filter {
    name   = "key-name"
    values = [aws_instance.test.key_name]
  }
}
`, rName, publicKey))
}

func testAccInstanceDataSourceConfig_VPC(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigLatestAmazonLinuxHvmEbsAmi(), testAccInstanceVPCBasicConfig(rName), `
resource "aws_instance" "test" {
  ami                         = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type               = "t2.small"
  subnet_id                   = aws_subnet.test.id
  associate_public_ip_address = true
  #tenancy                     = "dedicated"
  # pre-encoded base64 data
  user_data = "3dc39dda39be1205215e776bad998da361a5955d"
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`)
}

func testAccInstanceDataSourceConfig_PlacementGroup(rName string) string {
	return acctest.ConfigLatestAmazonLinuxHvmEbsAmi() + testAccInstanceVPCBasicConfig(rName) + fmt.Sprintf(`
resource "aws_placement_group" "test" {
  name     = %[1]q
  strategy = "cluster"
}

# Limitations: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/placement-groups.html#concepts-placement-groups
resource "aws_instance" "test" {
  ami                         = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type               = "c5.large"
  subnet_id                   = aws_subnet.test.id
  associate_public_ip_address = true
  placement_group             = aws_placement_group.test.name

  # pre-encoded base64 data
  user_data = "3dc39dda39be1205215e776bad998da361a5955d"
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`, rName)
}

func testAccInstanceDataSourceConfig_SecurityGroups(rInt int) string {
	return acctest.ConfigLatestAmazonLinuxHvmEbsAmi() + fmt.Sprintf(`
resource "aws_security_group" "tf_test_foo" {
  name        = "tf_test_foo-%d"
  description = "foo"

  ingress {
    protocol  = "icmp"
    from_port = -1
    to_port   = -1
    self      = true
  }
}

resource "aws_instance" "test" {
  ami             = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type   = "t2.small"
  security_groups = [aws_security_group.tf_test_foo.name]
  user_data       = "foo:-with-character's"
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`, rInt)
}

func testAccInstanceDataSourceConfig_VPCSecurityGroups(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccInstanceVPCBasicConfig(rName),
		testAccInstanceVPCSecurityGroupConfig(rName),
		`
resource "aws_instance" "test" {
  ami                    = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type          = "t2.micro"
  vpc_security_group_ids = [aws_security_group.test.id]
  subnet_id              = aws_subnet.test.id
  depends_on             = [aws_internet_gateway.test]
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`)
}

func testAccInstanceDataSourceConfig_getPasswordData(rName, publicKey string, val bool) string {
	return acctest.ConfigCompose(
		testAccLatestWindowsServer2016CoreAMIConfig(),
		fmt.Sprintf(`
resource "aws_key_pair" "test" {
  key_name   = %[1]q
  public_key = %[2]q
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.win2016core-ami.id
  instance_type = "t2.medium"
  key_name      = aws_key_pair.test.key_name
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id

  get_password_data = %[3]t
}
`, rName, publicKey, val))
}

func testAccInstanceDataSourceConfigGetUserData(rName string, getUserData bool) string {
	return acctest.ConfigLatestAmazonLinuxHvmEbsAmi() + testAccInstanceVPCBasicConfig(rName) + fmt.Sprintf(`
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.micro"
  subnet_id     = aws_subnet.test.id

  user_data = <<EUD
#!/bin/bash

echo "hello world"
EUD
}

data "aws_instance" "test" {
  get_user_data = %[2]t
  instance_id   = aws_instance.test.id
}
`, rName, getUserData)
}

func testAccInstanceDataSourceConfigGetUserDataNoUserData(rName string, getUserData bool) string {
	return acctest.ConfigLatestAmazonLinuxHvmEbsAmi() + testAccInstanceVPCBasicConfig(rName) + fmt.Sprintf(`
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.micro"
  subnet_id     = aws_subnet.test.id
}

data "aws_instance" "test" {
  get_user_data = %[2]t
  instance_id   = aws_instance.test.id
}
`, rName, getUserData)
}

func testAccInstanceDataSourceConfig_creditSpecification(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccInstanceVPCBasicConfig(rName), `
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.micro"
  subnet_id     = aws_subnet.test.id

  credit_specification {
    cpu_credits = "unlimited"
  }
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`)
}

func testAccInstanceDataSourceConfig_metadataOptions(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccInstanceVPCConfig(rName, false),
		acctest.AvailableEC2InstanceTypeForRegion("t3.micro", "t2.micro"),
		fmt.Sprintf(`
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type
  subnet_id     = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }

  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 2
    instance_metadata_tags      = "enabled"
  }
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`, rName))
}

func testAccInstanceDataSourceConfig_enclaveOptions(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccInstanceVPCConfig(rName, false),
		acctest.AvailableEC2InstanceTypeForRegion("c5a.xlarge", "c5.xlarge"),
		fmt.Sprintf(`
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type
  subnet_id     = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }

  enclave_options {
    enabled = true
  }
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`, rName))
}

func testAccInstanceDataSourceConfig_blockDeviceTags(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		acctest.AvailableEC2InstanceTypeForRegion("t3.micro", "t2.micro"),
		fmt.Sprintf(`
resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type

  tags = {
    Name = %[1]q
  }

  ebs_block_device {
    device_name = "/dev/xvdc"
    volume_size = 10

    tags = {
      Name   = %[1]q
      Factum = "SapereAude"
    }
  }

  root_block_device {
    tags = {
      Name   = %[1]q
      Factum = "VincitQuiSeVincit"
    }
  }
}

data "aws_instance" "test" {
  instance_id = aws_instance.test.id
}
`, rName))
}
