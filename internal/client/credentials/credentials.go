package credentials

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type CredentialOptions struct {
	TenantID string
	ClientID string

	UseOIDC           bool
	OIDCRequestToken  string
	OIDCRequestURL    string
	OIDCToken         string
	OIDCTokenFilePath string
}

func NewTokenCredential(options *CredentialOptions) (azcore.TokenCredential, error) {
	var credentials []azcore.TokenCredential

	credentialOptions := &azidentity.DefaultAzureCredentialOptions{
		TenantID: options.TenantID,
	}

	if options.UseOIDC {
		oidcCredential, err := newOidcCredential(options)
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

func newOidcCredential(options *CredentialOptions) (azcore.TokenCredential, error) {
	oidcOptions := &OidcCredentialOptions{
		ClientID: options.ClientID,
		TenantID: options.TenantID,

		RequestToken:  options.OIDCRequestToken,
		RequestUrl:    options.OIDCRequestURL,
		Token:         options.OIDCToken,
		TokenFilePath: options.OIDCTokenFilePath,
	}
	credential, err := NewOidcCredential(oidcOptions)
	if err != nil {
		return nil, err
	}
	return credential, nil
}
