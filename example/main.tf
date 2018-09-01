provider "kong" {
  admin_api_url = "http://localhost:8001"
}

resource "kong_service" "mockbin_service" {
  name = "mockbin"
  protocol = "http"
  host = "mockbin.org"
  port = 80
  path = "/request"
}

resource "kong_route" "mock" {
  service = {
    id = "${kong_service.mockbin_service.id}"
  },
  paths = ["/mock"]
}

resource "kong_service" "service_from_url" {
  name = "service-from-url",
  url = "https://foobar.org:8080/test"
}
