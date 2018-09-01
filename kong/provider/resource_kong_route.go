package provider

import (
	"github.com/alexashley/terraform-provider-kong/kong/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKongRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongRouteCreate,
		Read:   resourceKongRouteRead,
		Update: resourceKongRouteUpdate,
		Delete: resourceKongRouteDelete,
		Schema: map[string]*schema.Schema{
			"protocols": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"methods": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"hosts": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"paths": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"strip_path": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"preserve_host": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"service": &schema.Schema{
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"regex_priority": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceKongRouteCreate(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	associatedService := toStringMap(data.Get("service").(map[string]interface{}))
	route, err := kongClient.CreateRoute(client.KongRoute{
		Methods:      toStringArray(data.Get("methods").([]interface{})),
		Protocols:    toStringArray(data.Get("protocols").([]interface{})),
		Hosts:        toStringArray(data.Get("hosts").([]interface{})),
		Paths:        toStringArray(data.Get("paths").([]interface{})),
		StripPath:    data.Get("strip_path").(bool),
		PreserveHost: data.Get("preserve_host").(bool),
		Service: client.KongServiceReference{
			Id: associatedService["id"],
		},
	})

	if err != nil {
		return err
	}

	data.SetId(route.Id)
	data.Set("created_at", route.CreatedAt)
	data.Set("updated_at", route.UpdatedAt)

	return nil
}

func resourceKongRouteRead(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongRouteUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongRouteDelete(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	return kongClient.DeleteRoute(data.Id())
}

func toStringMap(data map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for key, value := range data {
		result[key] = value.(string)
	}

	return result
}

func toStringArray(data []interface{}) []string {
	result := make([]string, len(data))

	for index, value := range data {
		result[index] = value.(string)
	}

	return result
}
