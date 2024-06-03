package resources

import (
	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/acceptance"
	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/data"
)

var protoV6ProviderFactories = acceptance.NewProtoV6ProviderFactories(data.DataSources(), Resources())
