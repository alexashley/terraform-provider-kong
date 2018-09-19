package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

var supportedProtocols = []string{"http", "https"}
var supportedMethods = []string{
	"GET",
	"PUT",
	"POST",
	"DELETE",
	"PATCH",
	"HEAD",
	"OPTIONS",
	"TRACE",
	"CONNECT",
}

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
			"protocols": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DefaultFunc: func() (interface{}, error) {
					return supportedProtocols, nil
				},
				Optional: true,
				Computed: true,
			},
			"methods": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"hosts": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"paths": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"strip_path": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"preserve_host": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"regex_priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceKongRouteCreate(data *schema.ResourceData, meta interface{}) error {
	if err := validateRouteResource(data); err != nil {
		return err
	}

	kongClient := meta.(*kong.KongClient)

	route, err := kongClient.CreateRoute(mapToApiRoute(data))

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

	err := kongClient.UpdateRoute(mapToApiRoute(data))

	if err != nil {
		return err
	}

	return resourceKongRouteRead(data, meta)
}

func resourceKongRouteDelete(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	return kongClient.DeleteRoute(data.Id())
}

func mapToApiRoute(data *schema.ResourceData) *kong.KongRoute {
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

func validateRouteResource(data *schema.ResourceData) error {
	protocols := data.Get("protocols")
	invalidProtocols, _ := filterBySet(protocols, supportedProtocols)

	if len(invalidProtocols) > 0 {
		return fmt.Errorf("The supplied protocols are not supported by Kong: %s", strings.Join(invalidProtocols, ", "))
	}
	methods := data.Get("methods").([]interface{})

	invalidMethods, _ := filterBySet(methods, supportedMethods)

	if len(invalidMethods) > 0 {
		return fmt.Errorf("Invalid HTTP methods: %s", strings.Join(invalidMethods, ", "))
	}

	paths := data.Get("paths").([]interface{})
	hosts := data.Get("hosts").([]interface{})

	if len(paths) == 0 && len(hosts) == 0 && len(methods) == 0 {
		return fmt.Errorf("At least one of methods, paths, or hosts must be set in order for Kong to proxy traffic to this route.")
	}

	return nil
}
