package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/client"
	"github.com/alexashley/terraform-provider-kong/kong/provider/test_util"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

type requestTransformerCrud struct {
	headers     []string
	querystring []string
	body        []string
	uri         string
}

func randomRequestTransformerCrud(prefix string, includeUri bool) *requestTransformerCrud {
	crud := requestTransformerCrud{
		headers:     []string{fmt.Sprintf("%s-header:%s-header-value", prefix, prefix)},
		querystring: []string{fmt.Sprintf("%s-querystring:%s-querystring-value", prefix, prefix)},
		body:        []string{fmt.Sprintf("%s-body:%s-body-value", prefix, prefix)},
	}

	if includeUri {
		crud.uri = fmt.Sprintf("/%s-uri", prefix)
	}

	return &crud
}

func TestAccKongPluginRequestTransformerAdvanced_basic(t *testing.T) {
	var plugin client.KongPlugin

	expectedMethod := test_util.PickOne([]string{"GET", "PUT", "POST", "DELETE", "PATCH"})
	expectedPath := fmt.Sprintf("/kong-provider-test-acc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin_request_transformer_advanced",
			"kong_plugin_request_transformer_advanced.test",
			"request-transformer-advanced"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongPluginRequestTransformerAdvancedConfig_basic(expectedMethod, expectedPath),
				Check: resource.ComposeTestCheckFunc(
					test_util.TestAccCheckKongPluginExists(
						testAccProvider,
						"kong_plugin_request_transformer_advanced.test",
						&plugin,
					),
					test_util.TestAccKongPluginConfigAttributes(&plugin, &client.KongPlugin{
						Config: map[string]interface{}{
							"http_method": expectedMethod,
							"replace": map[string]interface{}{
								"uri":         expectedPath,
								"body":        map[string]interface{}{},
								"headers":     map[string]interface{}{},
								"querystring": map[string]interface{}{},
							},
							"add": map[string]interface{}{
								"body":        map[string]interface{}{},
								"headers":     map[string]interface{}{},
								"querystring": map[string]interface{}{},
							},
							"append": map[string]interface{}{
								"body":        map[string]interface{}{},
								"headers":     map[string]interface{}{},
								"querystring": map[string]interface{}{},
							},
							"remove": map[string]interface{}{
								"body":        map[string]interface{}{},
								"headers":     map[string]interface{}{},
								"querystring": map[string]interface{}{},
							},
							"rename": map[string]interface{}{
								"body":        map[string]interface{}{},
								"headers":     map[string]interface{}{},
								"querystring": map[string]interface{}{},
							},
						},
					}),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "service_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "route_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "consumer_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "enabled", "true"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "http_method", expectedMethod),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "replace_uri", expectedPath),
				),
			},
		},
	})
}

func TestAccKongPluginRequestTransformerAdvanced_update(t *testing.T) {
	var plugin client.KongPlugin

	expectedMethod := test_util.PickOne([]string{"GET", "PUT", "POST", "DELETE", "PATCH"})
	expectedPath := fmt.Sprintf("/kong-provider-test-acc-%s", acctest.RandString(5))

	updatedPath := fmt.Sprintf("%s-update-%s", expectedPath, acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin_request_transformer_advanced",
			"kong_plugin_request_transformer_advanced.test",
			"request-transformer-advanced"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongPluginRequestTransformerAdvancedConfig_basic(expectedMethod, expectedPath),
				Check: resource.ComposeTestCheckFunc(
					test_util.TestAccCheckKongPluginExists(
						testAccProvider,
						"kong_plugin_request_transformer_advanced.test",
						&plugin,
					),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "service_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "route_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "consumer_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "enabled", "true"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "http_method", expectedMethod),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "replace_uri", expectedPath),
				),
			},
			{
				Config: testAccKongPluginRequestTransformerAdvancedConfig_basic(expectedMethod, updatedPath),
				Check: resource.ComposeTestCheckFunc(
					test_util.TestAccCheckKongPluginExists(
						testAccProvider,
						"kong_plugin_request_transformer_advanced.test",
						&plugin,
					),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "service_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "route_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "consumer_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "enabled", "true"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "http_method", expectedMethod),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "replace_uri", updatedPath),
				),
			},
			{
				Config: testAccKongPluginRequestTransformerAdvancedConfig_basic(expectedMethod, expectedPath),
				Check: resource.ComposeTestCheckFunc(
					test_util.TestAccCheckKongPluginExists(
						testAccProvider,
						"kong_plugin_request_transformer_advanced.test",
						&plugin,
					),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "service_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "route_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "consumer_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "enabled", "true"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "http_method", expectedMethod),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "replace_uri", expectedPath),
				),
			},
		},
	})
}

