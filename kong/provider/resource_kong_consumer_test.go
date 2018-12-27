package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/alexashley/terraform-provider-kong/kong/provider/test_util"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"regexp"
	"testing"
)

func TestAccKongConsumer_basic(t *testing.T) {
	var consumer kong.KongConsumer
	username := "kong-provider-acc-test-" + acctest.RandString(5)
	customId := "kong-provider-acc-test-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongConsumerDestroy("kong_consumer.test"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongConsumerConfig_basic(username, customId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerExists("kong_consumer.test", &consumer),
					resource.TestCheckResourceAttr("kong_consumer.test", "username", username),
					resource.TestCheckResourceAttr("kong_consumer.test", "custom_id", customId),
				),
			},
		},
	})
}

func TestAccKongConsumer_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongConsumerDestroy("kong_consumer.test"),
		Steps: []resource.TestStep{
			{
				Config:      testAccKongConsumerConfig_basic("", ""),
				ExpectError: regexp.MustCompile("At least one of username or custom_id must be supplied"),
			},
		},
	})
}

func TestAccKongConsumer_update(t *testing.T) {
	var consumer kong.KongConsumer
	var updated kong.KongConsumer
	username := "kong-provider-acc-test-" + acctest.RandString(5)
	customId := "kong-provider-acc-test-" + acctest.RandString(5)

	updatedUsername := "kong-provider-acc-test-update-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongConsumerDestroy("kong_consumer.test"),
		Steps: []resource.TestStep{
			{
				Config: testAccKongConsumerConfig_basic(username, customId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerExists("kong_consumer.test", &consumer),
					testAccKongCustomerAttributes(&kong.KongConsumer{
						Username: username,
						CustomId: customId,
					}, &consumer),
					resource.TestCheckResourceAttr("kong_consumer.test", "username", username),
					resource.TestCheckResourceAttr("kong_consumer.test", "custom_id", customId),
				),
			},
			{
				Config: testAccKongConsumerConfig_basic(updatedUsername, customId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerExists("kong_consumer.test", &updated),
					testAccKongCustomerAttributes(&kong.KongConsumer{
						Username: updatedUsername,
						CustomId: customId,
					}, &updated),
					resource.TestCheckResourceAttr("kong_consumer.test", "username", updatedUsername),
					resource.TestCheckResourceAttr("kong_consumer.test", "custom_id", customId),
				),
			},
			{
				Config: testAccKongConsumerConfig_basic(username, customId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerExists("kong_consumer.test", &updated),
					testAccKongCustomerAttributes(&kong.KongConsumer{
						Username: username,
						CustomId: customId,
					}, &updated),
					resource.TestCheckResourceAttr("kong_consumer.test", "username", username),
					resource.TestCheckResourceAttr("kong_consumer.test", "custom_id", customId),
				),
			},
		},
	})
}

func testAccCheckKongConsumerDestroy(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := testAccProvider.Meta().(*kong.KongClient)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "kong_consumer" {
				continue
			}

			consumer, err := client.GetConsumer(state.RootModule().Resources[name].Primary.ID)

			if err == nil {
				return fmt.Errorf("consumer still exists: %s", consumer.Id)
			}

			if resourceDoesNotExistError(err) {
				return nil
			}

			return err
		}

		return nil
	}
}

func testAccCheckKongConsumerExists(name string, output *kong.KongConsumer) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		r, ok := state.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("consumer resource not found %s", name)
		}

		if r.Primary.ID == "" {
			return fmt.Errorf("no id set for %s", name)
		}

		client := testAccProvider.Meta().(*kong.KongClient)

		consumer, err := client.GetConsumer(r.Primary.ID)

		if err != nil {
			return err
		}

		*output = *consumer

		return nil
	}
}

func testAccKongCustomerAttributes(expected, actual *kong.KongConsumer) resource.TestCheckFunc {
	return func(_ *terraform.State) error {

		if expected.Username != actual.Username {
			return test_util.ExpectedAndActualError("Usernames don't match", expected.Username, actual.Username)
		}

		if expected.CustomId != actual.CustomId {
			return test_util.ExpectedAndActualError("Custom ids don't match", expected.CustomId, actual.CustomId)
		}

		return nil
	}
}

func testAccKongConsumerConfig_basic(username, customId string) string {
	return fmt.Sprintf(`
resource "kong_consumer" "test" {
	username = "%s"
	custom_id = "%s"
}
`, username, customId)
}
