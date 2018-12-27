package provider

import (
	"encoding/json"
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/provider/test_util"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"reflect"
	"regexp"
	"testing"
)

func TestAccKongPlugin_basic_update(t *testing.T) {
	configJson := randomNopPluginConfig()
	updatedConfigJson := randomNopPluginConfig()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin",
			"kong_plugin.generic-test",
			"nop"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongPlugin_basic("nop", configJson),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("kong_plugin.generic-test", "name", "nop"),
					checkJson("kong_plugin.generic-test", "config_json", configJson),
				),
			},
			{
				ResourceName:      "kong_plugin.generic-test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccKongPlugin_basic("nop", updatedConfigJson),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("kong_plugin.generic-test", "name", "nop"),
					checkJson("kong_plugin.generic-test", "config_json", updatedConfigJson),
				),
			},
		},
	})
}

func TestAccKongPlugin_specific_resource(t *testing.T) {
	configJson := randomNopPluginConfig()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccKongPlugin_basic("openid-connect", configJson),
				ExpectError: regexp.MustCompile("plugin openid-connect has a resource implementation: kong_plugin_openid_connect"),
			},
		},
	})
}

func TestAccKongPlugin_specific_resource_override(t *testing.T) {
	os.Setenv("TF_KONG_ALLOW_GENERIC_PLUGIN_OPENID_CONNECT", "YES")

	issuer := fmt.Sprintf("https://%s.com", acctest.RandString(10))
	expectedJson := fmt.Sprintf(`{"issuer": "%s", "auth_methods": ["bearer"]}`, issuer)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKongPlugin_oidc(issuer),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("kong_plugin.test-generic-oidc", "name", "openid-connect"),
					checkJson("kong_plugin.test-generic-oidc", "config_json", expectedJson),
				),
			},
		},
	})

	os.Unsetenv("TF_KONG_ALLOW_GENERIC_PLUGIN_OPENID_CONNECT")
}

func TestAccKongPlugin_invalid_json(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccKongPlugin_basic("nop", acctest.RandString(15)),
				ExpectError: regexp.MustCompile("\"config_json\" contains an invalid JSON:"),
			},
		},
	})
}

func randomNopPluginConfig() string {
	config := map[string]string{
		"arg1": acctest.RandString(10),
		"arg2": acctest.RandString(10),
		"arg3": acctest.RandString(10),
	}

	configJson, _ := json.Marshal(config)

	return string(configJson[:])
}

func testAccKongPlugin_basic(name, config string) string {
	return fmt.Sprintf(`
	resource "kong_service" "test-service" {
		name	= "mockbin-%s"
		url		= "https://mockbin.org/request"
	}

	resource "kong_plugin" "generic-test" {
		service_id	= "${kong_service.test-service.id}"
		name 		= "%s"
		config_json	= <<CONFIG
			%s
		CONFIG
	}`, acctest.RandString(5), name, config)
}

func testAccKongPlugin_oidc(issuer string) string {
	return fmt.Sprintf(`
	resource "kong_service" "test-service" {
		name	= "mockbin-%s"
		url		= "https://mockbin.org/request"
	}

	resource "kong_plugin" "test-generic-oidc" {
		name 		= "openid-connect",
		config_json	= <<CONFIG
		{
			"auth_methods": ["bearer"],
			"issuer": "%s"
		}
		CONFIG
	}`, acctest.RandString(10), issuer)
}

func checkJson(resourceName, attributeName, expectedJson string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		expected := make(map[string]interface{})
		actual := make(map[string]interface{})

		if err := json.Unmarshal([]byte(expectedJson), &expected); err != nil {
			return err
		}

		r, ok := state.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("%s not found", resourceName)
		}

		attribute, ok := r.Primary.Attributes[attributeName]

		if !ok {
			return fmt.Errorf("attribute %s not found on resource %s", attributeName, resourceName)
		}

		if err := json.Unmarshal([]byte(attribute), &actual); err != nil {
			return err
		}

		if !reflect.DeepEqual(expected, actual) {
			return test_util.ExpectedAndActualError(
				fmt.Sprintf("%s did not match", attributeName),
				fmt.Sprintf("%v", expected),
				fmt.Sprintf("%v", actual))
		}

		return nil
	}
}
