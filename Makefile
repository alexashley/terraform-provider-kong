all=build example testacc

.PHONY=all
MAKEFLAGS += --silent

KONG ?= "http://localhost:8001"

build:
	 GO111MODULE=on go build -o terraform-provider-kong

testacc:
	TF_ACC=1 KONG_ADMIN_API_URL=${KONG} go test ./kong/provider -v

example: build
	terraform init example
	terraform destroy example
	terraform apply example
	terraform plan example
