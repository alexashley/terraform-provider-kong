package provider

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKongPluginOpenidConnect() *schema.Resource {
	return CreateGenericPluginResource(&GenericPluginResource{
		AllowsConsumers: false,
		Name: func(data *schema.ResourceData) string {
			return "openid-connect"
		},
		AdditionalSchema: map[string]*schema.Schema{
			"anonymous": {
				Description:  "Anonymous consumer id. This is useful if you need to enable multiple auth plugins -- failing to authenticate will cause this consumer to be set.",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateUuid,
			},
			"auth_methods": {
				Description: "Allowed authentication methods",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"consumer_by": {
				Description: "A JWT claim used to lookup a Kong consumer. Used with consumer_claim to control the process of identifying a Kong consumer.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"consumer_claim": {
				Description: "JWT claims to use to map to a Kong consumer. Typically set to `sub`",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"issuer": {
				Description:  "URL of the OpenId Connect server",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateUrl,
			},
		},
		MapSchemaToPluginConfig: func(data *schema.ResourceData) (interface{}, error) {
			config := make(map[string]interface{})
			config["issuer"] = data.Get("issuer").(string)
			if value, ok := data.GetOk("anonymous"); ok {
				config["anonymous"] = value
			}

			if consumerClaim, ok := data.GetOk("consumer_claim"); ok {
				claims := setToStringArray(consumerClaim.(*schema.Set))
				config["consumer_claim"] = claims
			}

			if authMethods, ok := data.GetOk("auth_methods"); ok {
				methods := setToStringArray(authMethods.(*schema.Set))
				if err := validateAuthMethods(methods); err != nil {
					return nil, err
				}

				config["auth_methods"] = methods
			}

			if consumerBys, ok := data.GetOk("consumer_by"); ok {
				consumerBy := setToStringArray(consumerBys.(*schema.Set))
				if err := validateConsumerBy(consumerBy); err != nil {
					return nil, err
				}

				config["consumer_by"] = consumerBy
			}

			return config, nil
		},
		MapApiModelToResource: func(plugin *kong.KongPlugin, data *schema.ResourceData) error {
			config := plugin.Config.(map[string]interface{})

			data.Set("issuer", config["issuer"])

			optionals := []string{"anonymous", "auth_methods", "consumer_claim"}

			for _, attribute := range optionals {
				if attributeValue, ok := config[attribute]; ok {
					data.Set(attribute, attributeValue)
				}
			}

			return nil
		},
	})
}

func validateAuthMethods(authMethods []string) error {
	validAuthMethods := []string{
		"password",
		"client_credentials",
		"authorization_code",
		"bearer",
		"introspection",
		"kong_oauth2",
		"refresh_token",
		"session",
	}

	for _, method := range authMethods {
		match := false
		for _, validMethod := range validAuthMethods {
			if method == validMethod {
				match = true
				break
			}
		}
		if !match {
			return fmt.Errorf("%s is not a valid auth_method", method)
		}
	}

	return nil
}

func validateConsumerBy(consumerByFields []string) error {
	for _, consumerBy := range consumerByFields {
		if !(consumerBy == "username" || consumerBy == "consumer") {
			return fmt.Errorf("invalid value for consumer_by: must be one of custom_id or username")
		}
	}

	return nil
}
