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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
		Name: data.Get("name").(string),
		Url:  data.Get("url").(string),
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
