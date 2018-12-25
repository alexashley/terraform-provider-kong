provider "kong" {
  admin_api_url = "http://localhost:8001"
}

resource "kong_service" "mockbin" {
  name = "mockbin"
  url = "https://mockbin.org/request"
}

resource "kong_route" "mock" {
  service_id = "${kong_service.mockbin.id}"
  paths = [
    "/mock"]
}

resource "kong_plugin_request_transformer_advanced" "request-transformer-plugin-service" {
  service_id = "${kong_service.mockbin.id}"
  add_headers = ["x-parent-resource:service"]
  http_method = "GET",
  replace_uri = "/foobar"
}

variable "test-count-enabled" {
  default = "no"
  type = "string"
}

resource "kong_service" "test-count-service" {
  count = "${var.test-count-enabled == "yes" ? 1 : 0}"

  name = "foo"
  url = "http://foo.bar.org"
}

resource "kong_route" "child-route-test-count" {
  /*
    Need to set count for dependent resources of a conditional resource,
    otherwise TF will error:
    Resource 'kong_service.test-count-service' not found for variable 'kong_service.test-count-service.id'
  */
  count = "${var.test-count-enabled == "yes" ? 1 : 0}"

  service_id = "${kong_service.test-count-service.id}"
}

variable basic-auth-foo-config {
  type = "string"
  default =<<EOF
{
  "anonymous": "",
  "hide_credentials": true
}
EOF
}

resource "kong_plugin" "basic-auth-foo" {
  route_id = "${kong_route.mock.id}"
  name = "basic-auth"
  config_json = "${var.basic-auth-foo-config}"
}
