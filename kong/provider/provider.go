package provider

import (
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/hashicorp/terraform/helper/schema"
)

var resourcesMap = map[string]*schema.Resource{
	"kong_service":  resourceKongService(),
	"kong_route":    resourceKongRoute(),
	"kong_plugin":   resourceKongPlugin(),
	"kong_consumer": resourceKongConsumer(),
}

// these are separate from the other resources, since `kong_plugin` needs a way to determine if a specific resource implementation of a plugin exists.
// if it does, `kong_plugin` will error and force the use of the specialized resource
var pluginResourcesMap = map[string]*schema.Resource{
	"kong_plugin_request_transformer_advanced": resourceKongPluginRequestTransformerAdvanced(),
	"kong_plugin_openid_connect":               resourceKongPluginOpenidConnect(),
}

func KongProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: mergeResourceMaps(resourcesMap, pluginResourcesMap),
		Schema: map[string]*schema.Schema{
			"admin_api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KONG_ADMIN_API_URL", nil),
				Description: "Kong's admin api url. Usually bound to port 8001.",
			},
			"rbac_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KONG_RBAC_TOKEN", nil),
				Description: "An RBAC user's token. If your Kong EE installation uses RBAC, then this is required for the provider to interact with the admin API.",
			},
		},
		ConfigureFunc: func(data *schema.ResourceData) (interface{}, error) {
			adminApiUrl := data.Get("admin_api_url").(string)
			rbacToken := data.Get("rbac_token").(string)

			return kong.NewKongClient(kong.KongConfig{AdminApiUrl: adminApiUrl, RbacToken: rbacToken})
		},
	}
}

func mergeResourceMaps(a, b map[string]*schema.Resource) map[string]*schema.Resource {
	combined := make(map[string]*schema.Resource)

	for resourceName, resourceValue := range a {
		combined[resourceName] = resourceValue
	}

	for resourceName, resourceValue := range b {
		combined[resourceName] = resourceValue
	}

	return combined
}
