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
	return &schema.Resource{
		Create: resourceKongPluginCreate,
		Read:   resourceKongPluginRead,
		Update: resourceKongPluginUpdate,
		Delete: resourceKongPluginDelete,
		Importer: &schema.ResourceImporter{
			State: importResourceIfUuidIsValid, // TODO: change import to always import config_json over config
		},
		Schema: map[string]*schema.Schema{
			"service_id": {
				Description: "The service for which the plugin will run.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"route_id": {
				Description: "The route for which the plugin will run.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"consumer_id": {
				Description: "The consumer for which the plugin will run. Not supported by all plugins.",
				Type:        schema.TypeString,
				Optional:    true,
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
						errorMessage := fmt.Errorf("Plugin %s has a resource implementation: %s. This resource should be used instead.", name, pluginResourceName)

						errors = append(errors, errorMessage)
					}

					return warnings, errors
				},
			},
			"config": {
				Description:   "An object representing the plugin's configuration. At this time it's not possible to represent all valid plugin configurations with Terraform. Should this be a problem, you can use a specific plugin resource or the `config_json` field.	",
				Type:          schema.TypeMap,
				Elem:          &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:      true,
				ConflictsWith: []string{"config_json"},
			},
			"config_json": {
				Description:      "A JSON string containing the plugin's configuration. Can't be used with `config`.",
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"config"},
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},
			"enabled": {
				Description: "Turns the plugin on or off.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"created_at": {
				Description: "Unix timestamp representing the creation date",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceKongPluginCreate(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	config, err := getPluginConfig(data)

	if err != nil {
		return err
	}

	plugin, err := kongClient.CreatePlugin(&kong.KongPlugin{
		ServiceId:  data.Get("service_id").(string),
		RouteId:    data.Get("route_id").(string),
		ConsumerId: data.Get("consumer_id").(string),
		Name:       data.Get("name").(string),
		Config:     config,
		Enabled:    data.Get("enabled").(bool),
	})

	if err != nil {
		return err
	}

	data.SetId(plugin.Id)

	return resourceKongPluginRead(data, meta)
}

func resourceKongPluginRead(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	plugin, err := kongClient.GetPlugin(data.Id())

	if err != nil {
		if resourceDoesNotExistError(err) {
			data.SetId("")

			return nil
		}

		return err
	}

	data.Set("service_id", plugin.ServiceId)
	data.Set("route_id", plugin.RouteId)
	data.Set("consumer_id", plugin.ConsumerId)
	data.Set("name", plugin.Name)
	data.Set("enabled", plugin.Enabled)
	data.Set("created_at", plugin.CreatedAt)

	if data.Get("config_json").(string) != "" {
		configJson, _ := json.Marshal(plugin.Config)
		data.Set("config_json", string(configJson[:]))
	} else {
		configErr := data.Set("config", plugin.Config)

		if configErr != nil {
			// TF schema can only handle simple maps.
			// This is mainly for plugins where the **default** config cannot be created as map[string]string
			// For example, basic-auth has a `hide_credentials` flag, which cannot be converted to a string
			data.SetId("")

			return fmt.Errorf("%s; use the config_json field for more complex configurations; you may have to import the resource by-hand or delete and re-create using the config_json field", configErr)
		}
	}

	return nil
}

func resourceKongPluginUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongPluginDelete(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	return kongClient.DeletePlugin(data.Id())
}

func getPluginConfig(data *schema.ResourceData) (map[string]interface{}, error) {
	var config map[string]interface{}

	config = data.Get("config").(map[string]interface{})
	configJson := data.Get("config_json").(string)

	if len(config) != 0 {
		return config, nil
	}

	err := json.Unmarshal([]byte(configJson), &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func pluginIsAResource(pluginName string) (string, bool) {
	snakeCaseName := strings.Replace(pluginName, "-", "_", -1)
	resourceName := "kong_plugin_" + snakeCaseName

	_, ok := pluginResourcesMap[resourceName]

	return resourceName, ok
}
