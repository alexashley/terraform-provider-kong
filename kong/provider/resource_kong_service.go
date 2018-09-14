package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKongService() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongServiceCreate,
		Read:   resourceKongServiceRead,
		Update: resourceKongServiceUpdate,
		Delete: resourceKongServiceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKongServiceImport,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url": &schema.Schema{
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
			"connect_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60000,
			},
			"retries": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"read_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60000,
			},
			"write_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60000,
			},
		},
	}
}

func resourceKongServiceCreate(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	kongService := client.KongService{
		Name:           data.Get("name").(string),
		Url:            data.Get("url").(string),
		ConnectTimeout: data.Get("connect_timeout").(int),
		Retries:        data.Get("retries").(int),
		ReadTimeout:    data.Get("read_timeout").(int),
		WriteTimeout:   data.Get("write_timeout").(int),
	}

	service, err := kongClient.CreateService(kongService)

	if err != nil {
		return err
	}

	data.SetId(service.Id)

	return resourceKongServiceRead(data, meta)
}

func resourceKongServiceRead(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	service, err := kongClient.GetService(data.Id())

	if err != nil {
		httpError, ok := err.(*client.KongHttpError)

		if ok && httpError.StatusCode == 404 {
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
	return nil
}

func resourceKongServiceDelete(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	return kongClient.DeleteService(data.Id())
}

func resourceKongServiceImport(data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := data.Id()

	if !isValidUuid(id) {
		return nil, fmt.Errorf("%s is not a valid UUID", id)
	}

	return []*schema.ResourceData{data}, nil
}
