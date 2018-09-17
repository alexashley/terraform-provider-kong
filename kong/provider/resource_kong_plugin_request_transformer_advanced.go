package provider

import (
	"github.com/alexashley/terraform-provider-kong/kong/client"
	"github.com/alexashley/terraform-provider-kong/kong/util"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

type RequestTransformerCrud struct {
	Headers     []string `json:"headers,omitempty"`
	Querystring []string `json:"querystring,omitempty"`
	Body        []string `json:"body,omitempty"`
	Uri         string   `json:"uri,omitempty"`
}

type RequestTransformerConfig struct {
	HttpMethod string                 `json:"http_method,omitempty"`
	Remove     RequestTransformerCrud `json:"remove,omitempty"`
	Replace    RequestTransformerCrud `json:"replace,omitempty"`
	Rename     RequestTransformerCrud `json:"rename,omitempty"`
	Add        RequestTransformerCrud `json:"add,omitempty"`
	Append     RequestTransformerCrud `json:"append,omitempty"`
}

func resourceKongPluginRequestTransformerAdvanced() *schema.Resource {
	return CreateGenericPluginResource(&GenericPluginResource{
		Name: "request-transformer-advanced",
		AdditionalSchema: map[string]*schema.Schema{
			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GET",
					"PUT",
					"POST",
					"DELETE",
					"PATCH",
					"HEAD",
					"TRACE",
					"CONNECT",
					"OPTIONS",
				}, true),
			},
			"remove_headers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"remove_querystring": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"remove_body_params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"replace_headers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"replace_querystring": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"replace_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"replace_body_params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rename_headers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rename_querystring": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rename_body_params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"add_headers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"add_querystring": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"add_body_params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"append_headers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"append_querystring": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"append_body_params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		MapSchemaToPluginConfig: func(data *schema.ResourceData) interface{} {
			return RequestTransformerConfig{
				HttpMethod: data.Get("http_method").(string),
				Remove: RequestTransformerCrud{
					Headers:     getStringArray(data, "remove_headers"),
					Querystring: getStringArray(data, "remove_querystring"),
					Body:        getStringArray(data, "remove_body_params"),
				},
				Replace: RequestTransformerCrud{
					Headers:     getStringArray(data, "replace_headers"),
					Querystring: getStringArray(data, "replace_querystring"),
					Body:        getStringArray(data, "replace_body_params"),
					Uri:         data.Get("replace_uri").(string),
				},
				Rename: RequestTransformerCrud{
					Headers:     getStringArray(data, "rename_headers"),
					Querystring: getStringArray(data, "rename_querystring"),
					Body:        getStringArray(data, "rename_body_params"),
				},
				Add: RequestTransformerCrud{
					Headers:     getStringArray(data, "add_headers"),
					Querystring: getStringArray(data, "add_querystring"),
					Body:        getStringArray(data, "add_body_params"),
				},
				Append: RequestTransformerCrud{
					Headers:     getStringArray(data, "append_headers"),
					Querystring: getStringArray(data, "append_querystring"),
					Body:        getStringArray(data, "append_body_params"),
				},
			}
		},
		MapApiModelToResource: func(plugin *client.KongPlugin, data *schema.ResourceData) {
			pluginConfig := plugin.Config.(map[string]interface{})

			httpMethod, ok := data.GetOk("http_method")

			if ok {
				data.Set("http_method", httpMethod.(string))
			}

			prefixes := []string{"remove", "replace", "rename", "add", "append"}
			for i := range prefixes {
				prefix := prefixes[i]
				action := pluginConfig[prefix].(map[string]interface{})
				for key, value := range action {
					// Kong returns an empty {} when the array is empty
					// this is an inelegant workaround
					if _, ok = value.(map[string]interface{}); ok {
						value = &[]string{}
					}

					tfKeyName := prefix + "_" + key

					if key == "body" {
						tfKeyName = tfKeyName + "_params"
					}

					err := data.Set(tfKeyName, value)

					if err != nil {
						util.Log("Error setting thing: " + err.Error())
					}
				}
			}
		},
	})
}

func getStringArray(data *schema.ResourceData, keyName string) []string {
	value, ok := data.GetOk(keyName)

	if !ok {
		return []string{}
	}

	return toStringArrayFromInterface(value)
}
