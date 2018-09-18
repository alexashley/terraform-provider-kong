# kong-importer

A command-line utility to import Kong resources into Terraform.

## usage

### Creating Terraform config
`./kong-import generate --api-url=http://localhost:8001`


### Bulk import
`./kong-import run --plan=import.json --tf-dir=foo`
