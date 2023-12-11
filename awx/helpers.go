package awx

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	allAWXDataSources []func() datasource.DataSource
	allAWXResources   []func() resource.Resource
)

// registerDataSource may be called during package initialization to register a new data source with the provider.
func registerDataSource(fn func() datasource.DataSource) {
	allAWXDataSources = append(allAWXDataSources, fn)
}

// registerResource may be called during package initialization to register a new resource with the provider.
func registerResource(fn func() resource.Resource) {
	allAWXResources = append(allAWXResources, fn)
}
