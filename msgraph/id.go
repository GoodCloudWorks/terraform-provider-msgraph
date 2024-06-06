package msgraph

import (
	"fmt"

	"github.com/GoodCloudWorks/terraform-provider-msgraph/msgraph/id"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ensureParseIDString(value types.String) (*id.ID, diag.Diagnostics) {
	return ensureParseID(value.ValueString())
}

func ensureParseID(value string) (*id.ID, diag.Diagnostics) {
	id, err := id.Parse(value)
	if err != nil {
		return nil, errorDiagnostics(fmt.Sprintf("Failed to parse ID: %q", value), err.Error())
	}
	return id, noErrors()
}
