package provider

import (
	"encoding/json"
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
	"os"
	"strings"
)

func resourceKongPlugin() *schema.Resource {
	return CreateGenericPluginResource(&GenericPluginResource{
		AllowsConsumers: true,
		Name: func(data *schema.ResourceData) string {
			return data.Get("name").(string)
		},
		AdditionalSchema: map[string]*schema.Schema{
			"config_json": {
				Description:      "A JSON string containing the plugin's configuration.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},
			"name": {
				Description: "The plugin name, e.g. `basic-auth`",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: func(i interface{}, s string) (warnings []string, errors []error) {
					name := i.(string)

					pluginResourceName, ok := pluginIsAResource(name)

					escapeHatchEnvName := "TF_KONG_ALLOW_GENERIC_PLUGIN_" + strings.ToUpper(strings.Replace(name, "-", "_", -1))

					if ok && os.Getenv(escapeHatchEnvName) == "" {
						errorMessage := fmt.Errorf("plugin %s has a resource implementation: %s. this resource should be used instead", name, pluginResourceName)

						errors = append(errors, errorMessage)
					}

					return warnings, errors
				},
			},
		},
		MapSchemaToPluginConfig: func(data *schema.ResourceData) (interface{}, error) {
			configJson := data.Get("config_json").(string)
			var config map[string]interface{}
			if err := json.Unmarshal([]byte(configJson), &config); err != nil {
				return nil, err
			}

			return config, nil
		},
		MapApiModelToResource: func(plugin *kong.KongPlugin, data *schema.ResourceData) error {
			config := plugin.Config
			configJson, err := json.Marshal(config)

			if err != nil {
				return err
			}

			data.Set("config_json", string(configJson[:]))
			data.Set("name", plugin.Name)

			return nil
		},
	})
}

func pluginIsAResource(pluginName string) (string, bool) {
	snakeCaseName := strings.Replace(pluginName, "-", "_", -1)
	resourceName := "kong_plugin_" + snakeCaseName

	_, ok := pluginResourcesMap[resourceName]

	return resourceName, ok
}
