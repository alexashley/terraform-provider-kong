all=build tf-example

.PHONY=all
MAKEFLAGS += --silent

build:
	 GO111MODULE=on go build -o terraform-provider-kong

example: build
	terraform init examples/simple
	terraform destroy examples/simple
	terraform apply examples/simple
	terraform plan examples/simple

import: build
	terraform init examples/import
	terraform import examples/import
	terraform plan examples/import
