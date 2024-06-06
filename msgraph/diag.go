package msgraph

import "github.com/hashicorp/terraform-plugin-framework/diag"

func noErrors() diag.Diagnostics {
	return diag.Diagnostics{}
}

func errorDiagnostics(summary string, detail string) diag.Diagnostics {
	return diag.Diagnostics{
		diag.NewErrorDiagnostic(summary, detail),
	}
}
