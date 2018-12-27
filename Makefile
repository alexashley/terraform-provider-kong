.PHONY: build example testacc clean-docs build-docs docs
.ONESHELL: build-docs
MAKEFLAGS += --silent

KONG ?= "http://localhost:8001"

build:
	 GO111MODULE=on go build -o terraform-provider-kong

testacc:
	TF_ACC=1 KONG_ADMIN_API_URL=${KONG} go test ./kong/provider -v -coverprofile=coverage.out -covermode=count

testapi:
	go test ./kong/kong -v

example: build
	terraform init example
	terraform destroy example
	terraform apply example
	terraform plan example

clean-docs:
	rm -rf docs/*.md

build-docs:
	cd docsgen
	GO111MODULE=on go build -o docsgen

docs: clean-docs build-docs
	./docsgen/docsgen

coverage-report: testacc
	 go tool cover -html=coverage.out
