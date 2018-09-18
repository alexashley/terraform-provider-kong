package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/alexashley/terraform-provider-kong/kong/provider/test_util"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"strings"
	"testing"
)

func TestAccKongPluginIpHeaderRestriction_basic_global_plugin(t *testing.T) {
	var plugin kong.KongPlugin

	var expectedIpAddresses = []string{test_util.RandomIp(), test_util.RandomIp()}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin_ip_header_restriction",
			"kong_plugin_ip_header_restriction.test",
			"ip-header-restriction"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongPluginIpHeaderRestrictionConfig_basic(expectedIpAddresses),
				Check: resource.ComposeTestCheckFunc(
					test_util.TestAccCheckKongPluginExists(
						testAccProvider,
						"kong_plugin_ip_header_restriction.test",
						&plugin,
					),
					test_util.TestAccKongPluginConfigAttributes(&plugin, &kong.KongPlugin{
						Config: map[string]interface{}{
							"whitelist":       expectedIpAddresses,
							"ip_header":       "x-forwarded-for",
							"override_global": false,
						},
					}),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "service_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "route_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "consumer_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "enabled", "true"),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.#", "2"),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.0", expectedIpAddresses[0]),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.1", expectedIpAddresses[1]),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "ip_header", "x-forwarded-for"),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "override_global", "false"),
				),
			},
		},
	})
}

func TestAccKongPluginIpHeaderRestriction_basic_service_plugin(t *testing.T) {
	var plugin kong.KongPlugin

	var expectedIpAddresses = []string{test_util.RandomIp()}
	serviceName := "kong-terraform-acc-test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin_ip_header_restriction",
			"kong_plugin_ip_header_restriction.test",
			"ip-header-restriction"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongPluginIpHeaderRestrictionConfig_service_plugin(expectedIpAddresses, serviceName),
				Check: resource.ComposeTestCheckFunc(
					test_util.TestAccCheckKongPluginExists(
						testAccProvider,
						"kong_plugin_ip_header_restriction.test",
						&plugin,
					),
					resource.TestCheckResourceAttrSet("kong_plugin_ip_header_restriction.test", "service_id"),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "route_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "consumer_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.0", expectedIpAddresses[0]),
				),
			},
		},
	})
}

func TestAccKongPluginIpHeaderRestriction_basic_route_plugin(t *testing.T) {
	var plugin kong.KongPlugin

	var expectedIpAddresses = []string{test_util.RandomIp(), test_util.RandomIp(), test_util.RandomIp()}
	serviceName := "kong-terraform-acc-test-" + acctest.RandString(10)
	routePath := "/kong-terraform-acc-test" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin_ip_header_restriction",
			"kong_plugin_ip_header_restriction.test",
			"ip-header-restriction"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongPluginIpHeaderRestrictionConfig_route_plugin(expectedIpAddresses, serviceName, routePath),
				Check: resource.ComposeTestCheckFunc(
					test_util.TestAccCheckKongPluginExists(
						testAccProvider,
						"kong_plugin_ip_header_restriction.test",
						&plugin,
					),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "service_id", ""),
					resource.TestCheckResourceAttrSet("kong_plugin_ip_header_restriction.test", "route_id"),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "consumer_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.#", "3"),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.0", expectedIpAddresses[0]),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.1", expectedIpAddresses[1]),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.2", expectedIpAddresses[2]),
				),
			},
		},
	})
}

func TestAccKongPluginIpHeaderRestriction_update(t *testing.T) {
	var plugin kong.KongPlugin

	var expectedIpAddresses = []string{test_util.RandomIp(), test_util.RandomIp()}
	var updatedIpAddresses = []string{test_util.RandomIp()}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin_ip_header_restriction",
			"kong_plugin_ip_header_restriction.test",
			"ip-header-restriction"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongPluginIpHeaderRestrictionConfig_basic(expectedIpAddresses),
				Check: resource.ComposeTestCheckFunc(
					test_util.TestAccCheckKongPluginExists(
						testAccProvider,
						"kong_plugin_ip_header_restriction.test",
						&plugin,
					),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.#", "2"),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.0", expectedIpAddresses[0]),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.1", expectedIpAddresses[1]),
				),
			},
			{
				Config: testAccKongPluginIpHeaderRestrictionConfig_basic(updatedIpAddresses),
				Check: resource.ComposeTestCheckFunc(
					test_util.TestAccCheckKongPluginExists(
						testAccProvider,
						"kong_plugin_ip_header_restriction.test",
						&plugin,
					),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.0", updatedIpAddresses[0]),
				),
			},
			{
				Config: testAccKongPluginIpHeaderRestrictionConfig_basic(expectedIpAddresses),
				Check: resource.ComposeTestCheckFunc(
					test_util.TestAccCheckKongPluginExists(
						testAccProvider,
						"kong_plugin_ip_header_restriction.test",
						&plugin,
					),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.#", "2"),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.0", expectedIpAddresses[0]),
					resource.TestCheckResourceAttr("kong_plugin_ip_header_restriction.test", "whitelist.1", expectedIpAddresses[1]),
				),
			},
		},
	})
}

func testAccKongPluginIpHeaderRestrictionConfig_basic(ipAddresses []string) string {
	return fmt.Sprintf(`
resource "kong_plugin_ip_header_restriction" "test" {
	whitelist = ["%s"]
}`, strings.Join(ipAddresses, "\",\""))
}

func testAccKongPluginIpHeaderRestrictionConfig_service_plugin(ipAddresses []string, serviceName string) string {
	return fmt.Sprintf(`
resource "kong_service" "ip-plugin-test" {
	name = "%s"
	url = "http://mockbin.org"
}

resource "kong_plugin_ip_header_restriction" "test" {
	service_id = "${kong_service.ip-plugin-test.id}"
	whitelist = ["%s"]
}
`, serviceName, strings.Join(ipAddresses, "\",\""))
}

func testAccKongPluginIpHeaderRestrictionConfig_route_plugin(ipAddresses []string, serviceName, routePath string) string {
	return fmt.Sprintf(`
resource "kong_service" "ip-plugin-test" {
	name = "%s"
	url = "http://mockbin.org"
}

resource "kong_route" "ip-plugin-test" {
	service_id = "${kong_service.ip-plugin-test.id}"
	paths = ["%s"]
}

resource "kong_plugin_ip_header_restriction" "test" {
	route_id = "${kong_route.ip-plugin-test.id}"
	whitelist = ["%s"]
}
`, serviceName, routePath, strings.Join(ipAddresses, "\",\""))
}
