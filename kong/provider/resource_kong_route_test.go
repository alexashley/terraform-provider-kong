package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/alexashley/terraform-provider-kong/kong/provider/test_util"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"regexp"
	"sort"
	"strings"
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
					test_util.AssertValueInTerraformSet("kong_route.basic-route", "protocols", "http"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "methods.#", "3"),
					test_util.AssertValueInTerraformSet("kong_route.basic-route", "methods", "GET"),
					test_util.AssertValueInTerraformSet("kong_route.basic-route", "methods", "PUT"),
					test_util.AssertValueInTerraformSet("kong_route.basic-route", "methods", "DELETE"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "hosts.#", "0"),
					resource.TestCheckResourceAttr("kong_route.basic-route", "paths.#", "1"),
					test_util.AssertValueInTerraformSet("kong_route.basic-route", "paths", routePath),
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
					test_util.AssertValueInTerraformSet("kong_route.basic-route", "paths", routePath),
				),
			},
			{
				Config: testAccKongRouteConfig_basic(serviceName, updatedRoutePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRouteExists("kong_route.basic-route", &route),
					resource.TestCheckResourceAttr("kong_route.basic-route", "paths.#", "1"),
					test_util.AssertValueInTerraformSet("kong_route.basic-route", "paths", updatedRoutePath),
				),
			},
			{
				Config: testAccKongRouteConfig_basic(serviceName, routePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRouteExists("kong_route.basic-route", &route),
					resource.TestCheckResourceAttr("kong_route.basic-route", "paths.#", "1"),
					test_util.AssertValueInTerraformSet("kong_route.basic-route", "paths", routePath),
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
					test_util.AssertValueInTerraformSet("kong_route.host-route", "protocols", "http"),
					test_util.AssertValueInTerraformSet("kong_route.host-route", "protocols", "https"),
					resource.TestCheckResourceAttr("kong_route.host-route", "methods.#", "0"),
					resource.TestCheckResourceAttr("kong_route.host-route", "hosts.#", "1"),
					test_util.AssertValueInTerraformSet("kong_route.host-route", "hosts", host),
					resource.TestCheckResourceAttr("kong_route.host-route", "paths.#", "0"),
				),
			},
		},
	})
}

func TestAccKongRoute_validate_protocols(t *testing.T) {
	serviceName := fmt.Sprintf("kong-provider-acc-test-%s", acctest.RandString(5))
	invalidProtocol := acctest.RandString(10)
	config := fmt.Sprintf(`
resource "kong_service" "test" {
	name = "%s",
	url = "https://foobar.org"
}

resource "kong_route" "invalid-protocol-route" {
	service_id = "${kong_service.test.id}"

	paths = ["/foo"]
	protocols = ["%s"]
}`, serviceName, invalidProtocol)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("the supplied protocols are not supported by Kong: " + invalidProtocol),
			},
		},
	})
}

func TestAccKongRoute_validate_methods(t *testing.T) {
	serviceName := fmt.Sprintf("kong-provider-acc-test-%s", acctest.RandString(5))
	invalidMethods := []string{acctest.RandString(10), acctest.RandString(5)}
	sort.Slice(invalidMethods, func(i, j int) bool {
		return invalidMethods[i] < invalidMethods[j]
	})

	config := fmt.Sprintf(`
resource "kong_service" "test" {
	name = "%s"
	url = "http://foobar.org"
}

resource "kong_route" "invalid-methods-route" {
	service_id = "${kong_service.test.id}"

	paths = ["/foobar"]
	methods = ["%s", "%s"]
}`, serviceName, invalidMethods[0], invalidMethods[1])

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(fmt.Sprintf("invalid HTTP methods: %s", strings.Join(invalidMethods, ", "))),
			},
		},
	})
}

func TestAccKongRoute_validate_config(t *testing.T) {
	serviceName := fmt.Sprintf("kong-provider-acc-test-%s", acctest.RandString(5))

	config := fmt.Sprintf(`
resource "kong_service" "test" {
	name = "%s"
	url = "http://foobar.org"
}

resource "kong_route" "invalid-route" {
	service_id = "${kong_service.test.id}"
}`, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("at least one of methods, paths, or hosts must be set in order for Kong to proxy traffic to this route"),
			},
		},
	})
}

func testAccCheckKongRouteExists(name string, output *kong.KongRoute) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		r, ok := state.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("route resource not found: %s", name)
		}

		if r.Primary.ID == "" {
			return fmt.Errorf("no id set for %s", name)
		}

		client := testAccProvider.Meta().(*kong.KongClient)

		route, err := client.GetRoute(r.Primary.ID)

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
				return fmt.Errorf("route still exists: %s", route.Id)
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

		if !setsAreEqual(actualRoute.Protocols, expectedRoute.Protocols) {
			return test_util.ExpectedAndActualErrorStringSlice("Protocols don't match", expectedRoute.Protocols, actualRoute.Protocols)
		}

		if !setsAreEqual(actualRoute.Methods, expectedRoute.Methods) {
			return test_util.ExpectedAndActualErrorStringSlice("Methods don't match", expectedRoute.Methods, actualRoute.Methods)
		}

		if !setsAreEqual(actualRoute.Hosts, expectedRoute.Hosts) {
			return test_util.ExpectedAndActualErrorStringSlice("Paths don't match", expectedRoute.Hosts, actualRoute.Hosts)
		}

		if !setsAreEqual(actualRoute.Paths, expectedRoute.Paths) {
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

func setsAreEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		match := false
		for j := range b {
			if a[i] == b[j] {
				match = true
				break
			}
		}

		if !match {
			return false
		}
	}

	return true
}
