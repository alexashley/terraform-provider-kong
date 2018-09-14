package provider

import (
	"github.com/alexashley/terraform-provider-kong/kong/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKongPluginIpHeaderRestriction() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongPluginIpHeaderRestrictionCreate,
		Read:   resourceKongPluginIpHeaderRestrictionRead,
		Update: resourceKongPluginIpHeaderRestrictionUpdate,
		Delete: resourceKongPluginIpHeaderRestrictionDelete,
		Schema: map[string]*schema.Schema{
			"service_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"consumer_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"whitelist": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"ip_header": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "x-forwarded-for",
			},
			"override_global": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceKongPluginIpHeaderRestrictionCreate(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)
	plugin, err := kongClient.CreatePlugin(client.KongPlugin{
		ServiceId:  data.Get("service_id").(string),
		RouteId:    data.Get("route_id").(string),
		ConsumerId: data.Get("consumer_id").(string),
		Name:       "ip-header-restriction",
		Enabled:    data.Get("enabled").(bool),
		Config: map[string]interface{}{
			"whitelist":       data.Get("whitelist"),
			"ip_header":       data.Get("ip_header"),
			"override_global": data.Get("override_global"),
		},
	})

	if err != nil {
		return nil
	}

	data.SetId(plugin.Id)

	return resourceKongPluginIpHeaderRestrictionRead(data, meta)
}

func resourceKongPluginIpHeaderRestrictionRead(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	plugin, err := kongClient.GetPlugin(data.Id())

	if err != nil {
		data.SetId("")
		return nil
	}

	data.Set("service_id", plugin.ServiceId)
	data.Set("route_id", plugin.RouteId)
	data.Set("consumer_id", plugin.ConsumerId)
	data.Set("name", plugin.Name)
	data.Set("enabled", plugin.Enabled)
	data.Set("created_at", plugin.CreatedAt)
	data.Set("whitelist", toStringArray(plugin.Config["whitelist"].([]interface{})))
	data.Set("ip_header", plugin.Config["ip_header"].(string))
	data.Set("override_global", plugin.Config["override_global"].(bool))

	return nil
}

func resourceKongPluginIpHeaderRestrictionUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongPluginIpHeaderRestrictionDelete(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	return kongClient.DeletePlugin(data.Id())
}
