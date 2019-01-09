.PHONY: build example testacc clean-docs build-docs docs coverageacc docker-image docker-push import build-import
.ONESHELL: build-docs build-import
MAKEFLAGS += --silent

KONG ?= "http://localhost:8001"

IMAGE_VERSION="0.0.7"

build:
	 GO111MODULE=on go build -o terraform-provider-kong

testacc:
	TF_ACC=1 KONG_ADMIN_API_URL=${KONG} go test ./kong/provider -v -coverprofile=coverage.out -covermode=count

testapi:
	go test ./kong/kong -v

example: build
	terraform init examples/simple
	terraform destroy examples/simple
	terraform apply examples/simple
	terraform plan examples/simple

clean-docs:
	rm -rf docs/*.md

build-docs:
	cd docsgen
	GO111MODULE=on go build -o docsgen

docs: clean-docs build-docs
	./docsgen/docsgen

coverageacc: testacc
	 go tool cover -html=coverage.out

docker-image:
	docker build \
	-t alexashley/tf-provider-custom-kong:${IMAGE_VERSION} \
	--build-arg KONG_CUSTOM_PLUGINS=$$(find kong-docker/plugins/* -type d | sed "s|^\kong-docker/plugins/||" | paste -d, -s) \
	kong-docker

docker-push:
	docker push alexashley/tf-provider-custom-kong:${IMAGE_VERSION}

build-import:
	cd importer
	go build -o kong-import

import: build-import
	rm -f import-state.json
	terraform init examples/import
	terraform destroy examples/import
	./examples/import/create-resources-to-import.sh
	./importer/kong-import import -admin-api-url=http://localhost:8001 -tf-config=examples/import
	terraform plan examples/import
