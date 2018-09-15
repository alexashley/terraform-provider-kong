# terraform-provider-kong
[![CircleCI](https://circleci.com/gh/alexashley/terraform-provider-kong/tree/master.svg?style=svg)](https://circleci.com/gh/alexashley/terraform-provider-kong/tree/master)

A Terraform provider for the api gateway [Kong](https://github.com/Kong/kong).

## Resources

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

### Plugins (`kong_plugin`)

### Consumers


### Admins/RBAC

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
