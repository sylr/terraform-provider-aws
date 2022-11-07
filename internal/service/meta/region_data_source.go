package meta

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func init() {
	registerFrameworkDataSourceFactory(newDataSourceRegion)
}

// newDataSourceRegion instantiates a new DataSource for the aws_region data source.
func newDataSourceRegion(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &dataSourceRegion{}, nil
}

type dataSourceRegion struct {
	meta *conns.AWSClient
}

// Metadata should return the full name of the data source, such as
// examplecloud_thing.
func (d *dataSourceRegion) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) { // nosemgrep:ci.meta-in-func-name
	response.TypeName = "aws_region"
}

// GetSchema returns the schema for this data source.
func (d *dataSourceRegion) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	schema := tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"description": {
				Type:     types.StringType,
				Computed: true,
			},
			"endpoint": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
		},
	}

	return schema, nil
}

// Configure enables provider-level data or clients to be set in the
// provider-defined DataSource type. It is separately executed for each
// ReadDataSource RPC.
func (d *dataSourceRegion) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*conns.AWSClient); ok {
		d.meta = v
	}
}

// Read is called when the provider must read data source values in order to update state.
// Config values should be read from the ReadRequest and new state values set on the ReadResponse.
func (d *dataSourceRegion) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data dataSourceRegionData

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	var region *endpoints.Region

	if !data.Endpoint.IsNull() {
		matchingRegion, err := FindRegionByEndpoint(data.Endpoint.Value)

		if err != nil {
			response.Diagnostics.AddError("finding Region by endpoint", err.Error())

			return
		}

		region = matchingRegion
	}

	if !data.Name.IsNull() {
		matchingRegion, err := FindRegionByName(data.Name.Value)

		if err != nil {
			response.Diagnostics.AddError("finding Region by name", err.Error())

			return
		}

		if region != nil && region.ID() != matchingRegion.ID() {
			response.Diagnostics.AddError("multiple Regions matched", "use additional constraints to reduce matches to a single Region")

			return
		}

		region = matchingRegion
	}

	// Default to provider current region if no other filters matched
	if region == nil {
		matchingRegion, err := FindRegionByName(d.meta.Region)

		if err != nil {
			response.Diagnostics.AddError("finding Region by name", err.Error())

			return
		}

		region = matchingRegion
	}

	regionEndpointEC2, err := region.ResolveEndpoint(endpoints.Ec2ServiceID)

	if err != nil {
		response.Diagnostics.AddError("resolving EC2 endpoint", err.Error())

		return
	}

	data.Description = types.String{Value: region.Description()}
	data.Endpoint = types.String{Value: strings.TrimPrefix(regionEndpointEC2.URL, "https://")}
	data.ID = types.String{Value: region.ID()}
	data.Name = types.String{Value: region.ID()}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

type dataSourceRegionData struct {
	Description types.String `tfsdk:"description"`
	Endpoint    types.String `tfsdk:"endpoint"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
}

func FindRegionByEndpoint(endpoint string) (*endpoints.Region, error) {
	for _, partition := range endpoints.DefaultPartitions() {
		for _, region := range partition.Regions() {
			regionEndpointEC2, err := region.ResolveEndpoint(endpoints.Ec2ServiceID)

			if err != nil {
				return nil, err
			}

			if strings.TrimPrefix(regionEndpointEC2.URL, "https://") == endpoint {
				return &region, nil
			}
		}
	}

	return nil, fmt.Errorf("region not found for endpoint %q", endpoint)
}

func FindRegionByName(name string) (*endpoints.Region, error) {
	for _, partition := range endpoints.DefaultPartitions() {
		for _, region := range partition.Regions() {
			if region.ID() == name {
				return &region, nil
			}
		}
	}

	return nil, fmt.Errorf("region not found for name %q", name)
}
