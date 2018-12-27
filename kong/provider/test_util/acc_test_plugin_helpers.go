package test_util

import (
	"encoding/json"
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCheckGenericKongPluginDestroy(provider *schema.Provider, resourceType, resourceName, pluginName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		kong := provider.Meta().(*kong.KongClient)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			plugin, err := kong.GetPlugin(state.RootModule().Resources[resourceName].Primary.ID)

			if err == nil {
				return fmt.Errorf("kong plugin %s still exists: %s", pluginName, plugin.Id)
			}

			if ResourceDoesNotExistError(err) {
				return nil
			}

			return err
		}

		return nil
	}
}

func TestAccCheckKongPluginExists(provider *schema.Provider, resourceName string, output *kong.KongPlugin) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		r, ok := state.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Plugin resource not found: %s", resourceName)
		}

		if r.Primary.ID == "" {
			return fmt.Errorf("No id set for %s", resourceName)
		}

		kong := provider.Meta().(*kong.KongClient)

		plugin, err := kong.GetPlugin(r.Primary.ID)

		if err != nil {
			return err
		}

		*output = *plugin

		return nil
	}
}

func TestAccKongPluginConfigAttributes(actualPlugin, expectedPlugin *kong.KongPlugin) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		expectedConfigJson, _ := json.Marshal(expectedPlugin.Config)
		actualConfigJson, _ := json.Marshal(actualPlugin.Config)

		if string(expectedConfigJson[:]) != string(actualConfigJson[:]) {
			return ExpectedAndActualError("Plugin configs differ", string(expectedConfigJson[:]), string(actualConfigJson[:]))
		}

		return nil
	}
}
