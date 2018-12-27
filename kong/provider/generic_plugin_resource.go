package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/hashicorp/terraform/helper/schema"
)

type ResourceMapper func(plugin *kong.KongPlugin, data *schema.ResourceData) error
type ConfigMapper func(data *schema.ResourceData) (interface{}, error)
type GetName func(data *schema.ResourceData) string

type GenericPluginResource struct {
	AllowsConsumers         bool
	AdditionalSchema        map[string]*schema.Schema
	Name                    GetName
	MapSchemaToPluginConfig ConfigMapper
	MapApiModelToResource   ResourceMapper
}

var defaultPluginSchema = map[string]*schema.Schema{
	"service_id": {
		Description: "Unique identifier of the associated service.",
		Type:        schema.TypeString,
		Optional:    true,
	},
	"route_id": {
		Description: "Unique identifier of the associated route.",
		Type:        schema.TypeString,
		Optional:    true,
	},
	"consumer_id": {
		Description: "Unique identifier of the consumer for which this plugin will run. Not all plugins allow consumers",
		Type:        schema.TypeString,
		Optional:    true,
	},
	"enabled": {
		Description: "Toggle whether the plugin will run",
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
	},
	"created_at": {
		Description: "Unix timestamp representing when the plugin was created.",
		Type:        schema.TypeInt,
		Computed:    true,
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
		Schema: mergeSchemaWithDefaults(resource.AdditionalSchema, resource.AllowsConsumers),
	}
}

func mergeSchemaWithDefaults(pluginSchema map[string]*schema.Schema, allowsConsumers bool) map[string]*schema.Schema {
	mergedSchema := make(map[string]*schema.Schema)

	for key, value := range defaultPluginSchema {
		mergedSchema[key] = value
	}

	for key, value := range pluginSchema {
		mergedSchema[key] = value
	}

	if !allowsConsumers {
		delete(mergedSchema, "consumer_id")
	}

	return mergedSchema
}

func (resource *GenericPluginResource) resourceGenericPluginCreate(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	payload, err := resource.mapToApiModel(data)

	if err != nil {
		return err
	}

	plugin, err := kongClient.CreatePlugin(payload)

	if err != nil {
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

	if resource.AllowsConsumers {
		data.Set("consumer_id", plugin.ConsumerId)
	}

	data.Set("name", plugin.Name)
	data.Set("enabled", plugin.Enabled)
	data.Set("created_at", plugin.CreatedAt)

	return resource.MapApiModelToResource(plugin, data)
}

func (resource *GenericPluginResource) resourceGenericPluginUpdate(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	pluginToUpdate, err := resource.mapToApiModel(data)

	if err != nil {
		return err
	}

	err = kongClient.UpdatePlugin(pluginToUpdate)

	if err != nil {
		return err
	}

	return resource.resourceGenericPluginRead(data, meta)
}

func (resource *GenericPluginResource) resourceGenericPluginDelete(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	err := kongClient.DeletePlugin(data.Id())

	if resourceDoesNotExistError(err) {
		return nil
	}

	return err
}

func (resource *GenericPluginResource) mapToApiModel(data *schema.ResourceData) (*kong.KongPlugin, error) {
	config, err := resource.MapSchemaToPluginConfig(data)

	if err != nil {
		return nil, fmt.Errorf("error mapping TF resource to API model: %s", err)
	}

	plugin := kong.KongPlugin{
		Id:        data.Id(),
		ServiceId: data.Get("service_id").(string),
		RouteId:   data.Get("route_id").(string),
		Name:      resource.Name(data),
		Enabled:   data.Get("enabled").(bool),
		Config:    config,
	}

	if resource.AllowsConsumers {
		plugin.ConsumerId = data.Get("consumer_id").(string)
	}

	return &plugin, nil
}
