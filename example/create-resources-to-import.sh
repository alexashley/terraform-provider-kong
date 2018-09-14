#!/usr/bin/env bash
#!/usr/bin/env bash

set -eu

KONG_API="http://localhost:8001"

SERVICE_NAME="service-to-import"

service=$(curl -s -X POST ${KONG_API}/services \
            -d "name=${SERVICE_NAME}" \
            -d "url=http://mockbin.org/request")

service_id=$(jq -n ${service} | jq -r '.id')

echo "Made service with id: ${service_id}"

route=$(curl -s -X POST ${KONG_API}/routes/ \
            -d 'paths[]=/route-to-import' \
            -d "service.id=${service_id}")

route_id=$(jq -n ${route} | jq -r '.id')

echo "Made route with id ${route_id}"

basic_auth_plugin_id=$(curl -s -X POST ${KONG_API}/plugins \
  -d 'name=basic-auth' \
  -d "route_id=${route_id}" | jq -r '.id')

echo "Made basic-auth plugin with id ${basic_auth_plugin_id}"

ip_header_restriction_plugin_id=$(curl -s -X POST ${KONG_API}/plugins \
  -d 'name=ip-header-restriction' \
  -d "route_id=$(jq -n ${route} | jq -r '.id')" \
  -d 'config.whitelist=["244.213.135.114"]' | jq -r '.id')

echo "Made ip-header-restriction plugin with id ${ip_header_restriction_plugin_id}"
