package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func KongProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"kong_service": resourceKongService(),
		},
		Schema: map[string]*schema.Schema{
			"admin_api_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Kong's admin api url.",
			},
		},
		ConfigureFunc: func(data *schema.ResourceData) (interface{}, error) {
			adminApiUrl := data.Get("admin_api_url").(string)

			return NewKongClient(adminApiUrl)
		},
	}
}
