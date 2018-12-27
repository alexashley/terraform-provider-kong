# `terraform-provider-kong`

## Provider
~~~
provider "kong" {
  admin_api_url = "http://localhost:8001"
  rbac_token    = "foobar"
}
~~~

## Resources

- [`kong_consumer`](https://alexashley.github.io/terraform-provider-kong/kong_consumer)
- [`kong_plugin`](https://alexashley.github.io/terraform-provider-kong/kong_plugin)
- [`kong_plugin_openid_connect`](https://alexashley.github.io/terraform-provider-kong/kong_plugin_openid_connect)
- [`kong_plugin_request_transformer_advanced`](https://alexashley.github.io/terraform-provider-kong/kong_plugin_request_transformer_advanced)
- [`kong_route`](https://alexashley.github.io/terraform-provider-kong/kong_route)
- [`kong_service`](https://alexashley.github.io/terraform-provider-kong/kong_service)

[GitHub](https://github.com/alexashley/terraform-provider-kong)
