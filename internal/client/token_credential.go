package client

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func NewTokenCredential() (*azidentity.ChainedTokenCredential, error) {
	var credentials []azcore.TokenCredential

	//options := &azidentity.DefaultAzureCredentialOptions{}

	// environmentCredential, err := azidentity.NewEnvironmentCredential(&azidentity.EnvironmentCredentialOptions{
	// 	ClientOptions:            options.ClientOptions,
	// 	DisableInstanceDiscovery: options.DisableInstanceDiscovery,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// credentials = append(credentials, environmentCredential)

	defaultCredential, err := azidentity.NewDefaultAzureCredential(nil)
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
