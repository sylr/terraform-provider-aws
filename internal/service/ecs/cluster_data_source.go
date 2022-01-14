package ecs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceClusterRead,

		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"pending_tasks_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"running_tasks_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"registered_container_instances_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"setting": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ECSConn

	clusterName := d.Get("cluster_name").(string)
	desc, err := FindClusterByNameOrARN(conn, d.Get("cluster_name").(string))

	if err != nil {
		return fmt.Errorf("error reading ECS Cluster (%s): %w", clusterName, err)
	}

	if len(desc.Clusters) == 0 {
		return fmt.Errorf("no matches found for name: %s", clusterName)
	}

	if len(desc.Clusters) > 1 {
		return fmt.Errorf("multiple matches found for name: %s", clusterName)
	}

	cluster := desc.Clusters[0]
	d.SetId(aws.StringValue(cluster.ClusterArn))
	d.Set("arn", cluster.ClusterArn)
	d.Set("status", cluster.Status)
	d.Set("pending_tasks_count", cluster.PendingTasksCount)
	d.Set("running_tasks_count", cluster.RunningTasksCount)
	d.Set("registered_container_instances_count", cluster.RegisteredContainerInstancesCount)

	if err := d.Set("setting", flattenClusterSettings(cluster.Settings)); err != nil {
		return fmt.Errorf("error setting setting: %w", err)
	}

	return nil
}
