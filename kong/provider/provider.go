package provider

import (
	"github.com/alexashley/terraform-provider-kong/kong/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func KongProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"kong_service": resourceKongService(),
			"kong_route":   resourceKongRoute(),
			"kong_plugin":  resourceKongPlugin(),
		},
		Schema: map[string]*schema.Schema{
			"admin_api_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Kong's admin api url. Usually bound to port 8001.",
			},
			"rbac_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An RBAC user's token. If your Kong EE installation uses RBAC, then this is required for the provider to interact with the admin API.",
			},
		},
		ConfigureFunc: func(data *schema.ResourceData) (interface{}, error) {
			adminApiUrl := data.Get("admin_api_url").(string)
			rbacToken := data.Get("rbac_token").(string)

			return client.NewKongClient(client.KongConfig{AdminApiUrl: adminApiUrl, RbacToken: rbacToken})
		},
	}
}
