package credentials

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type CredentialOptions struct {
	TenantID string
	ClientID string

	UseMSI  bool
	UseCLI  bool
	UseOIDC bool

	OIDCRequestToken  string
	OIDCRequestURL    string
	OIDCToken         string
	OIDCTokenFilePath string
}

func NewTokenCredential(options *CredentialOptions) (azcore.TokenCredential, error) {
	var credentials []azcore.TokenCredential

	if options.UseOIDC {
		oidcCredential, err := newOidcCredential(options)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, oidcCredential)
	}

	if options.UseMSI {
		msiCredential, err := newMsiCredential(options)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, msiCredential)
	}

	if options.UseCLI {
		cliCredential, err := newCliCredential(options)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, cliCredential)
	}

	envCredential, err := newEnvironmentCredential()
	if err != nil {
		return nil, err
	}
	credentials = append(credentials, envCredential)

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

func newMsiCredential(options *CredentialOptions) (azcore.TokenCredential, error) {
	msiOptions := &azidentity.ManagedIdentityCredentialOptions{}

	if options.ClientID != "" {
		msiOptions.ID = azidentity.ClientID(options.ClientID)
	}

	return azidentity.NewManagedIdentityCredential(msiOptions)

}

func newCliCredential(options *CredentialOptions) (azcore.TokenCredential, error) {
	cliOptions := &azidentity.AzureCLICredentialOptions{
		TenantID: options.TenantID,
	}
	return azidentity.NewAzureCLICredential(cliOptions)
}

func newEnvironmentCredential() (azcore.TokenCredential, error) {
	envOptions := &azidentity.EnvironmentCredentialOptions{}
	return azidentity.NewEnvironmentCredential(envOptions)
}
