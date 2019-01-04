provider "kong" {
  admin_api_url = "http://localhost:8001"
}

resource "kong_service" "import" {
  name = "import"
  url  = "http://mockbin.org/request"
}

resource "kong_route" "import" {
  service_id = "${kong_service.import.id}"
  paths      = [
    "/import"
  ]
}

resource "kong_plugin" "import_basic_auth" {
  route_id    = "${kong_route.import.id}"
  name        = "basic-auth"
  config_json = <<CONFIG
  {
    "hide_credentials": false,
    "anonymous": ""
  }
  CONFIG
}

resource "kong_plugin_openid_connect" "import" {
  service_id   = "${kong_service.import.id}"
  issuer       = "http://foo.org"
  auth_methods = [
    "bearer"]
}

resource "kong_plugin_request_transformer_advanced" "global" {}

resource "kong_consumer" "steve_irwin" {
  username = "steve-irwin"
}
