# `terraform-provider-kong`

## Provider
```hcl
provider "kong" {
  admin_api_url = "http://localhost:8001"
  rbac_token = "foobar"
}
```

## Resources

- [`kong_consumer`](/docs/kong_consumer.md)
- [`kong_plugin`](/docs/kong_plugin.md)
- [`kong_plugin_openid_connect`](/docs/kong_plugin_openid_connect.md)
- [`kong_plugin_request_transformer_advanced`](/docs/kong_plugin_request_transformer_advanced.md)
- [`kong_route`](/docs/kong_route.md)
- [`kong_service`](/docs/kong_service.md)
