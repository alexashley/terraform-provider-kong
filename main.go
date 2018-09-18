package main

import (
	"github.com/alexashley/terraform-provider-kong/kong/importer"
	"github.com/alexashley/terraform-provider-kong/kong/provider"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"os"
)

func main() {
	if os.Getenv("KONG_IMPORTER") != "" {
		importer.RunCli()
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: func() terraform.ResourceProvider {
				return provider.KongProvider()
			},
		})
	}
}
