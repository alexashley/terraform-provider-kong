# kong-importer

A command-line utility to import Kong resources into Terraform. 
It doesn't currently support HCL generation (although I would like to add this feature), but it can help ease the pain of managing an existing Kong instance with Terraform.
Simply point it at a Kong installation and it will attempt to bulk import all consumers, services, routes, and plugins it can find.

## example

`./kong-import import -admin-api-url=http://localhost:8001 -tf-config=examples/import`

## flags

| flag            | description                                                                                   | default                 |
|-----------------|-----------------------------------------------------------------------------------------------|-------------------------|
| `admin-api-url` | Kong's admin api url. Usually listening on port 8001.                                         | `http://localhost:8001` |
| `rbac-token`    | Kong EE RBAC token. Only necessary if your Kong Enterprise installation is secured with RBAC. |                         |
| `dry-run`       | Discover all resources and print what resources would be imported                             | `false`                 |
| `state`         | Holds the current import state and any exclusions                                             | `import-state.json`     |
| `tf-config`     | Path to Terraform config directory                                                            |                         |

## resource naming

As a general rule, any spaces or hyphens are replaced with underscores.

- `kong_consumer`: First tries to use the `username`, then falls back to `custom_id`
- `kong_service`: Uses the `name` field off the service object. Any spaces are replaced with underscores.
- `kong_route`: Named by exploding the first path or host and prepending the service name. If the route begins with the service name, it will not be included twice.
    - service: `web-store`, route: `/v2/products/` -> `kong_route.web_store_v2_products` 
    - service: `products`, route: `/products` -> `kong_route.products`
- `kong_plugin`: The plugin name is prefixed with the associated service, route, or consumer name. Global plugins are prefixed with `global`
    - no service, route, or consumer: `kong_plugin.global_basic_auth`
    - service named `foo`: `kong_plugin.foo_basic_auth`
- other plugin resources: Similar to `kong_plugin`, but the resource name is simply the service or route name.    
    - `kong_plugin_openid_connect.products`
    - `kong_plugin_request_transformer_advanced.web_store_v2_products`

## failure/retries & ignoring resources

To keep track of progress, the tool stores the UUIDs of successfully imported resources in a file called `import-state.json`.
If there's a failure, the result of stderr is printed to the console to aid in debugging; once the issue has been resolved, it should be safe to run again.

If you wish to exclude certain resources from import, you can add their ids to the appropriate place in the file and they'll be ignored by the importer.

```json
{
  "consumers": [
    "ca6f8710-929c-4a75-b9f9-5933caf40929"
  ],
  "plugins": [],
  "routes": [],
  "services": [
    "2358625f-2070-495b-8c0a-3ca3eb5af472"
  ]
}
```

