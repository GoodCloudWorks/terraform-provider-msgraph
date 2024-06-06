package msgraph

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/GoodCloudWorks/terraform-provider-msgraph/msgraph/dynamic"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	mimeTypeApplicationJson = "application/json"

	httpStatusNotFound = 404

	apiVersionPath = "{api_version}/"
)

func get(http *resty.Request, url string) (*resty.Response, error) {
	return http.Get(apiVersionPath + url)
}

func post(http *resty.Request, url string) (*resty.Response, error) {
	return http.Post(apiVersionPath + url)
}

func patch(http *resty.Request, url string) (*resty.Response, error) {
	return http.Patch(apiVersionPath + url)
}

func delete(http *resty.Request, url string) (*resty.Response, error) {
	return http.Delete(apiVersionPath + url)
}

func ensureHttpResponseSucceeded(response *resty.Response, err error) diag.Diagnostics {
	if err != nil {
		return errorDiagnostics(fmt.Sprintf("Request failed for: %s %q", response.Request.Method, response.Request.URL), err.Error())
	}

	if response.IsError() {
		return errorDiagnostics(fmt.Sprintf("Request failed with %d for: %s %q", response.StatusCode(), response.Request.Method, response.Request.URL), string(response.Body()))
	}

	return noErrors()
}

func ensureResponseAsDynamic(response *resty.Response) (types.Dynamic, diag.Diagnostics) {
	content, err := dynamic.FromJSONImplied(response.Body())
	if err != nil {
		return types.DynamicNull(), errorDiagnostics(fmt.Sprintf("Parse request body failed for: %s %q", response.Request.Method, response.Request.URL), string(response.Body()))
	}

	return content, noErrors()
}

func ensureIsValidPath(value string) (string, diag.Diagnostics) {
	url, err := url.Parse(value)
	if err != nil {
		return "", errorDiagnostics(fmt.Sprintf("Failed to parse url: %q", value), err.Error())
	}
	return url.Path, noErrors()
}

func ensureIsValidPathString(value types.String) (string, diag.Diagnostics) {
	return ensureIsValidPath(value.ValueString())
}

func ensureRequestSetBodyFromDynamic(request *resty.Request, value types.Dynamic) diag.Diagnostics {
	body, err := dynamic.ToJSON(value)
	if err != nil {
		return errorDiagnostics("Failed to marshal request body to JSON.", err.Error())
	}

	request.
		SetHeader("Accept", mimeTypeApplicationJson).
		SetHeader("Content-Type", mimeTypeApplicationJson).
		SetBody(body)

	return noErrors()
}

func ensureResponseHasObjectID(response *resty.Response) (string, diag.Diagnostics) {
	var content struct {
		ID string `json:"id"`
	}

	body := response.Body()
	err := json.Unmarshal(body, &content)
	if err != nil {
		return "", errorDiagnostics(fmt.Sprintf("Failed to parse response body for: %s %q", response.Request.Method, response.Request.URL), string(body))
	}

	return content.ID, noErrors()
}

func ensureGetObjectAsDynamic(http *resty.Request, url string) (types.Dynamic, diag.Diagnostics) {
	response, err := get(http, url)
	diags := ensureHttpResponseSucceeded(response, err)
	if diags.HasError() {
		return types.DynamicNull(), diags
	}

	return ensureResponseAsDynamic(response)
}
