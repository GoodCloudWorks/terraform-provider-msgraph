package client

import (
	"context"
	"encoding/json"
	"fmt"

	"terraform-provider-msgraph/internal/client/credentials"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/go-resty/resty/v2"
)

type MsGraphClientOptions struct {
	ApiVersion  string
	Credentials *credentials.CredentialOptions
}

type MsGraphClient struct {
	Options    *MsGraphClientOptions
	credential azcore.TokenCredential
	resty      *resty.Client
}

func NewMsGraphClient(options *MsGraphClientOptions) (*MsGraphClient, error) {
	credential, err := credentials.NewTokenCredential(options.Credentials)
	if err != nil {
		return nil, err
	}

	client := resty.New()
	client.BaseURL = "https://graph.microsoft.com/"
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		token, err := credential.GetToken(req.Context(), policy.TokenRequestOptions{
			Scopes: []string{"https://graph.microsoft.com/.default"},
		})
		if err != nil {
			return err
		}
		req.SetAuthScheme("Bearer").SetAuthToken(token.Token)
		return nil
	})
	msGraphClient := &MsGraphClient{
		Options:    options,
		credential: credential,
		resty:      client,
	}
	return msGraphClient, nil
}

func (client *MsGraphClient) GetAccessToken(context context.Context) (string, error) {
	token, err := client.credential.GetToken(context, policy.TokenRequestOptions{
		Scopes: []string{"https://graph.microsoft.com/.default"},
	})
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func (client *MsGraphClient) Get(context context.Context, path string) (*interface{}, error) {
	response, err := client.resty.R().SetContext(context).Get(path)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() == 404 {
		return nil, fmt.Errorf("resource %q not found", path)
	}

	if response.StatusCode() == 403 {
		return nil, fmt.Errorf("access denied to resource %q", path)
	}

	if response.IsError() {
		return nil, fmt.Errorf("failed to get resource %q: %d %s", path, response.StatusCode(), string(response.Body()))
	}

	var result interface{}
	err = json.Unmarshal(response.Body(), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
