# [kong_plugin_openid_connect](https://github.com/alexashley/terraform-provider-kong/tree/master/kong/provider/resource_kong_plugin_openid_connect.go)
A resource for the Kong Enterprise plugin [`openid-connect`](https://docs.konghq.com/hub/kong-inc/openid-connect/).
Due to the complexity of the plugin, only a subset of the functionality is currently exposed; for example, it cannot be configured as a relying party.

### Example usage

~~~hcl
resource "kong_plugin_openid_connect" "oidc-route" {
  route_id      = "${kong_route.mock.id}"
  auth_methods  = ["bearer"]
  issuer        = "https://oidc.example.com/auth/"
}
~~~

~~~js
const a = async () => Promise.resolve('a');

console.log(await a());
~~~

### Fields Reference
The following fields are supported:


| field     | explanation     | type      | default     | required                         |
| :-------- | :-------------- | :-------- | :---------- | :------------------------------- |
|`issuer`|URL of the OpenId Connect server |`string`| | Y|
|`anonymous`|Anonymous consumer id. This is useful if you need to enable multiple auth plugins -- failing to authenticate will cause this consumer to be set. |`string`| | N|
|`auth_methods`|Allowed authentication methods |`set[string]`| | N|
|`consumer_claim`|JWT claims to use to map to a Kong consumer. Typically set to `sub` |`set[string]`| | N|
|`enabled`|Toggle whether the plugin will run |`bool`| true| N|
|`route_id`|Unique identifier of the associated route. |`string`| | N|
|`service_id`|Unique identifier of the associated service. |`string`| | N|


### Computed Fields
The following computed attributes are also available:

| field     | explanation     | type    |
|-----------|-----------------|---------|
|`created_at`|Unix timestamp representing when the plugin was created. |int|

### Import
Existing plugins can be imported: `terraform import kong_plugin_openid_connect <plugin UUID>`

[GitHub](https://github.com/alexashley/terraform-provider-kong)
