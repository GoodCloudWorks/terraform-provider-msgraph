package data

import (
	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/acceptance"
	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/resources"
)

var protoV6ProviderFactories = acceptance.NewProtoV6ProviderFactories(DataSources(), resources.Resources())
