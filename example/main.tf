provider "kong" {
  admin_api_url = "http://localhost:8001"
}

resource "kong_service" "mockbin-service" {
  name = "mockbin"
  protocol = "http"
  host = "mockbin.org"
  port = 80
  path = "/request"
}
