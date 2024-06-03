package provider

import (
	"context"
	"fmt"

	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/client"
	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/credentials"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type msGraphProviderClient struct {
	scopes     []string
	resty      *resty.Client
	credential azcore.TokenCredential
}

var _ client.MsGraphClient = &msGraphProviderClient{}

func (client *msGraphProviderClient) GetToken(context context.Context) (string, error) {
	token, err := client.credential.GetToken(context, policy.TokenRequestOptions{
		Scopes: client.scopes,
	})
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func (client *msGraphProviderClient) R(context context.Context, apiVersion types.String) *resty.Request {
	request := client.resty.R().SetContext(context)
	if !apiVersion.IsNull() {
		request.SetPathParam("api_version", apiVersion.ValueString())
	}
	return request
}

func (*msGraphProviderClient) URL(id types.String) string {
	return fmt.Sprintf("{api_version}/%s", id.ValueString())
}

func (data *MsGraphProviderData) NewClient() (*msGraphProviderClient, error) {
	credentialOptions := &credentials.CredentialOptions{
		TenantID: data.TenantID.ValueString(),
		ClientID: data.ClientID.ValueString(),

		UseOIDC: data.UseOIDC.ValueBool(),
		UseMSI:  data.UseMSI.ValueBool(),
		UseCLI:  data.UseCLI.ValueBool(),

		OIDCRequestToken:  data.OIDCRequestToken.ValueString(),
		OIDCRequestURL:    data.OIDCRequestURL.ValueString(),
		OIDCToken:         data.OIDCToken.ValueString(),
		OIDCTokenFilePath: data.OIDCTokenFilePath.ValueString(),
	}

	credential, err := credentials.NewCredential(credentialOptions)
	if err != nil {
		return nil, err
	}

	var scopes []string
	for _, scope := range data.Scopes.Elements() {
		scopes = append(scopes, scope.(types.String).ValueString())
	}

	client := resty.New()
	client.BaseURL = "https://graph.microsoft.com/"
	client.SetPathParam("api_version", data.ApiVersion.ValueString())
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		token, err := credential.GetToken(req.Context(), policy.TokenRequestOptions{
			Scopes: scopes,
		})
		if err != nil {
			return err
		}
		req.SetAuthScheme("Bearer").SetAuthToken(token.Token)
		return nil
	})

	providerClient := &msGraphProviderClient{
		scopes:     scopes,
		resty:      client,
		credential: credential,
	}

	return providerClient, nil
}
