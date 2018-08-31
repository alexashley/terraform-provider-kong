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
				Default: "http",
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port": &schema.Schema{
				Default: 80,
				Type:     schema.TypeInt,
				Optional: true,
			},
			"url": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceKongServiceCreate(d *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	kongService := client.KongService{
		Name:     d.Get("name").(string),
		Protocol: d.Get("protocol").(string),
		Host:     d.Get("host").(string),
		Port:     d.Get("port").(int),
		Path:     d.Get("path").(string),
		Url:	  d.Get("url").(string),
	}

	id, err := kongClient.CreateService(kongService)

	if err != nil {
		return err
	}

	d.SetId(id)

	return nil
}

func resourceKongServiceRead(d *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	service, err := kongClient.GetService(d.Id())

	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("name", service.Name)
	d.Set("protocol", service.Protocol)
	d.Set("host", service.Host)
	d.Set("port", service.Port)
	d.Set("path", service.Path)
	d.Set("url", service.Url)

	return nil
}

func resourceKongServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongServiceDelete(d *schema.ResourceData, meta interface{}) error {
	kongClient := meta.(*client.KongClient)

	return kongClient.DeleteService(d.Id())
}
