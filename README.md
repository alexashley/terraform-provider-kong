# terraform-provider-kong
[![CircleCI](https://circleci.com/gh/alexashley/terraform-provider-kong/tree/master.svg?style=svg)](https://circleci.com/gh/alexashley/terraform-provider-kong/tree/master) | [Documentation](/docs/index.md)

A Terraform provider for the api gateway [Kong](https://github.com/Kong/kong).

## Features

- Supports Kong's Enterprise authentication
- Resources for individual plugins, including EE-only plugins 
- Import individual consumers, services, routes, or plugins
- [WIP] Bulk import tool to ease migrating existing infrastructure into Terraform 

## Development

- install Go 1.11 (this project uses [Go modules](https://github.com/golang/go/wiki/Modules#installing-and-activating-module-support))
- install Terraform.
- `docker-compose up --build -d` to stand up an instance of Kong 0.14 CE w/ Postgres.
- `make build` to create the provider.
- `make example` to run the provider against Kong.

### Manually testing imported resources

- Run `./examples/import/create-resources.sh` to create an example service, route, and a few plugins.
- Then `terraform import -config=examples/import <resource>.<resource-name> <UUID>` should import the resource (ex: `terraform import -config=examples/import kong_service.service-to-import e86f981e-a580-4bd6-aef3-1324adfcc12c`).
- Afterwards `terraform destroy examples/import` will remove the resources.
