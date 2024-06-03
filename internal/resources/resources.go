package resources

import "github.com/hashicorp/terraform-plugin-framework/resource"

func Resources() []func() resource.Resource {
	return []func() resource.Resource{}
}
