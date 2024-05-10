package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/go-resty/resty/v2"
)

type MsGraphClientOptions struct {
	ApiVersion string
	TenantID   string
	ClientID   string

	UseOIDC           bool
	OIDCRequestToken  string
	OIDCRequestURL    string
	OIDCToken         string
	OIDCTokenFilePath string
}

type MsGraphClient struct {
	resty *resty.Client
}

func NewMsGraphClient(options *MsGraphClientOptions) (*MsGraphClient, error) {
	credential, err := newTokenCredential(options)
	if err != nil {
		return nil, err
	}

	client := resty.New()
	client.BaseURL = "https://graph.microsoft.com/" + options.ApiVersion + "/"
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
	return &MsGraphClient{resty: client}, nil
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

func newTokenCredential(options *MsGraphClientOptions) (*azidentity.ChainedTokenCredential, error) {
	var credentials []azcore.TokenCredential

	credentialOptions := &azidentity.DefaultAzureCredentialOptions{
		TenantID: options.TenantID,
	}

	if options.UseOIDC {
		oidcOptions := &oidcCredentialOptions{
			ClientID: options.ClientID,
			TenantID: options.TenantID,

			RequestToken:  options.OIDCRequestToken,
			RequestUrl:    options.OIDCRequestURL,
			Token:         options.OIDCToken,
			TokenFilePath: options.OIDCTokenFilePath,
		}
		oidcCredential, err := newOidcCredential(oidcOptions)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, oidcCredential)
	}

	defaultCredential, err := azidentity.NewDefaultAzureCredential(credentialOptions)
	if err != nil {
		return nil, err
	}
	credentials = append(credentials, defaultCredential)

	chainedCredentials, err := azidentity.NewChainedTokenCredential(credentials, nil)
	if err != nil {
		return nil, err
	}

	return chainedCredentials, nil
}
