package main

import (
	"github.com/alexashley/terraform-provider-kong/kong/provider"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return provider.KongProvider()
		},
	})
}
