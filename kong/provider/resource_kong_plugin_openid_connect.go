package provider

import (
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKongPluginOpenIdConnect() *schema.Resource {
	clientResource :=  &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type: schema.TypeString,
				Required: true,
			},
			"secret": {
				Type: schema.TypeString,
				Optional: true,
			},
			"redirect_uri": {
				Type: schema.TypeString,
				Optional: true,
			},
			"login_redirect_uri": {
				Type: schema.TypeString,
				Optional: true,
			},
			"logout_redirect_uri": {
				Type: schema.TypeString,
				Optional: true,
			},
			"forbidden_redirect_uri": {
				Type: schema.TypeString,
				Optional: true,
			},
			"unauthorized_redirect_uri": {
				Type: schema.TypeString,
				Optional: true,
			},
			"unexpected_redirect_uri": {
				Type: schema.TypeString,
				Optional: true,
			},
		},
	}

	return CreateGenericPluginResource(&GenericPluginResource{
		Name: "",
		AdditionalSchema: map[string]*schema.Schema{
			"issuer": {
				Type: schema.TypeString,
				Required: true,
				// ValidateFunc: IsUrl
			},
			"client_arg": {
			 	Type: schema.TypeString,
				Optional: true,
				Default: "client_id",
			},
			"clients": {
				Type: schema.TypeList, // the order of clients is important to the Kong plugin
				Elem: clientResource,
			},
			"forbidden_destroy_session": {
				Type: schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"scopes": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					return []string{"open_id"}, nil
				},
				Set: schema.HashString,
			},
			"scopes_required": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Set: schema.HashString,
			},
			"response_mode": {
				Type: schema.TypeString,
				Optional: true,
				Default: "query",
			},
			"auth_methods": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					return []string{
						"password",
						"client_credentials",
						"authorization_code",
						"bearer",
						"introspection",
						"kong_oauth2",
						"refresh_token",
						"session",
					}, nil
				},
			},
			"audience": {
				// not sure if order is important to Kong for audience/audience_required
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ConflictsWith: []string{"audience_required"},
			},
			"audience_required": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ConflictsWith: []string{"audience"},
			},
			"domains": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Set: schema.HashString,
			},
			"max_age_seconds": {
				Type: schema.TypeInt,
				Optional: true,
				Default: 0,
			},
			"authorization_cookie": {
				Type: schema.TypeMap,
				Elem: makeCookieResource("authorization", 600),
				Optional: true,
			},
			"session": {
				Type: schema.TypeMap,
				Elem: &schema.Resource{

				},
			},
			"session_cookie_name": {
				Type: schema.TypeString,
			},
			"session_cookie_lifetime": {
				Type: schema.TypeString,
			},
			"session_storage": {
				Type: schema.TypeString,
			},
			"session_memcache_prefix": {
				Type: schema.TypeString,
			},
			"session_memcache_socket": {
				Type: schema.TypeString,
			},
			"session_memcache_host": {
				Type: schema.TypeString,
			},
			"session_memcache_port": {
				Type: schema.TypeString,
			},
			"session_redis_prefix": {
				Type: schema.TypeString,
			},
			"session_redis_socket": {
				Type: schema.TypeString,
			},
			"session_redis_host": {
				Type: schema.TypeString,
			},
			"session_redis_port": {
				Type: schema.TypeString,
			},
			"session_redis_auth": {
				Type: schema.TypeString,
			},
			"extra_jwks_uris": {
				Type: schema.TypeString,
			},
			"jwt_session_cookie": {
				Type: schema.TypeString,
			},
			"jwt_session_claim": {
				Type: schema.TypeString,
			},
			"reverify": {
				Type: schema.TypeString,
			},
			"bearer_token_param_type": {
				Type: schema.TypeString,
			},
			"client_credentials_param_type": {
				Type: schema.TypeString,
			},
			"password_param_type": {
				Type: schema.TypeString,
			},
			"id_token_param_name": {
				Type: schema.TypeString,
			},
			"id_token_param_type": {
				Type: schema.TypeString,
			},
			"discovery_headers_names": {
				Type: schema.TypeString,
			},
			"discovery_headers_values": {
				Type: schema.TypeString,
			},
			"authorization_query_args_names": {
				Type: schema.TypeString,
			},
			"authorization_query_args_values": {
				Type: schema.TypeString,
			},
			"authorization_query_args_client": {
				Type: schema.TypeString,
			},
			"token_post_args_names": {
				Type: schema.TypeString,
			},
			"token_post_args_values": {
				Type: schema.TypeString,
			},
			"token_headers_client": {
				Type: schema.TypeString,
			},
			"token_headers_replay": {
				Type: schema.TypeString,
			},
			"token_headers_prefix": {
				Type: schema.TypeString,
			},
			"token_headers_grants": {
				Type: schema.TypeString,
			},
			"upstream_headers_claims": {
				Type: schema.TypeString,
			},
			"upstream_headers_names": {
				Type: schema.TypeString,
			},
			"downstream_headers_claims": {
				Type: schema.TypeString,
			},
			"downstream_headers_names": {
				Type: schema.TypeString,
			},
			"upstream_access_token_header": {
				Type: schema.TypeString,
			},
			"downstream_access_token_header": {
				Type: schema.TypeString,
			},
			"upstream_access_token_jwk_header": {
				Type: schema.TypeString,
			},
			"downstream_access_token_jwk_header": {
				Type: schema.TypeString,
			},
			"upstream_id_token_header": {
				Type: schema.TypeString,
			},
			"downstream_id_token_header": {
				Type: schema.TypeString,
			},
			"upstream_id_token_jwk_header": {
				Type: schema.TypeString,
			},
			"downstream_id_token_jwk_header": {
				Type: schema.TypeString,
			},
			"upstream_refresh_token_header": {
				Type: schema.TypeString,
			},
			"downstream_refresh_token_header": {
				Type: schema.TypeString,
			},
			"upstream_user_info_header": {
				Type: schema.TypeString,
			},
			"downstream_user_info_header": {
				Type: schema.TypeString,
			},
			"upstream_introspection_header": {
				Type: schema.TypeString,
			},
			"downstream_introspection_header": {
				Type: schema.TypeString,
			},
			"introspect_jwt_tokens": {
				Type: schema.TypeString,
			},
			"introspection_endpoint": {
				Type: schema.TypeString,
			},
			"introspection_hint": {
				Type: schema.TypeString,
			},
			"introspection_headers_names": {
				Type: schema.TypeString,
			},
			"introspection_headers_values": {
				Type: schema.TypeString,
			},
			"login_methods": {
				Type: schema.TypeString,
			},
			"login_action": {
				Type: schema.TypeString,
			},
			"login_tokens": {
				Type: schema.TypeString,
			},
			"login_redirect_mode": {
				Type: schema.TypeString,
			},
			"logout_query_arg": {
				Type: schema.TypeString,
			},
			"logout_post_arg": {
				Type: schema.TypeString,
			},
			"logout_uri_suffix": {
				Type: schema.TypeString,
			},
			"logout_methods": {
				Type: schema.TypeString,
			},
			"logout_revoke": {
				Type: schema.TypeString,
			},
			"revocation_endpoint": {
				Type: schema.TypeString,
			},
			"end_session_endpoint": {
				Type: schema.TypeString,
			},
			"token_exchange_endpoint": {
				Type: schema.TypeString,
			},
			"consumer_claim": {
				Type: schema.TypeString,
			},
			"consumer_by": {
				Type: schema.TypeString,
			},
			"credential_claim": {
				Type: schema.TypeString,
			},
			"anonymous": {
				Type: schema.TypeString,
			},
			"run_on_preflight": {
				Type: schema.TypeString,
			},
			"leeway": {
				Type: schema.TypeString,
			},
			"verify_parameters": {
				Type: schema.TypeString,
			},
			"verify_nonce": {
				Type: schema.TypeString,
			},
			"verify_signature": {
				Type: schema.TypeString,
			},
			"verify_claims": {
				Type: schema.TypeString,
			},
			"cache_ttl": {
				Type: schema.TypeString,
			},
			"cache_introspection": {
				Type: schema.TypeString,
			},
			"cache_token_exchange": {
				Type: schema.TypeString,
			},
			"cache_tokens": {
				Type: schema.TypeString,
			},
			"cache_user_info": {
				Type: schema.TypeString,
			},
			"hide_credentials": {
				Type: schema.TypeString,
			},
			"http_version": {
				Type: schema.TypeString,
			},
			"ssl_verify": {
				Type: schema.TypeString,
			},
			"timeout": {
				Type: schema.TypeString,
			},
		},
		MapApiModelToResource: func(plugin *kong.KongPlugin, data *schema.ResourceData) {

		},
		MapSchemaToPluginConfig: func(data *schema.ResourceData) interface{} {
			return nil
		},
	})
}

func makeCookieResource(name string, lifetimeInSeconds int) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Optional: true,
				Default: name,
			},
			"lifetime_seconds": {
				Type: schema.TypeInt,
				Optional: true,
				Default: lifetimeInSeconds,
			},
		},
	}
}
