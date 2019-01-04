#!/usr/bin/env bash

set -eu

KONG_API="http://localhost:8001"

SERVICE_NAME="import"

consumer=$(curl -s -X POST ${KONG_API}/consumers \
            -d "username=steve-irwin")

consumer_id=$(jq -n ${consumer} | jq -r '.id')

echo "Made consumer with id ${consumer_id}"

service=$(curl -s -X POST ${KONG_API}/services \
            -d "name=${SERVICE_NAME}" \
            -d "url=http://mockbin.org/request")

service_id=$(jq -n ${service} | jq -r '.id')

echo "Made service with id: ${service_id}"

route=$(curl -s -X POST ${KONG_API}/routes/ \
            -d 'paths[]=/import' \
            -d "service.id=${service_id}")
route_id=$(jq -n ${route} | jq -r '.id')

echo "Made route with id ${route_id}"


basic_auth_plugin_id=$(curl -s -X POST ${KONG_API}/plugins \
  -d 'name=basic-auth' \
  -d "route_id=${route_id}" | jq -r '.id')

echo "Made basic-auth plugin with id ${basic_auth_plugin_id}"

oidc_plugin_id=$(curl -s -X POST ${KONG_API}/plugins \
  -d 'name=openid-connect' \
  -d 'config.issuer=http://foo.org' \
  -d 'config.auth_methods=bearer' \
  -d "service_id=${service_id}" | jq -r '.id')

echo "Made openid-connect plugin with id ${oidc_plugin_id}"

request_transformer_id=$(curl -s -X POST ${KONG_API}/plugins \
  -d 'name=request-transformer-advanced' | jq -r '.id')

echo "Made request-transformer-advanced plugin with id ${oidc_plugin_id}"


