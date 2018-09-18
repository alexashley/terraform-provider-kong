package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/alexashley/terraform-provider-kong/kong/provider/test_util"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKongRoute_basic(t *testing.T) {
	serviceName := fmt.Sprintf("kong-provider-acc-test-%s", acctest.RandString(5))
	routePath := fmt.Sprintf("/kong-provider-acc-test-route-%s", acctest.RandString(5))

	var route kong.KongRoute

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongRouteDestroy("kong_route.basic-route"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongRouteConfig_basic(serviceName, routePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRouteExists("kong_route.basic-route", &route),
					testAccCheckKongRouteAttributes(&route, &kong.KongRoute{
						Protocols: []string{"http"},
						Paths:     []string{routePath},
						Methods:   []string{"GET", "PUT", "DELETE"},
					}),
					resource.TestCheckResourceAttr("kong_route.basic-route", "protocols.#", "1"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "protocols.0", "http"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "methods.#", "3"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "methods.0", "GET"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "methods.1", "PUT"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "methods.2", "DELETE"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "hosts.#", "0"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "paths.#", "1"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "paths.0", routePath),
					resource.TestCheckResourceAttr("kong_route.basic-route", "strip_path", "true"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "preserve_host", "false"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "regex_priority", "0"),
				),
			},
		},
	})
}

func TestAccKongRoute_update(t *testing.T) {
	serviceName := fmt.Sprintf("kong-provider-acc-test-%s", acctest.RandString(5))
	routePath := fmt.Sprintf("/kong-provider-acc-test-route-%s", acctest.RandString(5))
	updatedRoutePath := fmt.Sprintf("%s-update-%s", routePath, acctest.RandString(5))
	var route kong.KongRoute

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongRouteDestroy("kong_route.basic-route"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongRouteConfig_basic(serviceName, routePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRouteExists("kong_route.basic-route", &route),
					resource.TestCheckResourceAttr("kong_route.basic-route", "paths.#", "1"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "paths.0", routePath),
				),
			},
			{
				Config: testAccKongRouteConfig_basic(serviceName, updatedRoutePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRouteExists("kong_route.basic-route", &route),
					resource.TestCheckResourceAttr("kong_route.basic-route", "paths.#", "1"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "paths.0", updatedRoutePath),
				),
			},
			{
				Config: testAccKongRouteConfig_basic(serviceName, routePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRouteExists("kong_route.basic-route", &route),
					resource.TestCheckResourceAttr("kong_route.basic-route", "paths.#", "1"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "paths.0", routePath),
				),
			},
		},
	})
}

func TestAccKongRoute_host(t *testing.T) {
	serviceName := fmt.Sprintf("kong-provider-acc-test-%s", acctest.RandString(5))
	host := fmt.Sprintf("kong-provider-acc-test-route-%s.org", acctest.RandString(5))

	var route kong.KongRoute

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongRouteDestroy("kong_route.host-route"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongRouteConfig_host(serviceName, host),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRouteExists("kong_route.host-route", &route),
					testAccCheckKongRouteAttributes(&route, &kong.KongRoute{
						Hosts:     []string{host},
						Protocols: []string{"http", "https"},
					}),
					resource.TestCheckResourceAttr("kong_route.host-route", "protocols.#", "2"),
					resource.TestCheckResourceAttr("kong_route.host-route", "protocols.0", "http"),
					resource.TestCheckResourceAttr("kong_route.host-route", "protocols.1", "https"),
					resource.TestCheckResourceAttr("kong_route.host-route", "methods.#", "0"),
					resource.TestCheckResourceAttr("kong_route.host-route", "hosts.#", "1"),
					resource.TestCheckResourceAttr("kong_route.host-route", "hosts.0", host),
					resource.TestCheckResourceAttr("kong_route.host-route", "paths.#", "0"),
				),
			},
		},
	})
}

func testAccCheckKongRouteExists(name string, output *kong.KongRoute) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		r, ok := state.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Route resource not found: %s", name)
		}

		if r.Primary.ID == "" {
			return fmt.Errorf("No id set for %s", name)
		}

		kong := testAccProvider.Meta().(*kong.KongClient)

		route, err := kong.GetRoute(r.Primary.ID)

		if err != nil {
			return err
		}

		*output = *route

		return nil
	}
}

func testAccCheckKongRouteDestroy(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := testAccProvider.Meta().(*kong.KongClient)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "kong_route" {
				continue
			}
			route, err := client.GetRoute(state.RootModule().Resources[name].Primary.ID)

			if err == nil {
				return fmt.Errorf("Route still exists: %s", route.Id)
			}

			kongError, ok := err.(*kong.HttpError)

			if !ok {
				return err
			}

			if kongError.StatusCode != 404 {
				return kongError
			}

			return nil
		}

		return nil
	}
}

func testAccCheckKongRouteAttributes(actualRoute *kong.KongRoute, expectedRoute *kong.KongRoute) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if !slicesAreEqual(actualRoute.Protocols, expectedRoute.Protocols) {
			return test_util.ExpectedAndActualErrorStringSlice("Protocols don't match", expectedRoute.Protocols, actualRoute.Protocols)
		}

		if !slicesAreEqual(actualRoute.Methods, expectedRoute.Methods) {
			return test_util.ExpectedAndActualErrorStringSlice("Methods don't match", expectedRoute.Methods, actualRoute.Methods)
		}

		if !slicesAreEqual(actualRoute.Hosts, expectedRoute.Hosts) {
			return test_util.ExpectedAndActualErrorStringSlice("Paths don't match", expectedRoute.Hosts, actualRoute.Hosts)
		}

		if !slicesAreEqual(actualRoute.Paths, expectedRoute.Paths) {
			return test_util.ExpectedAndActualErrorStringSlice("Paths don't match", expectedRoute.Paths, actualRoute.Paths)
		}

		return nil
	}
}

func testAccKongRouteConfig_basic(serviceName, path string) string {
	return fmt.Sprintf(`
resource "kong_service" "basic-service" {
	name = "%s"
	url = "http://foobar.org:5555"
}

resource "kong_route" "basic-route" {
	service_id = "${kong_service.basic-service.id}"
	paths = ["%s"],
	methods = ["GET", "PUT", "DELETE"]
	protocols = ["http"]
}
`, serviceName, path)
}

func testAccKongRouteConfig_host(serviceName, host string) string {
	return fmt.Sprintf(`
resource "kong_service" "host-service" {
	name = "%s"
	url = "http://foobar.org:5555"
}

resource "kong_route" "host-route" {
	service_id = "${kong_service.host-service.id}"
	hosts = ["%s"],
}`, serviceName, host)
}

func slicesAreEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
