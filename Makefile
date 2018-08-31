all=build tf-example

.PHONY=all
MAKEFLAGS += --silent

build:
	 GO111MODULE=on go build -o terraform-provider-kong

tf-example: build
	terraform destroy example
	terraform init example
	terraform apply example
