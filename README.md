# terraform-provider-kong
[![CircleCI](https://circleci.com/gh/alexashley/terraform-provider-kong/tree/master.svg?style=svg)](https://circleci.com/gh/alexashley/terraform-provider-kong/tree/master)

A Terraform provider for the api gateway [Kong](https://github.com/Kong/kong).

## Features

- Supports Kong's enterprise-edition RBAC authentication
- Resources for individual plugins, including EE-only plugins 
- Import individual consumers, services, routes, or plugins
- [WIP] Bulk import tool to ease migrating existing infrastructure into Terraform 

## CE Resources

### Provider
```hcl
provider "kong" {
  admin_api_url = "http://localhost:8001"
  rbac_token = "foobar" // Only available with the enterprise edition
}
```

### Services (`kong_service`)
A representation of Kong's [service object](https://docs.konghq.com/0.14.x/admin-api/#service-object).

#### Schema
| field             | explanation                                                                          | type   | default |
|-------------------|--------------------------------------------------------------------------------------|--------|---------|
| `name`            | The service name.                                                                    | `string` | N/A     |
| `url`             | The url for the service. It encapsulates protocol, host, port and path.              | `string` | N/A     |
| `connect_timeout` | Time in milliseconds to connect to the upstream server.                              | `int`    | 60000   |
| `read_timeout`    | Time in milliseconds between two  read operations to the upstream server.            | `int`    | 60000   |
| `write_timeout`   | Time in milliseconds between two successive write operations to the upstream server. | `int`    | 60000   |
| `retries`         | Number of times Kong will try to proxy if there's an error.                        | `int`    | 5       |

#### Example
```hcl
resource "kong_service" "mockbin" {
  name = "mockbin"
  url = "https://mockbin.org/request"
}
```

#### Import
Existing Kong services can be imported:

`terraform import kong_service.name-of-service-to-import <service UUID>`

### Routes (`kong_route`)
A representation of Kong's [route object](https://docs.konghq.com/0.14.x/admin-api/#route-object).
Services can have many routes, but a route corresponds to just one service.

For more information on `regex_priority`, see the [Kong docs](https://docs.konghq.com/0.14.x/proxy/#evaluation-order).

#### Schema
| field            | explanation                                                                                                                                    | type     | default             |
|------------------|------------------------------------------------------------------------------------------------------------------------------------------------|----------|---------------------|
| `protocols`      | Protocols that Kong will proxy to this route                                                                                                   | `[]string` | `["http", "https"]` |
| `methods`        | HTTP methods that Kong will proxy to this route                                                                                                | `[]string` | `[]`                |
| `paths`          | List of path prefixes that will match this route                                                                                               | `[]string` | `[]`                |
| `strip_path`     | If the route is matched by path, this flag indicates whether the matched path should be removed from the upstream request.                     | `bool`     | true                |
| `preserve_host`  | If the route is matched by the `Host` header, this flag indicates if the `Host` header should be set to the matched value.                     | `bool`     | false               |
| `regex_priority` | Determines the order that paths defined by regexs are evaluated. Larger numbers are evaluated first, but simple path prefixes take precedence. | `int`      | 0                   |
| `service_id`     | The associated service.                                                                                                                        | `string`   | N/A                 |
#### Example
```hcl
resource "kong_route" "mock" {
  service_id = "${kong_service.mockbin.id}"
  paths = ["/mock"]
}
```
#### Import
Existing Kong routes can be imported into Terraform:

`terraform import kong_route.name-of-route-to-import <route UUID>`

### Consumers (`kong_consumer`)
A representation of Kong's [consumer object](https://docs.konghq.com/0.14.x/admin-api/#consumer-object).

#### Schema
| field       | explanation                                                                                     | type     | default |
|-------------|-------------------------------------------------------------------------------------------------|----------|---------|
| `username`  | A unique username for the consumer. Either the username or the custom_id (or both) must be set. | `string` | N/A     |
| `custom_id` | A unique identifier for the consumer.                                                           | `string` | N/A     |

#### Example

```hcl
resource "kong_consumer" "crocodile-hunter" {
  username = "steve-irwin"
}
```

#### Import 
Existing Kong consumers can also be imported into Terraform state:

`terraform import kong_consumer.crocodile-hunter <consumer UUID>`

### Plugins (`kong_plugin`)
A resource for Kong's [plugin object](https://docs.konghq.com/0.14.x/admin-api/#plugin-object).
Certain plugins are available as resources and should be used instead of this generic resource, as they provide in-depth validation, as well as greater type safety and ease of use.  If there's a specific resource for the plugin you're adding, the provider will require you to use that resource instead of the generic plugin.

Plugins can run for a service, route, consumer or globally; if multiple of the same plugin are configured, only one will run for a request, generally the most specific configuration. There's more information on [plugin precedence](https://docs.konghq.com/0.14.x/admin-api/#precedence) in the docs.

#### Schema
| field         | explanation                                                                                                                                                                                                                                  | type                | default |
|---------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------|---------|
| `service_id`  | The service for which the plugin will run.                                                                                                                                                                                                   | `string`            | N/A     |
| `route_id`    | The route for which the plugin will run.                                                                                                                                                                                                     | `string`            | N/A     |
| `consumer_id` | The consumer for which the plugin will run. Not supported by all plugins.                                                                                                                                                                    | `string`            | N/A     |
| `name`        | The plugin name, e.g. `basic-auth`                                                                                                                                                                                                           | `string`            | N/A     |
| `config`      | An object representing the plugin's configuration.  At this time it's not possible to represent all valid plugin configurations with Terraform. Should this be a problem, you can use a specific plugin resource or the `config_json` field. | `map[string]string` | N/A     |
| `config_json` | A JSON string containing the plugin's configuration. Can't be used with `config`                                                                                                                                                             | `string`            | N/A     |
| `enabled`     | Turns the plugin on or off.                                                                                                                                                                                                                  | `bool`              | `true`   |

Note that if you use the `config_json` field you'll need to provide the full plugin configuration, including defaults from Kong; otherwise, Terraform will detect that there changes that need to be applied.

#### Example

```hcl
resource "kong_plugin" "basic-auth-plugin" {
  route_id = "${kong_route.foo-route.id}"
  name = "basic-auth"
  config_json = <<EOF
{
  "anonymous": "",
  "hide_credentials": true
}
EOF
}
```

#### Import
Existing plugins can be imported:
`terraform import kong_plugin.basic-auth-plugin <plugin UUID>`


## EE Resources

### Request Transformer Advanced (`kong_plugin_request_transformer_advanced`)
A resource for the [`request-transformer-advanced`](https://docs.konghq.com/hub/kong-inc/request-transformer-advanced/) plugin.

#### Schema

In addition to the fields below, the resource also shares most of the `kong_plugin` config, with the exception of the `config` and `config_json` fields.

| field                 | explanation                                                                                                                       | type       | default |
|-----------------------|-----------------------------------------------------------------------------------------------------------------------------------|------------|---------|
| `http_method`         | Method that will be used for the upstream request.                                                                                | `string`   | N/A     |
| `remove_headers`      | Header key:value pairs to scrub from the request                                                                                  | `[]string` | N/A     |
| `remove_querystring`  | Querystring key:value pairs to scrub from the request                                                                             | `[]string` | N/A     |
| `remove_body_params`  | Body parameters to scrub from the request.                                                                                        | `[]string` | N/A     |
| `replace_headers`     | Header key:value pairs. If the header is set, its value will be replaced. Otherwise it will be ignored                            | `[]string` | N/A     |
| `replace_uri`         | Rewrites the path to the upstream request.                                                                                        | `string`   | N/A     |
| `replace_body_params` | Body parameters to replace in the request. If the param is set, its value will be replaced. Otherwise it will be ignored.         | `[]string` | `true   |
| `rename_headers`      | Header key:value pairs. If the header is set, it will be renamed. The value will remain unchanged.                                | `[]string` | N/A     |
| `rename_querystring`  | Querystring key:value pairs.  If the querystring is in the request, the field will be renamed but the value will remain the same. | `[]string` | N/A     |
| `rename_body_params`  | Body parameters to rename in the request.                                                                                         | `[]string` | N/A     |
| `add_querystring`     | Querystring key:value pairs to add to the request. Ignored if the query is already set.                                           | `[]string` | N/A     |
| `add_headers`         | Header key:value pairs to add to the request. Ignored if the header is already set.                                               | `[]string` | N/A     |
| `add_body_params`     | Body parameters to add to the request. Ignored if already set.                                                                    | `[]string` | N/A     |
| `append_headers`      | Header key:value pairs to append to the request. The header is added if it's not already present                                  | `[]string` | N/A     |
| `append_querystring` | Querystring key:value pairs to append to the request. The query is added if it's not already present | `[]string` | N/A
| `append_body_params` | Body parameters to append to the request. The parameter is set if it's not already in the request | `[]string` | N/A

#### Example
```hcl
resource "kong_plugin_request_transformer_advanced" "request-transformer-plugin-service" {
  service_id = "${kong_service.mockbin.id}"
  add_headers = ["x-parent-resource:service"]
  http_method = "GET",
  replace_uri = "/foobar"
}
```

#### Import

To import an existing instance of the plugin:
`terraform import kong_plugin_request_transformer_advanced.req-transformer <plugin UUID>`

### Admins/RBAC

TBD

### Unsupported Resources

No plans to support the legacy [API resource](https://docs.konghq.com/0.12.x/admin-api/#api-object).

I also don't intend to add support for the following resources, since I don't currently use them:

- [Certificate](https://docs.konghq.com/0.14.x/admin-api/#certificate-object)
- [SNI](https://docs.konghq.com/0.14.x/admin-api/#sni-objects)
- [Upstream](https://docs.konghq.com/0.14.x/admin-api/#upstream-objects)
- [Target](https://docs.konghq.com/0.14.x/admin-api/#target-object)

That said, PRs are welcome and may I try to implement them once the other resources are mature.

## Development

- install Go 1.11 (this project uses [Go modules](https://github.com/golang/go/wiki/Modules#installing-and-activating-module-support))
- install Terraform.
- `docker-compose up --build -d` to stand up an instance of Kong 0.14 CE w/ Postgres.
- `make build` to create the provider.
- `make example` to run the provider against Kong.

### Manually testing imported resources

- Run `./examples/import/create-resources.sh` to create an example service, route, and a few plugins.
- Then `terraform import -config=examples/import <resource>.<resource-name> <UUID>` should import the resource (ex: `terraform import -config=examples/import kong_service.service-to-import e86f981e-a580-4bd6-aef3-1324adfcc12c`).
- Afterwards `terraform destroy examples/import` will remove the resources.
