---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_configuration_rules"
description: |-
  Use this data source to get a list of configuration rules within HuaweiCloud.
---

# huaweicloud_cdn_configuration_rules

Use this data source to get a list of configuration rules within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}

data "huaweicloud_cdn_configuration_rules" "test" {
  domain_name = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the configuration rules are located.  
  If omitted, the provider-level region will be used.

* `domain_name` - (Required, String) Specifies the accelerated domain name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the configuration rules.  
  The [rules](#cdn_configuration_rules_rules) structure is documented below.

<a name="cdn_configuration_rules_rules"></a>
The `rules` block supports:

* `id` - The ID of the configuration rule.

* `name` - The name of the configuration rule.

* `status` - Whether the rule is enabled.  
  The valid values are as follows:
  + **on**
  + **off**

* `priority` - The priority of the configuration rule.  
  The valid value is range from `1` to `100`.

* `conditions` - The trigger conditions of the current rule, in JSON format.

* `actions` - The list of actions to be performed when the rules are met.  
  The [actions](#cdn_configuration_rules_rules_actions) structure is documented below.

<a name="cdn_configuration_rules_rules_actions"></a>
The `actions` block supports:

* `flexible_origin` - The list of flexible origin configurations.  
  The [flexible_origin](#cdn_configuration_rules_rules_actions_flexible_origin) structure is documented below.

* `origin_request_header` - The list of origin request header configurations.  
  The [origin_request_header](#cdn_configuration_rules_rules_actions_origin_request_header) structure is documented below.

* `http_response_header` - The list of HTTP response header configurations.  
  The [http_response_header](#cdn_configuration_rules_rules_actions_http_response_header) structure is documented below.

* `access_control` - The access control configuration.  
  The [access_control](#cdn_configuration_rules_rules_actions_access_control) structure is documented below.

* `request_limit_rule` - The request rate limit configuration.  
  The [request_limit_rule](#cdn_configuration_rules_rules_actions_request_limit_rule) structure is documented below.

* `origin_request_url_rewrite` - The origin request URL rewrite configuration.  
  The [origin_request_url_rewrite](#cdn_configuration_rules_rules_actions_origin_request_url_rewrite) structure is
  documented below.

* `cache_rule` - The cache rule configuration.  
  The [cache_rule](#cdn_configuration_rules_rules_actions_cache_rule) structure is documented below.

* `request_url_rewrite` - The access URL rewrite configuration.  
  The [request_url_rewrite](#cdn_configuration_rules_rules_actions_request_url_rewrite) structure is documented below.

* `browser_cache_rule` - The browser cache rule configuration.  
  The [browser_cache_rule](#cdn_configuration_rules_rules_actions_browser_cache_rule) structure is documented below.

* `error_code_cache` - The list of error code cache configurations.  
  The [error_code_cache](#cdn_configuration_rules_rules_actions_error_code_cache) structure is documented below.

* `origin_range` - The list of origin range configurations.  
  The [origin_range](#cdn_configuration_rules_rules_actions_origin_range) structure is documented below.

<a name="cdn_configuration_rules_rules_actions_flexible_origin"></a>
The `flexible_origin` block supports:

* `sources_type` - The source type.  
  The valid values are as follows:
  + **ipaddr**
  + **domain**
  + **obs_bucket**
  + **third_bucket**

* `ip_or_domain` - The origin IP or domain name.

* `priority` - The origin priority.  
  The valid value is range from `1` to `100`.

* `weight` - The weight.  
  The valid value is range from `1` to `100`.

* `obs_bucket_type` - The OBS bucket type.  
  The valid values are as follows:
  + **private**
  + **public**

* `bucket_access_key` - The third-party object storage access key.

* `bucket_secret_key` - The third-party object storage secret key.

* `bucket_region` - The third-party object storage region.

* `bucket_name` - The third-party object storage name.

* `host_name` - The origin HOST.

* `origin_protocol` - The origin protocol.  
  The valid values are as follows:
  + **follow**
  + **http**
  + **https**

* `http_port` - The HTTP port number.  
  The valid value is range from `1` to `65,535`.

* `https_port` - The HTTPS port number.  
  The valid value is range from `1` to `65,535`.

<a name="cdn_configuration_rules_rules_actions_origin_request_header"></a>
The `origin_request_header` block supports:

* `name` - The back-to-origin request header parameter name.

* `action` - The back-to-origin request header setting type.  
  The valid values are as follows:
  + **delete**
  + **set**

* `value` - The back-to-origin request header parameter value.

<a name="cdn_configuration_rules_rules_actions_http_response_header"></a>
The `http_response_header` block supports:

* `name` - The HTTP response header parameter name.

* `action` - The operation type of setting HTTP response header.  
  The valid values are as follows:
  + **set**
  + **delete**

* `value` - The HTTP response header parameter value.

<a name="cdn_configuration_rules_rules_actions_access_control"></a>
The `access_control` block supports:

* `type` - The access control type.  
  The valid values are as follows:
  + **block**
  + **trust**

<a name="cdn_configuration_rules_rules_actions_request_limit_rules"></a>
The `request_limit_rules` block supports:

* `limit_rate_after` - The rate limit condition.

* `limit_rate_value` - The rate limit value.

<a name="cdn_configuration_rules_rules_actions_origin_request_url_rewrite"></a>
The `origin_request_url_rewrite` block supports:

* `rewrite_type` - The rewrite type.  
  The valid values are as follows:
  + **simple**
  + **wildcard**
  + **regex**

* `target_url` - The target URL.

* `source_url` - The source URL to be rewritten.

<a name="cdn_configuration_rules_rules_actions_cache_rule"></a>
The `cache_rule` block supports:

* `ttl` - The cache expiration time.

* `ttl_unit` - The cache expiration time unit.  
  The valid values are as follows:
  + **s**
  + **m**
  + **h**
  + **d**

* `follow_origin` - The cache expiration time source.  
  The valid values are as follows:
  + **off**
  + **on**
  + **min_ttl**

* `force_cache` - Whether to enable forced caching.  
  The valid values are as follows:
  + **on**
  + **off**

<a name="cdn_configuration_rules_rules_actions_request_url_rewrite"></a>
The `request_url_rewrite` block supports:

* `redirect_url` - The redirect URL.

* `execution_mode` - The execution mode.  
  The valid values are as follows:
  + **redirect**
  + **break**

* `redirect_status_code` - The redirect status code.  
  The valid values are as follows:
  + `301`
  + `302`
  + `303`
  + `307`

* `redirect_host` - The redirect host.

<a name="cdn_configuration_rules_rules_actions_browser_cache_rule"></a>
The `browser_cache_rule` block supports:

* `cache_type` - The cache effective type.  
  The valid values are as follows:
  + **follow_origin**
  + **ttl**
  + **never**

* `ttl` - The cache expiration time.

* `ttl_unit` - The cache expiration time unit.  
  The valid values are as follows:
  + **s**
  + **m**
  + **h**
  + **d**

<a name="cdn_configuration_rules_rules_actions_error_code_cache"></a>
The `error_code_cache` block supports:

* `code` - The error code to be cached.

* `ttl` - The error code cache time.

<a name="cdn_configuration_rules_rules_actions_origin_range"></a>
The `origin_range` block supports:

* `status` - The origin range status.
