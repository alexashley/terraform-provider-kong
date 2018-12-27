# kong_plugin_request_transformer_advanced
A resource for the [`request-transformer-advanced`](https://docs.konghq.com/hub/kong-inc/request-transformer-advanced/) plugin.

### Example usage

```hcl
resource  "kong_plugin_request_transformer_advanced" "request-transformer-plugin-service" {
  service_id = "${kong_service.mockbin.id}"
  add_headers = ["x-parent-resource:service"]
  http_method = "GET",
  replace_uri = "/foobar"
}

```

### Fields Reference
The following fields are supported:


| field     | explanation     | type      | default     | required                         |
| :-------- | :-------------- | :-------- | :---------- | :------------------------------- |
|`add_body_params`|Body parameters to add to the request. Ignored if already set. |`set[string]`| N/A| N|
|`add_headers`|Header key:value pairs to add to the request. Ignored if the header is already set. |`set[string]`| N/A| N|
|`add_querystring`|Querystring key:value pairs to add to the request. Ignored if the query is already set. |`set[string]`| N/A| N|
|`append_body_params`|Body parameters to append to the request. The parameter is set if it's not already in the request |`set[string]`| N/A| N|
|`append_headers`|Header key:value pairs to append to the request. The header is added if it's not already present |`set[string]`| N/A| N|
|`append_querystring`|Querystring key:value pairs to append to the request. The query is added if it's not already present |`set[string]`| N/A| N|
|`consumer_id`|Unique identifier of the consumer for which this plugin will run. Not all plugins allow consumers |`string`| N/A| N|
|`enabled`|Toggle whether the plugin will run |`bool`| true| N|
|`http_method`|Method that will be used for the upstream request. |`string`| N/A| N|
|`remove_body_params`|Body parameters to scrub from the request. |`set[string]`| N/A| N|
|`remove_headers`|Header key:value pairs to scrub from the request. |`set[string]`| N/A| N|
|`remove_querystring`|Querystring key:value pairs to scrub from the request. |`set[string]`| N/A| N|
|`rename_body_params`|Body parameters to rename in the request. |`set[string]`| N/A| N|
|`rename_headers`|Header key:value pairs. If the header is set, it will be renamed. The value will remain unchanged. |`set[string]`| N/A| N|
|`rename_querystring`|Querystring key:value pairs. If the querystring is in the request, the field will be renamed but the value will remain the same. |`set[string]`| N/A| N|
|`replace_body_params`|Body parameters to replace in the request. If the param is set, its value will be replaced. Otherwise it will be ignored. |`set[string]`| N/A| N|
|`replace_headers`|Header key:value pairs. If the header is set, its value will be replaced. Otherwise it will be ignored |`set[string]`| N/A| N|
|`replace_querystring`|Querystring key:value pairs to replace if the key is set in the request. |`set[string]`| N/A| N|
|`replace_uri`|Rewrites the path to the upstream request. |`string`| N/A| N|
|`route_id`|Unique identifier of the associated route. |`string`| N/A| N|
|`service_id`|Unique identifier of the associated service. |`string`| N/A| N|


### Computed Fields
The following computed attributes are also available:

| field     | explanation     | type    |
|-----------|-----------------|---------|
|`created_at`|Unix timestamp representing when the plugin was created. |int|

### Import
To import an existing instance of the plugin: `terraform import kong_plugin_request_transformer_advanced.req-transformer <plugin UUID>`
