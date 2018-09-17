package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/client"
	"github.com/alexashley/terraform-provider-kong/kong/provider/test_util"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKongService_basic(t *testing.T) {
	serviceName := fmt.Sprintf("kong-provider-acc-test-%s", acctest.RandString(5))

	var service client.KongService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongServiceDestroy("kong_service.basic-service"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongServiceConfig_basic(serviceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.basic-service", &service),
					testAccCheckKongServiceAttributes(&service, &client.KongService{
						Name: serviceName,
						Url:  "http://foobar.org:5555",
					}),
					resource.TestCheckResourceAttr("kong_service.basic-service", "name", serviceName),
					resource.TestCheckResourceAttr("kong_service.basic-service", "url", "http://foobar.org:5555"),
				),
			},
		},
	})
}

func TestAccKongService_update(t *testing.T) {
	serviceName := fmt.Sprintf("kong-provider-acc-test-%s", acctest.RandString(5))
	var service client.KongService

	updatedServiceName := fmt.Sprintf("%s-update-%s", serviceName, acctest.RandString(5))
	updatedUrl := fmt.Sprintf("http://%s.com", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongServiceDestroy("kong_service.basic-service"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongServiceConfig_basic(serviceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.basic-service", &service),
					testAccCheckKongServiceAttributes(&service, &client.KongService{
						Name: serviceName,
						Url:  "http://foobar.org:5555",
					}),
					resource.TestCheckResourceAttr("kong_service.basic-service", "name", serviceName),
					resource.TestCheckResourceAttr("kong_service.basic-service", "url", "http://foobar.org:5555"),
				),
			},
			{
				Config: testAccKongServiceConfig_update(updatedServiceName, updatedUrl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.basic-service", &service),
					testAccCheckKongServiceAttributes(&service, &client.KongService{
						Name: updatedServiceName,
						Url:  updatedUrl,
					}),
					resource.TestCheckResourceAttr("kong_service.basic-service", "name", updatedServiceName),
					resource.TestCheckResourceAttr("kong_service.basic-service", "url", updatedUrl),
				),
			},
			{
				Config: testAccKongServiceConfig_basic(serviceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.basic-service", &service),
					testAccCheckKongServiceAttributes(&service, &client.KongService{
						Name: serviceName,
						Url:  "http://foobar.org:5555",
					}),
					resource.TestCheckResourceAttr("kong_service.basic-service", "name", serviceName),
					resource.TestCheckResourceAttr("kong_service.basic-service", "url", "http://foobar.org:5555"),
				),
			},
		},
	})
}

func TestAccKongService_defaults(t *testing.T) {
	serviceName := fmt.Sprintf("kong-provider-acc-test-%s", acctest.RandString(5))

	var service client.KongService
	expectedService := client.KongService{
		Name:           serviceName,
		ConnectTimeout: 60000,
		Retries:        5,
		ReadTimeout:    60000,
		WriteTimeout:   60000,
		Url:            "http://foobar.org:5555",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongServiceDestroy("kong_service.basic-service"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongServiceConfig_basic(serviceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.basic-service", &service),
					testAccCheckKongServiceAttributes(&service, &expectedService),
					resource.TestCheckResourceAttr("kong_service.basic-service", "name", serviceName),
					resource.TestCheckResourceAttr("kong_service.basic-service", "url", "http://foobar.org:5555"),
				),
			},
		},
	})
}

func TestAccKongService_override_defaults(t *testing.T) {
	serviceName := fmt.Sprintf("kong-provider-acc-test-%s", acctest.RandString(5))

	var service client.KongService
	expectedService := client.KongService{
		Name:           serviceName,
		ConnectTimeout: 30000,
		Retries:        2,
		ReadTimeout:    30000,
		WriteTimeout:   30000,
		Url:            "http://foobar.org:5555",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongServiceDestroy("kong_service.override-defaults-service"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongServiceConfig_overrideDefaults(&expectedService),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.override-defaults-service", &service),
					testAccCheckKongServiceAttributes(&service, &expectedService),
					resource.TestCheckResourceAttr("kong_service.override-defaults-service", "name", serviceName),
					resource.TestCheckResourceAttr("kong_service.override-defaults-service", "url", "http://foobar.org:5555"),
				),
			},
		},
	})
}

