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

resource "kong_service" "service-from-url" {
  name = "service-from-url",
  url = "https://foobar.org:8080/test"
}
