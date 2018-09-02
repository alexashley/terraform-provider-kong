all=build tf-example

.PHONY=all
MAKEFLAGS += --silent

build:
	 GO111MODULE=on go build -o terraform-provider-kong

tf-example: build
	terraform init example
	terraform destroy example
	terraform apply example
	terraform plan example
