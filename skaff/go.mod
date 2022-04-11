module github.com/hashicorp/terraform-provider-aws/skaff

go 1.17

require (
	github.com/hashicorp/terraform-provider-aws v1.60.1-0.20220322001452-8f7a597d0c24
	github.com/spf13/cobra v1.4.0
)

require (
	github.com/aws/aws-sdk-go v1.43.26 // indirect
	github.com/aws/aws-sdk-go-v2 v1.16.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.8 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53domains v1.12.2 // indirect
	github.com/aws/smithy-go v1.11.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/hashicorp/terraform-provider-aws => ../
