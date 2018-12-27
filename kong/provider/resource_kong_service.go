package provider

import (
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKongService() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongServiceCreate,
		Read:   resourceKongServiceRead,
		Update: resourceKongServiceUpdate,
		Delete: resourceKongServiceDelete,
		Importer: &schema.ResourceImporter{
			State: importResourceIfUuidIsValid,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The service name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"url": {
				Description: "The url for the service. It encapsulates protocol, host, port and path",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validateUrl,
			},
			"created_at": {
				Description: "Unix timestamp representing the time the service was created.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"updated_at": {
				Description: "Unix timestamp for the last time the service was updated.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"connect_timeout": {
				Description: "Time in milliseconds to connect to the upstream server.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60000,
			},
			"retries": {
				Description: "Number of times Kong will try to proxy if there's an error.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
			},
			"read_timeout": {
				Description: "Time in milliseconds between two read operations to the upstream server.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60000,
			},
			"write_timeout": {
				Description: "Time in milliseconds between two successive write operations to the upstream server.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60000,
			},
		},
	}
}

func resourceKongServiceCreate(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	kongService := mapToApiClient(data)

	service, err := kongClient.CreateService(&kongService)

	if err != nil {
		return err
	}

	data.SetId(service.Id)

	return resourceKongServiceRead(data, meta)
}

func resourceKongServiceRead(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	service, err := kongClient.GetService(data.Id())

	if err != nil {
		if resourceDoesNotExistError(err) {
			data.SetId("")

			return nil
		}

		return err
	}

	data.Set("name", service.Name)
	data.Set("url", service.Url)
	data.Set("created_at", service.CreatedAt)
	data.Set("updated_at", service.UpdatedAt)
	data.Set("connect_timeout", service.ConnectTimeout)
	data.Set("retries", service.Retries)
	data.Set("read_timeout", service.ReadTimeout)
	data.Set("write_timeout", service.WriteTimeout)

	return nil
}

func resourceKongServiceUpdate(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	serviceToUpdate := mapToApiClient(data)
	err := kongClient.UpdateService(&serviceToUpdate)

	if err != nil {
		return err
	}

	return resourceKongServiceRead(data, meta)
}

func resourceKongServiceDelete(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*kong.KongClient)

	err := kongClient.DeleteService(data.Id())
	if resourceDoesNotExistError(err) {
		return nil
	}

	return err
}

func mapToApiClient(data *schema.ResourceData) kong.KongService {
	return kong.KongService{
		Id:             data.Id(),
		Name:           data.Get("name").(string),
		Url:            data.Get("url").(string),
		ConnectTimeout: data.Get("connect_timeout").(int),
		Retries:        data.Get("retries").(int),
		ReadTimeout:    data.Get("read_timeout").(int),
		WriteTimeout:   data.Get("write_timeout").(int),
	}
}
