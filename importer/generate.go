package main

import (
	"context"
	"flag"
	"fmt"
	kong "github.com/alexashley/terraform-provider-kong/kong/client"
	"github.com/google/subcommands"
)

type generateCommand struct {
	adminApiUrl string
	rbacToken   string
}

func (*generateCommand) Name() string {
	return "generate"
}

func (*generateCommand) Synopsis() string {
	return "Generate HCL from a Kong instance"
}

func (*generateCommand) Usage() string {
	return `generate -admin-api-url=https://kong-admin.foo.com`
}

func (cmd *generateCommand) SetFlags(flags *flag.FlagSet) {
	flags.StringVar(
		&cmd.adminApiUrl,
		"admin-api-url",
		"http://localhost:8001",
		"Kong's admin api url. Usually listening on port 8001.",
	)
	flags.StringVar(
		&cmd.rbacToken,
		"rbac-token",
		"",
		"Kong EE RBAC token. Only necessary if your Kong EE installation is secured with RBAC.",
	)
}

//type Hcl string
//type pluginAdapter func(plugin *kong.KongPlugin) Hcl
//
//var pluginAdapters = map[string] pluginAdapter{
//	"ip-header-restriction": func(plugin *kong.KongPlugin) Hcl {
//		return "foo"
//	},
//}

func (cmd *generateCommand) Execute(_ context.Context, flags *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	fmt.Println("Generating Terraform config for: " + cmd.adminApiUrl)

	c, _ := kong.NewKongClient(kong.KongConfig{
		AdminApiUrl: cmd.adminApiUrl,
		RbacToken:   cmd.rbacToken,
	})

	return subcommands.ExitSuccess
}
