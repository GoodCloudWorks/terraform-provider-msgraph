// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	internal_provider "terraform-provider-msgraph/internal/provider"

	terraform_provider "github.com/hashicorp/terraform-plugin-framework/provider"

	terraform_provider_server "github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name msgraph
func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := terraform_provider_server.ServeOpts{
		Address: "registry.terraform.io/goodcloudworks/msgraph",
		Debug:   debug,
	}

	err := terraform_provider_server.Serve(context.Background(), func() terraform_provider.Provider {
		return &internal_provider.MsGraphProvider{}
	}, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
