package provider

import (
	"github.com/alexashley/terraform-provider-kong/kong/kong"
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
		AllowsConsumers: true,
		Name: func(data *schema.ResourceData) string {
			return "request-transformer-advanced"
		},
		AdditionalSchema: map[string]*schema.Schema{
			"http_method": {
				Description: "Method that will be used for the upstream request.",
				Type:        schema.TypeString,
				Optional:    true,
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
				Description: "Header key:value pairs to scrub from the request.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				// The hash function needs to be specified -- the default is to hash the entire schema
				// This is technically fine, but we don't really care about the schema, but rather the value
				// https://github.com/hashicorp/terraform/blob/c753df6a933ef25197101664b1d7cbb5d115701d/helper/schema/schema.go#L290
				Set: schema.HashString,
			},
			"remove_querystring": {
				Description: "Querystring key:value pairs to scrub from the request.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"remove_body_params": {
				Description: "Body parameters to scrub from the request.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"replace_headers": {
				Description: "Header key:value pairs. If the header is set, its value will be replaced. Otherwise it will be ignored",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"replace_querystring": {
				Description: "Querystring key:value pairs to replace if the key is set in the request.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"replace_uri": {
				Description: "Rewrites the path to the upstream request.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"replace_body_params": {
				Description: "Body parameters to replace in the request. If the param is set, its value will be replaced. Otherwise it will be ignored.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"rename_headers": {
				Description: "Header key:value pairs. If the header is set, it will be renamed. The value will remain unchanged.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"rename_querystring": {
				Description: "Querystring key:value pairs. If the querystring is in the request, the field will be renamed but the value will remain the same.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"rename_body_params": {
				Description: "Body parameters to rename in the request.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"add_headers": {
				Description: "Header key:value pairs to add to the request. Ignored if the header is already set.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"add_querystring": {
				Description: "Querystring key:value pairs to add to the request. Ignored if the query is already set.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"add_body_params": {
				Description: "Body parameters to add to the request. Ignored if already set.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"append_headers": {
				Description: "Header key:value pairs to append to the request. The header is added if it's not already present",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"append_querystring": {
				Description: "Querystring key:value pairs to append to the request. The query is added if it's not already present",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"append_body_params": {
				Description: "Body parameters to append to the request. The parameter is set if it's not already in the request",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
		},
		MapSchemaToPluginConfig: func(data *schema.ResourceData) (interface{}, error) {
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
			}, nil
		},
		MapApiModelToResource: func(plugin *kong.KongPlugin, data *schema.ResourceData) error {
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

			return nil
		},
	})
}

func getStringArray(data *schema.ResourceData, keyName string) []string {
	value, ok := data.GetOk(keyName)

	if !ok {
		return []string{}
	}

	set := value.(*schema.Set)

	return toStringArrayFromInterface(set.List())
}
