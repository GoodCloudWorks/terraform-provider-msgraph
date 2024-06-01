package provider

import (
	"context"
	"fmt"
	"terraform-provider-msgraph/internal/credentials"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MsGraphClient interface {
	GetToken(context context.Context) (string, error)
	R(context context.Context, apiVersion types.String) *resty.Request
	URL(id types.String) string
}

type MsGraphProviderClient struct {
	scopes     []string
	resty      *resty.Client
	credential azcore.TokenCredential
}

var _ MsGraphClient = &MsGraphProviderClient{}

func (client *MsGraphProviderClient) GetToken(context context.Context) (string, error) {
	token, err := client.credential.GetToken(context, policy.TokenRequestOptions{
		Scopes: client.scopes,
	})
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func (client *MsGraphProviderClient) R(context context.Context, apiVersion types.String) *resty.Request {
	request := client.resty.R().SetContext(context)
	if !apiVersion.IsNull() {
		request.SetPathParam("api_version", apiVersion.ValueString())
	}
	return request
}

func (*MsGraphProviderClient) URL(id types.String) string {
	return fmt.Sprintf("{api_version}/%s", id.ValueString())
}

func (data *MsGraphProviderData) NewClient() (*MsGraphProviderClient, error) {
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

	credential, err := credentials.NewTokenCredential(credentialOptions)
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

	providerClient := &MsGraphProviderClient{
		scopes:     scopes,
		resty:      client,
		credential: credential,
	}

	return providerClient, nil
}
