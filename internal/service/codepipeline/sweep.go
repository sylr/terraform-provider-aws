//go:build sweep
// +build sweep

package codepipeline

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_codepipeline", &resource.Sweeper{
		Name: "aws_codepipeline",
		F:    sweepPipelines,
	})
}

func sweepPipelines(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	input := &codepipeline.ListPipelinesInput{}
	conn := client.(*conns.AWSClient).CodePipelineConn
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListPipelinesPages(input, func(page *codepipeline.ListPipelinesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Pipelines {
			r := ResourcePipeline()
			d := r.Data(nil)

			d.SetId(aws.StringValue(v.Name))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Codepipeline Pipeline sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing Codepipeline Pipelines (%s): %w", region, err)
	}

	err = sweep.SweepOrchestrator(sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Codepipeline Pipelines (%s): %w", region, err)
	}

	return nil
}
