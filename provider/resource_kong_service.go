package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type KongService struct {
	id       string
	name     string
	protocol string
	host     string
	port     int
	path     string
}

func resourceKongService() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongServiceCreate,
		Read:   resourceKongServiceRead,
		Update: resourceKongServiceUpdate,
		Delete: resourceKongServiceDelete,

		/* // not here: url property
		{
			"id": "4e13f54a-bbf1-47a8-8777-255fed7116f2", (computed)
			"created_at": 1488869076800, (computed)
			"updated_at": 1488869076800, (computed)
			"connect_timeout": 60000, (optional)
			"protocol": "http", (required, one of (http (default), https)
			"host": "example.org", (required)
			"port": 80, (required, default)
			"path": "/api", (optional)
			"name": "example-service", (optional)
			"retries": 5, (optional, default)
			"read_timeout": 60000 (optional, default),
			"write_timeout": 60000 (optional, default)
		}
		*/

		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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

//resp, err := http.PostForm("http://example.com/form",
//url.Values{"key": {"Value"}, "id": {"123"}}

func resourceKongServiceCreate(d *schema.ResourceData, meta interface{}) error {
	kong := meta.(*Kong)

	kongService := KongService{
		name:     d.Get("name").(string),
		protocol: d.Get("protocol").(string),
		host:     d.Get("host").(string),
		port:     d.Get("port").(int),
		path:     d.Get("path").(string),
	}

	id, err := kong.createService(kongService)

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
