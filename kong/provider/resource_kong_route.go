package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/hashicorp/terraform/helper/schema"
	"sort"
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
				Description: "Protocols that Kong will proxy to this route",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DefaultFunc: func() (interface{}, error) {
					return supportedProtocols, nil
				},
				Optional: true,
				Computed: true,
				Set:      schema.HashString,
			},
			"methods": {
				Description: "HTTP verbs that Kong will proxy to this route",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Set:      schema.HashString,
			},
			"hosts": {
				Description: "Host header values that should be matched to this route.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Set:      schema.HashString,
			},
			"paths": {
				Description: "List of path prefixes that will match this route",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Set:      schema.HashString,
			},
			"strip_path": {
				Description: " If the route is matched by path, this flag indicates whether the matched path should be removed from the upstream request.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"preserve_host": {
				Description: " If the route is matched by the `Host` header, this flag indicates if the `Host` header should be set to the matched value.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"service_id": {
				Description: "Unique identifier of the associated service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"created_at": {
				Description: "Unix timestamp of creation date.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"updated_at": {
				Description: "Unix timestamp of last edit date.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"regex_priority": {
				Description: "Determines the order that paths defined by regexes are evaluated.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
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

	err := kongClient.DeleteRoute(data.Id())

	if resourceDoesNotExistError(err) {
		return nil
	}

	return err
}

func mapToApiRoute(data *schema.ResourceData) *kong.KongRoute {
	return &kong.KongRoute{
		Id:           data.Id(),
		Methods:      toStringArray(data.Get("methods").(*schema.Set).List()),
		Protocols:    toStringArray(data.Get("protocols").(*schema.Set).List()),
		Hosts:        toStringArray(data.Get("hosts").(*schema.Set).List()),
		Paths:        toStringArray(data.Get("paths").(*schema.Set).List()),
		StripPath:    data.Get("strip_path").(bool),
		PreserveHost: data.Get("preserve_host").(bool),
		Service: kong.KongServiceReference{
			Id: data.Get("service_id").(string),
		},
	}
}

func validateRouteResource(data *schema.ResourceData) error {
	protocols := data.Get("protocols").(*schema.Set).List()
	invalidProtocols, _ := filterBySet(protocols, supportedProtocols)

	if len(invalidProtocols) > 0 {
		sort.Slice(invalidProtocols, func(i, j int) bool {
			return invalidProtocols[i] < invalidProtocols[j]
		})

		return fmt.Errorf("the supplied protocols are not supported by Kong: %s", strings.Join(invalidProtocols, ", "))
	}
	methods := data.Get("methods").(*schema.Set).List()

	invalidMethods, _ := filterBySet(methods, supportedMethods)

	if len(invalidMethods) > 0 {
		sort.Slice(invalidMethods, func(i, j int) bool {
			return invalidMethods[i] < invalidMethods[j]
		})

		return fmt.Errorf("invalid HTTP methods: %s", strings.Join(invalidMethods, ", "))
	}

	paths := data.Get("paths").(*schema.Set).List()
	hosts := data.Get("hosts").(*schema.Set).List()

	if len(paths) == 0 && len(hosts) == 0 && len(methods) == 0 {
		return fmt.Errorf("at least one of methods, paths, or hosts must be set in order for Kong to proxy traffic to this route")
	}

	return nil
}
