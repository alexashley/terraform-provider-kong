provider "kong" {
  admin_api_url = "http://localhost:8001"
}

resource "kong_service" "mockbin" {
  name = "mockbin"
  url = "https://mockbin.org/request"
}

resource "kong_route" "mock" {
  service = {
    id = "${kong_service.mockbin.id}"
  },
  paths = ["/mock"]
}

resource "kong_plugin" "response_transformer_plugin_route" {
  route_id = "${kong_route.mock.id}"
  name = "response-transformer"
  config = {
    add.json = "added-by-terraform:true",
    add.headers = "x-parent-resource:route"
  }
}
