package provider

import (
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/alexashley/terraform-provider-kong/kong/util"
	"github.com/hashicorp/terraform/helper/schema"
)

type ResourceMapper func(plugin *kong.KongPlugin, data *schema.ResourceData)
type ConfigMapper func(data *schema.ResourceData) interface{}

type GenericPluginResource struct {
	Name                    string
	AdditionalSchema        map[string]*schema.Schema
	MapSchemaToPluginConfig ConfigMapper
	MapApiModelToResource   ResourceMapper
}

var defaultPluginSchema = map[string]*schema.Schema{
	"service_id": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"route_id": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"consumer_id": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"created_at": {
		Type:     schema.TypeInt,
		Computed: true,
	},
}

func CreateGenericPluginResource(resource *GenericPluginResource) *schema.Resource {
	return &schema.Resource{
		Create: resource.resourceGenericPluginCreate,
		Read:   resource.resourceGenericPluginRead,
		Update: resource.resourceGenericPluginUpdate,
		Delete: resource.resourceGenericPluginDelete,
		Importer: &schema.ResourceImporter{
			State: importResourceIfUuidIsValid,
		},
		Schema: mergeSchemaWithDefaults(resource.AdditionalSchema),
	}
}

func mergeSchemaWithDefaults(pluginSchema map[string]*schema.Schema) map[string]*schema.Schema {
	mergedSchema := make(map[string]*schema.Schema)

	for key, value := range defaultPluginSchema {
		mergedSchema[key] = value
	}

	for key, value := range pluginSchema {
		mergedSchema[key] = value
	}

	return mergedSchema
}

func (resource *GenericPluginResource) resourceGenericPluginCreate(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)
	plugin, err := kongClient.CreatePlugin(&kong.KongPlugin{
		ServiceId:  data.Get("service_id").(string),
		RouteId:    data.Get("route_id").(string),
		ConsumerId: data.Get("consumer_id").(string),
		Name:       resource.Name,
		Enabled:    data.Get("enabled").(bool),
		Config:     resource.MapSchemaToPluginConfig(data),
	})

	if err != nil {
		util.Log("ERROR " + err.Error())
		return err
	}

	data.SetId(plugin.Id)

	return resource.resourceGenericPluginRead(data, meta)
}

func (resource *GenericPluginResource) resourceGenericPluginRead(data *schema.ResourceData, meta interface{}) error {
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

	resource.MapApiModelToResource(plugin, data)

	return nil
}

func (resource *GenericPluginResource) resourceGenericPluginUpdate(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	pluginToUpdate := resource.mapToApiModel(data)

	err := kongClient.UpdatePlugin(pluginToUpdate)

	if err != nil {
		return err
	}

	return resource.resourceGenericPluginRead(data, meta)
}

func (resource *GenericPluginResource) resourceGenericPluginDelete(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	return kongClient.DeletePlugin(data.Id())
}

func (resource *GenericPluginResource) mapToApiModel(data *schema.ResourceData) *kong.KongPlugin {
	return &kong.KongPlugin{
		Id:         data.Id(),
		ServiceId:  data.Get("service_id").(string),
		RouteId:    data.Get("route_id").(string),
		ConsumerId: data.Get("consumer_id").(string),
		Name:       resource.Name,
		Enabled:    data.Get("enabled").(bool),
		Config:     resource.MapSchemaToPluginConfig(data),
	}
}
