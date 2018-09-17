package provider

import (
	"github.com/alexashley/terraform-provider-kong/kong/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKongPluginIpHeaderRestriction() *schema.Resource {
	return CreateGenericPluginResource(&GenericPluginResource{
		Name: "ip-header-restriction",
		AdditionalSchema: map[string]*schema.Schema{
			"whitelist": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"ip_header": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "x-forwarded-for",
			},
			"override_global": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		MapSchemaToPluginConfig: func(data *schema.ResourceData) interface{} {
			return map[string]interface{}{
				"whitelist":       data.Get("whitelist"),
				"ip_header":       data.Get("ip_header"),
				"override_global": data.Get("override_global"),
			}
		},
		MapApiModelToResource: func(plugin *client.KongPlugin, data *schema.ResourceData) {
			pluginConfig := plugin.Config.(map[string]interface{})

			data.Set("whitelist", toStringArray(pluginConfig["whitelist"].([]interface{})))
			data.Set("ip_header", pluginConfig["ip_header"].(string))
			data.Set("override_global", pluginConfig["override_global"].(bool))
		},
	})
}
