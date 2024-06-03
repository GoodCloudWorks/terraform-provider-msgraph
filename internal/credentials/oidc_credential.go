package credentials

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/go-resty/resty/v2"
)

type oidcCredential struct {
	requestToken  string
	requestUrl    string
	token         string
	tokenFilePath string
	resty         *resty.Client
}

func newOidcCredential(options *CredentialOptions) (azcore.TokenCredential, error) {
	oidcCredential := &oidcCredential{
		requestToken:  options.OIDCRequestToken,
		requestUrl:    options.OIDCRequestURL,
		token:         options.OIDCToken,
		tokenFilePath: options.OIDCTokenFilePath,
		resty:         resty.New(),
	}

	credential, err := azidentity.NewClientAssertionCredential(
		options.TenantID,
		options.ClientID,
		oidcCredential.getAssertion,
		&azidentity.ClientAssertionCredentialOptions{},
	)
	if err != nil {
		return nil, err
	}

	return credential, nil
}

func (credential *oidcCredential) getAssertion(ctx context.Context) (string, error) {
	if credential.token != "" {
		return credential.token, nil
	}

	if credential.tokenFilePath != "" {
		idTokenData, err := os.ReadFile(credential.tokenFilePath)
		if err != nil {
			return "", fmt.Errorf("getAssertion: reading token file: %v", err)
		}

		return string(idTokenData), nil
	}

	var tokenResponse struct {
		Count *int    `json:"count"`
		Value *string `json:"value"`
	}

	url, err := url.Parse(credential.requestUrl)
	if err != nil {
		return "", fmt.Errorf("getAssertion: cannot parse URL: %s", err)
	}

	if url.Query().Get("audience") == "" {
		query := url.Query()
		query.Set("audience", "api://AzureADTokenExchange")
		url.RawQuery = query.Encode()
	}

	response, err := credential.resty.R().
		SetContext(ctx).
		SetAuthScheme("Bearer").
		SetAuthToken(credential.requestToken).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetResult(&tokenResponse).
		Get(url.String())
	if err != nil {
		return "", fmt.Errorf("getAssertion: cannot request token: %v", err)
	}

	if c := response.StatusCode(); c < 200 || c > 299 {
		return "", fmt.Errorf("getAssertion: received HTTP status %d with response: %s", c, string(response.Body()))
	}

	if tokenResponse.Value == nil {
		return "", fmt.Errorf("getAssertion: nil JWT assertion received from OIDC provider")
	}

	return *tokenResponse.Value, nil
}