func TestAccKongPluginRequestTransformerAdvanced_all(t *testing.T) {
	var plugin client.KongPlugin

	expectedMethod := test_util.PickOne([]string{"GET", "PUT", "POST", "DELETE", "PATCH"})

	expectedRemoveCrud := randomRequestTransformerCrud("remove", false)
	expectedReplaceCrud := randomRequestTransformerCrud("replace", true)
	expectedRenameCrud := randomRequestTransformerCrud("rename", false)
	expectedAddCrud := randomRequestTransformerCrud("add", false)
	expectedAppendCrud := randomRequestTransformerCrud("append", false)

	//conf := testAccKongPluginRequestTransformerAdvancedConfig_complex(
	//	expectedMethod,
	//	expectedRemoveCrud,
	//	expectedReplaceCrud,
	//	expectedRenameCrud,
	//	expectedAppendCrud,
	//	expectedAddCrud,
	//)
	//
	//util.Log(conf)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: test_util.TestAccCheckGenericKongPluginDestroy(
			testAccProvider,
			"kong_plugin_request_transformer_advanced",
			"kong_plugin_request_transformer_advanced.test",
			"request-transformer-advanced"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongPluginRequestTransformerAdvancedConfig_complex(
					expectedMethod,
					expectedRemoveCrud,
					expectedReplaceCrud,
					expectedRenameCrud,
					expectedAddCrud,
					expectedAppendCrud,
				),
				Check: resource.ComposeTestCheckFunc(
					test_util.TestAccCheckKongPluginExists(
						testAccProvider,
						"kong_plugin_request_transformer_advanced.test",
						&plugin,
					),
					test_util.TestAccKongPluginConfigAttributes(&plugin, &client.KongPlugin{
						Config: map[string]interface{}{
							"http_method": expectedMethod,
							"replace": map[string]interface{}{
								"uri":         expectedReplaceCrud.uri,
								"body":        expectedReplaceCrud.body,
								"headers":     expectedReplaceCrud.headers,
								"querystring": expectedReplaceCrud.querystring,
							},
							"add": map[string]interface{}{
								"body":        expectedAddCrud.body,
								"headers":     expectedAddCrud.headers,
								"querystring": expectedAddCrud.querystring,
							},
							"append": map[string]interface{}{
								"body":        expectedAppendCrud.body,
								"headers":     expectedAppendCrud.headers,
								"querystring": expectedAppendCrud.querystring,
							},
							"remove": map[string]interface{}{
								"body":        expectedRemoveCrud.body,
								"headers":     expectedRemoveCrud.headers,
								"querystring": expectedRemoveCrud.querystring,
							},
							"rename": map[string]interface{}{
								"body":        expectedRenameCrud.body,
								"headers":     expectedRenameCrud.headers,
								"querystring": expectedRenameCrud.querystring,
							},
						},
					}),
					resource.TestCheckResourceAttrSet("kong_plugin_request_transformer_advanced.test", "service_id"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "route_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "consumer_id", ""),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "enabled", "true"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "http_method", expectedMethod),
					// replace
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "replace_uri", expectedReplaceCrud.uri),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "replace_body_params.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "replace_body_params.0", expectedReplaceCrud.body[0]),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "replace_headers.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "replace_headers.0", expectedReplaceCrud.headers[0]),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "replace_querystring.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "replace_querystring.0", expectedReplaceCrud.querystring[0]),
					// remove
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "remove_body_params.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "remove_body_params.0", expectedRemoveCrud.body[0]),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "remove_headers.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "remove_headers.0", expectedRemoveCrud.headers[0]),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "remove_querystring.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "remove_querystring.0", expectedRemoveCrud.querystring[0]),
					// rename
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "rename_body_params.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "rename_body_params.0", expectedRenameCrud.body[0]),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "rename_headers.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "rename_headers.0", expectedRenameCrud.headers[0]),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "rename_querystring.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "rename_querystring.0", expectedRenameCrud.querystring[0]),
					// add
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "add_body_params.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "add_body_params.0", expectedAddCrud.body[0]),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "add_headers.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "add_headers.0", expectedAddCrud.headers[0]),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "add_querystring.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "add_querystring.0", expectedAddCrud.querystring[0]),
					// append
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "append_body_params.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "append_body_params.0", expectedAppendCrud.body[0]),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "append_headers.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "append_headers.0", expectedAppendCrud.headers[0]),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "append_querystring.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin_request_transformer_advanced.test", "append_querystring.0", expectedAppendCrud.querystring[0]),
				),
			},
		},
	})
}

func testAccKongPluginRequestTransformerAdvancedConfig_basic(httpMethod, replaceUri string) string {
	return fmt.Sprintf(`
resource "kong_plugin_request_transformer_advanced" "test" {
	http_method = "%s"
	replace_uri = "%s"
}`, httpMethod, replaceUri)
}

func testAccKongPluginRequestTransformerAdvancedConfig_complex(
	httpMethod string,
	remove, replace, rename, add, append *requestTransformerCrud) string {

	return fmt.Sprintf(`
resource "kong_service" "test-service" {
	name = "test-service-acc-request-transformer-advanced"
	url = "http://mockbin.org"
}

resource "kong_plugin_request_transformer_advanced" "test" {
	service_id = "${kong_service.test-service.id}"

	http_method = "%s"

	remove_headers = ["%s"]
	remove_querystring = ["%s"]
	remove_body_params = ["%s"]
	
	replace_headers = ["%s"] 
	replace_querystring = ["%s"]
	replace_body_params = ["%s"]
	replace_uri = "%s"

	rename_headers = ["%s"]
	rename_querystring = ["%s"]
	rename_body_params = ["%s"]

	add_headers = ["%s"]
	add_querystring = ["%s"]
	add_body_params = ["%s"]

	append_headers = ["%s"]
	append_querystring = ["%s"]
	append_body_params = ["%s"]
}
`,
		httpMethod,

		remove.headers[0],
		remove.querystring[0],
		remove.body[0],

		replace.headers[0],
		replace.querystring[0],
		replace.body[0],
		replace.uri,

		rename.headers[0],
		rename.querystring[0],
		rename.body[0],

		add.headers[0],
		add.querystring[0],
		add.body[0],

		append.headers[0],
		append.querystring[0],
		append.body[0])
}
