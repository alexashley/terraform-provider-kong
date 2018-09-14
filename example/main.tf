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

resource "kong_plugin" "response-transformer-plugin-route" {
  route_id = "${kong_route.mock.id}"
  name = "response-transformer"
  config = {
    add.json = "added-by-terraform:true",
    add.headers = "x-parent-resource:route"
  }
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

resource "kong_plugin_ip_header_restriction" "global-ip-header-restriction" {
  whitelist = [
    "230.188.209.188",
    "141.201.15.252",
    "142.121.207.218"
  ]
}

resource "kong_plugin_ip_header_restriction" "route-ip-header-restriction" {
  route_id = "${kong_route.mock.id}"
  whitelist = ["181.28.140.88"]
}


// imported resources
// run ./create-resources-to-import.sh to create them
//resource "kong_service" "service-to-import" {
//  name = "service-to-import"
//  url = "http://mockbin.org/request"
//}

