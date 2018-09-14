all=build example

.PHONY=all
MAKEFLAGS += --silent

build:
	 GO111MODULE=on go build -o terraform-provider-kong

example: build
	terraform init example
	terraform destroy example
	terraform apply example
	terraform plan example
