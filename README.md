# terraform-provider-kong

A Terraform provider for the api gateway [Kong](https://github.com/Kong/kong).

## development

- install Go 1.11 (this project uses [Go modules](https://github.com/golang/go/wiki/Modules#installing-and-activating-module-support))
- install Terraform.
- `docker-compose up --build -d` to stand up an instance of Kong 0.14 CE w/ Postgres.
- `make build` to create the provider.
- `make tf-example` to run the provider against Kong.



