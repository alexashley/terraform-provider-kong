provider "kong" {
  admin_api_url = "http://localhost:8001"
}

resource "kong_service" "mockbin_service" {
  name = "mockbin"
  url = "https://mockbin.org/request"
}

resource "kong_route" "mock" {
  service = {
    id = "${kong_service.mockbin_service.id}"
  },
  paths = ["/mock"]
}