func TestAccKongService_http_should_not_surface_port_80_in_url(t *testing.T) {
	serviceName := fmt.Sprintf("kong-provider-acc-test-%s", acctest.RandString(5))

	var service client.KongService

	httpUrl := fmt.Sprintf("http://%s.com", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongServiceDestroy("kong_service.http-service"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongServiceConfig_http(serviceName, httpUrl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.http-service", &service),
					resource.TestCheckResourceAttr("kong_service.http-service", "name", serviceName),
					resource.TestCheckResourceAttr("kong_service.http-service", "url", httpUrl),
				),
			},
		},
	})
}

func TestAccKongService_https_should_not_surface_port_443_in_url(t *testing.T) {
	serviceName := fmt.Sprintf("kong-provider-acc-test-%s", acctest.RandString(5))

	var service client.KongService

	httpsUrl := fmt.Sprintf("https://%s.org", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongServiceDestroy("kong_service.https-service"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongServiceConfig_https(serviceName, httpsUrl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.https-service", &service),
					resource.TestCheckResourceAttr("kong_service.https-service", "name", serviceName),
					resource.TestCheckResourceAttr("kong_service.https-service", "url", httpsUrl),
				),
			},
		},
	})
}

func testAccCheckKongServiceExists(name string, output *client.KongService) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		r, ok := state.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Service resource not found: %s", name)
		}

		if r.Primary.ID == "" {
			return fmt.Errorf("No id set for %s", name)
		}

		kong := testAccProvider.Meta().(*client.KongClient)

		service, err := kong.GetService(r.Primary.ID)

		if err != nil {
			return err
		}

		*output = *service

		return nil
	}
}

func testAccCheckKongServiceDestroy(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		kong := testAccProvider.Meta().(*client.KongClient)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "kong_service" {
				continue
			}
			service, err := kong.GetService(state.RootModule().Resources[name].Primary.ID)

			if err == nil {
				return fmt.Errorf("Service still exists: %s", service.Id)
			}

			kongError, ok := err.(*client.KongHttpError)

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

func testAccCheckKongServiceAttributes(actualService *client.KongService, expectedService *client.KongService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if actualService.Name != expectedService.Name {
			return test_util.ExpectedAndActualError("Kong service name is wrong", expectedService.Name, actualService.Name)
		}

		if actualService.Url != expectedService.Url {
			return test_util.ExpectedAndActualError("Kong service url is wrong", expectedService.Url, actualService.Url)
		}

		// test other fields, ignoring if the value is empty in expectedService
		if expectedService.WriteTimeout != 0 && actualService.WriteTimeout != expectedService.WriteTimeout {
			return test_util.ExpectedAndActualErrorInt("Kong service write_timeout is wrong", expectedService.WriteTimeout, actualService.WriteTimeout)
		}

		if expectedService.ConnectTimeout != 0 && actualService.ConnectTimeout != expectedService.ConnectTimeout {
			return test_util.ExpectedAndActualErrorInt("Kong service connect_timeout is wrong", expectedService.ConnectTimeout, actualService.ConnectTimeout)
		}

		if expectedService.Retries != 0 && actualService.Retries != expectedService.Retries {
			return test_util.ExpectedAndActualErrorInt("Kong service retries is wrong", expectedService.Retries, actualService.Retries)
		}

		return nil
	}
}

func testAccKongServiceConfig_basic(serviceName string) string {
	return fmt.Sprintf(`
resource "kong_service" "basic-service" {
	name = "%s"
	url = "http://foobar.org:5555"
}`, serviceName)
}

func testAccKongServiceConfig_update(serviceName string, serviceUrl string) string {
	return fmt.Sprintf(`
resource "kong_service" "basic-service" {
	name = "%s",
	url = "%s"
}
`, serviceName, serviceUrl)
}

func testAccKongServiceConfig_http(serviceName string, url string) string {
	return fmt.Sprintf(`
resource "kong_service" "http-service" {
	name = "%s"
	url = "%s"
}`, serviceName, url)
}

func testAccKongServiceConfig_https(serviceName string, url string) string {
	return fmt.Sprintf(`
resource "kong_service" "https-service" {
	name = "%s"
	url = "%s"
}`, serviceName, url)
}

func testAccKongServiceConfig_overrideDefaults(service *client.KongService) string {
	return fmt.Sprintf(`
resource "kong_service" "override-defaults-service" {
	name = "%s"
	url = "%s"
	connect_timeout = %d
	retries = %d
	read_timeout = %d
	write_timeout = %d
}`, service.Name, service.Url, service.ConnectTimeout, service.Retries, service.ReadTimeout, service.WriteTimeout)
}
