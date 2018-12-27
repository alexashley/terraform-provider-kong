# kong_route
A representation of Kong's [route object](https://docs.konghq.com/0.14.x/admin-api/#route-object).
Services can have many routes, but a route corresponds to just one service.

For more information on `regex_priority`, see the [Kong docs](https://docs.konghq.com/0.14.x/proxy/#evaluation-order).

### Example usage

```hcl
resource "kong_route" "mock" {
  service_id = "${kong_service.mockbin.id}"
  paths = ["/mock"]
}

```

### Fields Reference
The following fields are supported:

| field     | explanation     | type      | default     | required                         |
| :-------- | :-------------- | :-------- | :---------- | :------------------------------- |
|`service_id`|Unique identifier of the associated service. |`string`| N/A| Y|
|`hosts`|Host header values that should be matched to this route. |`set[string]`| N/A| N|
|`methods`|HTTP verbs that Kong will proxy to this route |`set[string]`| N/A| N|
|`paths`|List of path prefixes that will match this route |`set[string]`| N/A| N|
|`preserve_host`| If the route is matched by the `Host` header, this flag indicates if the `Host` header should be set to the matched value. |`bool`| false| N|
|`protocols`|Protocols that Kong will proxy to this route |`set[string]`| N/A| N|
|`regex_priority`|Determines the order that paths defined by regexes are evaluated. |`int`| 0| N|
|`strip_path`| If the route is matched by path, this flag indicates whether the matched path should be removed from the upstream request. |`bool`| true| N|
### Computed Fields
The following computed attributes are also available:

| field     | explanation     | type    |
|-----------|-----------------|---------|
|`created_at`|Unix timestamp of creation date. |int|
|`updated_at`|Unix timestamp of last edit date. |int|

### Import
Existing Kong routes can be imported into Terraform:
`terraform import kong_route.name-of-route-to-import <route UUID>`
