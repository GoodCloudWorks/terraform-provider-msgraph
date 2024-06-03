// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	msgraph "github.com/GoodCloudWorks/terraform-provider-msgraph/internal/provider"

	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/data"
	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/resources"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name msgraph
func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/goodcloudworks/msgraph",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), func() provider.Provider {
		return msgraph.New(data.DataSources(), resources.Resources())
	}, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
