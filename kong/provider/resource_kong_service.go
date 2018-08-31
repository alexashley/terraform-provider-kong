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
				Required: false,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
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
	}

	id, err := kongClient.CreateService(kongService)

	if err != nil {
		return err
	}

	d.SetId(id)

	return nil
}

func resourceKongServiceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKongServiceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
