package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/provider/test_util"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"regexp"
	"testing"
)

func TestAccKongPluginOpenIdConnect_basic_update(t *testing.T) {
	issuer := fmt.Sprintf("https://%s.com", acctest.RandString(10))
	updatedIssuer := fmt.Sprintf("https://%s.com", acctest.RandString(10))
	anonymousConsumer := uuid.New().String()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin_openid_connect",
			"kong_plugin_openid_connect.oidc-test",
			"openid-connect"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongPluginOpenIdConnect_basic(issuer),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("kong_plugin_openid_connect.oidc-test", "issuer", issuer),
					resource.TestCheckResourceAttr("kong_plugin_openid_connect.oidc-test", "auth_methods.#", "1"),
					test_util.AssertValueInTerraformSet("kong_plugin_openid_connect.oidc-test", "auth_methods", "bearer"),
					resource.TestCheckResourceAttr("kong_plugin_openid_connect.oidc-test", "consumer_claim.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_openid_connect.oidc-test", "consumer_by.#", "1"),
					test_util.AssertValueInTerraformSet("kong_plugin_openid_connect.oidc-test", "consumer_claim", "sub"),
					resource.TestCheckResourceAttrSet("kong_plugin_openid_connect.oidc-test", "service_id"),
				),
			},
			{
				ResourceName:      "kong_plugin_openid_connect.oidc-test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccKongPluginOpenIdConnect_anonymous(updatedIssuer, anonymousConsumer),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("kong_plugin_openid_connect.oidc-test", "issuer", updatedIssuer),
					resource.TestCheckResourceAttr("kong_plugin_openid_connect.oidc-test", "anonymous", anonymousConsumer),
				),
			},
		},
	})
}

func TestAccKongPluginOpenIdConnect_validate_anonymous(t *testing.T) {
	issuer := fmt.Sprintf("https://%s.com", acctest.RandString(10))
	anonymousConsumer := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin_openid_connect",
			"kong_plugin_openid_connect.oidc-test",
			"openid-connect"),
		Steps: []resource.TestStep{
			{
				Config:      testAccKongPluginOpenIdConnect_anonymous(issuer, anonymousConsumer),
				ExpectError: regexp.MustCompile(anonymousConsumer + " is not a valid UUID"),
			},
		},
	})
}

func TestAccKongPluginOpenIdConnect_validate_issuer(t *testing.T) {
	issuer := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin_openid_connect",
			"kong_plugin_openid_connect.oidc-test",
			"openid-connect"),
		Steps: []resource.TestStep{
			{
				Config:      testAccKongPluginOpenIdConnect_basic(issuer),
				ExpectError: regexp.MustCompile("error parsing url"),
			},
		},
	})
}

func TestAccKongPluginOpenIdConnect_validate_auth_methods(t *testing.T) {
	issuer := fmt.Sprintf("https://%s.com", acctest.RandString(10))
	authMethod := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin_openid_connect",
			"kong_plugin_openid_connect.oidc-test",
			"openid-connect"),
		Steps: []resource.TestStep{
			{
				Config:      testAccKongPluginOpenIdConnect_auth(issuer, authMethod),
				ExpectError: regexp.MustCompile(authMethod + " is not a valid auth_method"),
			},
		},
	})
}

func TestAccKongPluginOpenIdConnect_validate_consumer_by(t *testing.T) {
	issuer := fmt.Sprintf("https://%s.com", acctest.RandString(10))
	consumerBy := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin_openid_connect",
			"kong_plugin_openid_connect.oidc-test",
			"openid-connect"),
		Steps: []resource.TestStep{
			{
				Config:      testAccKongPluginOpenIdConnect_consumer_by(issuer, consumerBy),
				ExpectError: regexp.MustCompile("invalid value for consumer_by: must be one of custom_id or username"),
			},
		},
	})
}

func testAccKongPluginOpenIdConnect_basic(issuer string) string {
	return fmt.Sprintf(`
	resource "kong_service" "test-service" {
		name	= "mockbin-%s"
		url		= "https://mockbin.org/request"
	}

	resource "kong_plugin_openid_connect" "oidc-test" {
		service_id		= "${kong_service.test-service.id}"
		issuer 			= "%s"
		auth_methods 	= ["bearer"]
		consumer_claim	= ["sub"]
		consumer_by		= ["username"]
	}
`, acctest.RandString(5), issuer)
}

func testAccKongPluginOpenIdConnect_anonymous(issuer, anonymous string) string {
	return fmt.Sprintf(`
	resource "kong_service" "test-service" {
		name	= "mockbin-%s"
		url		= "https://mockbin.org/request"
	}

	resource "kong_plugin_openid_connect" "oidc-test" {
		service_id		= "${kong_service.test-service.id}"
		issuer 			= "%s"
		auth_methods 	= ["bearer"]
		consumer_claim	= ["sub"]
		consumer_by		= ["username"]
		anonymous		= "%s"
	}
`, acctest.RandString(5), issuer, anonymous)
}

func testAccKongPluginOpenIdConnect_auth(issuer, authMethod string) string {
	return fmt.Sprintf(`
	resource "kong_service" "test-service" {
		name	= "mockbin-%s"
		url		= "https://mockbin.org/request"
	}

	resource "kong_plugin_openid_connect" "oidc-test" {
		service_id		= "${kong_service.test-service.id}"
		issuer 			= "%s"
		auth_methods 	= ["%s"]
	}
`, acctest.RandString(5), issuer, authMethod)
}

func testAccKongPluginOpenIdConnect_consumer_by(issuer, consumerBy string) string {
	return fmt.Sprintf(`
	resource "kong_service" "test-service" {
		name	= "mockbin-%s"
		url		= "https://mockbin.org/request"
	}

	resource "kong_plugin_openid_connect" "oidc-test" {
		service_id		= "${kong_service.test-service.id}"
		issuer 			= "%s"
		consumer_by 	= ["%s"]
		
	}
`, acctest.RandString(5), issuer, consumerBy)
}
