package provider

import (
	"github.com/alexashley/terraform-provider-kong/kong/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKongService() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongServiceCreate,
		Read:   resourceKongServiceRead,
		Update: resourceKongServiceUpdate,
		Delete: resourceKongServiceDelete,
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": &schema.Schema{
				Default:  "",
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": &schema.Schema{
				Default:  "http",
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port": &schema.Schema{
				Default:  80,
				Type:     schema.TypeInt,
				Optional: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceKongServiceCreate(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	kongService := client.KongService{
		Name:     data.Get("name").(string),
		Protocol: data.Get("protocol").(string),
		Host:     data.Get("host").(string),
		Port:     data.Get("port").(int),
		Path:     data.Get("path").(string),
		Url:      data.Get("url").(string),
	}

	id, err := kongClient.CreateService(kongService)

	if err != nil {
		return err
	}

	data.SetId(id)

	return nil
}

func resourceKongServiceRead(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	service, err := kongClient.GetService(data.Id())

	if err != nil {
		data.SetId("")
		return nil
	}

	data.Set("name", service.Name)
	data.Set("protocol", service.Protocol)
	data.Set("host", service.Host)
	data.Set("port", service.Port)
	data.Set("path", service.Path)
	data.Set("url", service.Url)

	return nil
}

func resourceKongServiceUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongServiceDelete(data *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	return kongClient.DeleteService(data.Id())
}
