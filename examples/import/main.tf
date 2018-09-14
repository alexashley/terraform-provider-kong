provider "kong" {
  admin_api_url = "http://localhost:8001"
}

resource "kong_service" "service-to-import" {
  name = "service-to-import"
  url = "http://mockbin.org/request"
}
