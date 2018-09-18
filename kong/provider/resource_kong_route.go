package provider

import (
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKongRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongRouteCreate,
		Read:   resourceKongRouteRead,
		Update: resourceKongRouteUpdate,
		Delete: resourceKongRouteDelete,
		Importer: &schema.ResourceImporter{
			State: importResourceIfUuidIsValid,
		},
		Schema: map[string]*schema.Schema{
			"protocols": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DefaultFunc: func() (interface{}, error) {
					return []string{"http", "https"}, nil
				},
				Optional: true,
				Computed: true,
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
			"service_id": &schema.Schema{
				Type:     schema.TypeString,
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
	kongClient := meta.(*kong.KongClient)

	route, err := kongClient.CreateRoute(mapToApi(data))

	if err != nil {
		return err
	}

	data.SetId(route.Id)

	return resourceKongRouteRead(data, meta)
}

func resourceKongRouteRead(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	route, err := kongClient.GetRoute(data.Id())

	if err != nil {
		if resourceDoesNotExistError(err) {
			data.SetId("")

			return nil
		}

		return err
	}

	data.Set("hosts", route.Hosts)
	data.Set("protocols", route.Protocols)
	data.Set("methods", route.Methods)
	data.Set("paths", route.Paths)
	data.Set("strip_path", route.StripPath)
	data.Set("preserve_host", route.PreserveHost)
	data.Set("service_id", route.Service.Id)
	data.Set("created_at", route.CreatedAt)
	data.Set("updated_at", route.UpdatedAt)
	data.Set("regex_priority", route.RegexPriority)

	return nil
}

func resourceKongRouteUpdate(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	err := kongClient.UpdateRoute(mapToApi(data))

	if err != nil {
		return err
	}

	return resourceKongRouteRead(data, meta)
}

func resourceKongRouteDelete(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	return kongClient.DeleteRoute(data.Id())
}

func mapToApi(data *schema.ResourceData) *kong.KongRoute {
	return &kong.KongRoute{
		Id:           data.Id(),
		Methods:      toStringArray(data.Get("methods").([]interface{})),
		Protocols:    toStringArray(data.Get("protocols").([]interface{})),
		Hosts:        toStringArray(data.Get("hosts").([]interface{})),
		Paths:        toStringArray(data.Get("paths").([]interface{})),
		StripPath:    data.Get("strip_path").(bool),
		PreserveHost: data.Get("preserve_host").(bool),
		Service: kong.KongServiceReference{
			Id: data.Get("service_id").(string),
		},
	}
}
