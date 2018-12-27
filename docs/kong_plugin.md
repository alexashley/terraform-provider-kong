# [kong_plugin](https://github.com/alexashley/terraform-provider-kong/tree/master/kong/provider/resource_kong_plugin.go)
A resource for Kong's [plugin object](https://docs.konghq.com/0.14.x/admin-api/#plugin-object).
In addition to this generic resource, some plugins have resources of their own.
It's recommended to use these custom resources as they provide in-depth validation, as well as greater type safety and ease of use. If there's a specific resource for the plugin you're adding, the provider will require you to use that resource instead of the generic plugin.

Plugins can run for a service, route, consumer or globally; if multiple of the same plugin are configured, only one will run for a request, generally the most specific configuration. There's more information on [plugin precedence](https://docs.konghq.com/0.14.x/admin-api/#precedence) in the docs.

Note that if you use the `config_json` field you'll need to provide the full plugin configuration, including defaults from Kong; otherwise, Terraform will detect that there changes that need to be applied.

### Example usage

~~~
resource "kong_plugin" "basic-auth-plugin" {
  route_id      = "${kong_route.foo-route.id}"
  name          = "basic-auth"
  config_json   = <<EOF
  {
    "anonymous": "",
    "hide_credentials": true
  }
  EOF
}
~~~

### Fields Reference
The following fields are supported:


| field     | explanation     | type      | default     | required                         |
| :-------- | :-------------- | :-------- | :---------- | :------------------------------- |
|`config_json`|A JSON string containing the plugin's configuration. |`string`| | Y|
|`name`|The plugin name, e.g. `basic-auth` |`string`| | Y|
|`consumer_id`|Unique identifier of the consumer for which this plugin will run. Not all plugins allow consumers |`string`| | N|
|`enabled`|Toggle whether the plugin will run |`bool`| true| N|
|`route_id`|Unique identifier of the associated route. |`string`| | N|
|`service_id`|Unique identifier of the associated service. |`string`| | N|


### Computed Fields
The following computed attributes are also available:

| field     | explanation     | type    |
|-----------|-----------------|---------|
|`created_at`|Unix timestamp representing when the plugin was created. |int|

### Import
Existing plugins can be imported: `terraform import kong_plugin.basic-auth-plugin <plugin UUID>`

[GitHub](https://github.com/alexashley/terraform-provider-kong)
