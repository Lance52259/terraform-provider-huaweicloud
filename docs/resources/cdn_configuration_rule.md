---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_configuration_rule"
description: |-
  Manages a configuration rule resource within HuaweiCloud.
---

# huaweicloud_cdn_configuration_rule

Manages a configuration rule resource within HuaweiCloud.

## Example Usage

### Create a CDN configuration rule with basic actions

```hcl
variable "domain_name" {}

resource "huaweicloud_cdn_configuration_rule" "test" {
  domain_name = var.domain_name
  name        = "test_rule"
  status      = "on"
  priority    = 1

  conditions = jsonencode({
    match_type  = "file_extension"
    match_value = ".jpg,.png"
  })

  actions {
    cache_rule {
      ttl          = 3600
      ttl_unit     = "s"
      follow_origin = "off"
      force_cache   = "on"
    }
  }
}
```

### Create a CDN configuration rule with flexible origin

```hcl
variable "domain_name" {}
variable "ip_or_domain" {}

resource "huaweicloud_cdn_configuration_rule" "test" {
  domain_name = var.domain_name
  name        = "test_rule"
  status      = "on"
  priority    = 1

  conditions = jsonencode({
    match_type  = "catalog"
    match_value = "/test/"
  })

  actions {
    flexible_origin {
      sources_type = "ipaddr"
      ip_or_domain = var.ip_or_domain
      priority     = 1
      weight       = 50
      http_port    = 80
      https_port   = 443
      host_name    = "example.com"
      origin_protocol = "follow"
    }
  }
}
```

### Create a CDN configuration rule with multiple actions

```hcl
variable "domain_name" {}

resource "huaweicloud_cdn_configuration_rule" "test" {
  domain_name = var.domain_name
  name        = "test_rule"
  status      = "on"
  priority    = 1

  conditions = jsonencode({
    match_type  = "file_path"
    match_value = "/test/*"
  })

  actions {
    origin_request_header {
      name   = "X-Forwarded-For"
      action = "set"
      value  = "192.168.1.1"
    }

    http_response_header {
      name   = "X-Custom-Header"
      action = "set"
      value  = "custom-value"
    }

    access_control {
      type = "block"
    }

    request_limit_rule {
      limit_rate_after = 50
      limit_rate_value = 1048576
    }

    origin_request_url_rewrite {
      rewrite_type = "simple"
      source_url   = "/old/path"
      target_url   = "/new/path"
    }

    request_url_rewrite {
      redirect_url         = "/redirect"
      execution_mode       = "redirect"
      redirect_status_code = 302
      redirect_host        = "example.com"
    }

    browser_cache_rule {
      cache_type = "ttl"
      ttl        = 1800
      ttl_unit   = "s"
    }

    error_code_cache {
      code = 404
      ttl  = 300
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `domain_name` - (Required, String, NonUpdatable) Specifies the accelerated domain name.  
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the configuration rule.  
  The valid length is limit from `1` to `50`.

* `status` - (Required, String) Whether to enable the configuration rule.  
  The valid values are as follows:
  + **on**
  + **off**

* `priority` - (Required, Int) Specifies the priority of the configuration rule.
  The valid value is range from `1` to `100`.

* `conditions` - (Optional, String) Specifies the trigger conditions of the configuration rule, in JSON format.

* `actions` - (Required, List) Specifies the list of actions to be performed when the configuration rule is met.
  The [actions](#cdn_configuration_rule_actions) structure is documented below.

<a name="cdn_configuration_rule_actions"></a>
The `actions` block supports:

* `flexible_origin` - (Optional, List) Specifies the list of flexible origin configurations.  
  The [flexible_origin](#cdn_configuration_rule_actions_flexible_origin) structure is documented below.

* `origin_request_header` - (Optional, List) Specifies the list of origin request header configurations.  
  The [origin_request_header](#cdn_configuration_rule_actions_origin_request_header) structure is documented below.

* `http_response_header` - (Optional, List) Specifies the list of HTTP response header configurations.  
  The [http_response_header](#cdn_configuration_rule_actions_http_response_header) structure is documented below.

* `access_control` - (Optional, List) Specifies the access control configuration.  
  The [access_control](#cdn_configuration_rule_actions_access_control) structure is documented below.

* `request_limit_rule` - (Optional, List) Specifies the request rate limit configuration.  
  The [request_limit_rule](#cdn_configuration_rule_actions_request_limit_rule) structure is documented below.

* `origin_request_url_rewrite` - (Optional, List) Specifies the origin request URL rewrite configuration.  
  The [origin_request_url_rewrite](#cdn_configuration_rule_actions_origin_request_url_rewrite) structure is documented below.

* `cache_rule` - (Optional, List) Specifies the cache rule configuration.  
  The [cache_rule](#cdn_configuration_rule_actions_cache_rule) structure is documented below.

* `request_url_rewrite` - (Optional, List) Specifies the access URL rewrite configuration.  
  The [request_url_rewrite](#cdn_configuration_rule_actions_request_url_rewrite) structure is documented below.

* `browser_cache_rule` - (Optional, List) Specifies the browser cache rule configuration.  
  The [browser_cache_rule](#cdn_configuration_rule_actions_browser_cache_rule) structure is documented below.

* `error_code_cache` - (Optional, List) Specifies the list of error code cache configurations.  
  The [error_code_cache](#cdn_configuration_rule_actions_error_code_cache) structure is documented below.

* `origin_range` - (Optional, List) Specifies the origin range configuration.  
  The [error_code_cache](#cdn_configuration_rule_actions_origin_range) structure is documented below.

<a name="cdn_configuration_rule_actions_flexible_origin"></a>
The `flexible_origin` block supports:

* `sources_type` - (Required, String) Specifies the source type.  
  The valid values are as follows:
  + **ipaddr**
  + **domain**
  + **obs_bucket**
  + **third_bucket**

* `ip_or_domain` - (Required, String) Specifies the origin IP or domain name.

* `priority` - (Required, Int) Specifies the origin priority.  
  The valid value is range from `1` to `100`.

* `weight` - (Required, Int) Specifies the origin weight.  
  The valid value is range from `1` to `100`.

* `obs_bucket_type` - (Optional, String) Specifies the OBS bucket type.  
  The valid values are as follows:
  + **private**
  + **public**

* `bucket_access_key` - (Optional, String) Specifies the third-party object storage access key.

* `bucket_secret_key` - (Optional, String) Specifies the third-party object storage secret key.

* `bucket_region` - (Optional, String) Specifies the third-party object storage region.

* `bucket_name` - (Optional, String) Specifies the third-party object storage name.

* `host_name` - (Optional, String) Specifies the origin host name.

* `origin_protocol` - (Optional, String) Specifies the origin protocol.  
  The valid values are as follows:
  + **follow**
  + **http**
  + **https**

* `http_port` - (Optional, Int) Specifies the HTTP port number.  
  The valid value is range from `1` to `65,535`.

* `https_port` - (Optional, Int) Specifies the HTTPS port number.  
  The valid value is range from `1` to `65,535`.

<a name="cdn_configuration_rule_actions_origin_request_header"></a>
The `origin_request_header` block supports:

* `name` - (Required, String) Specifies the back-to-origin request header parameter name.

* `action` - (Required, String) Specifies the back-to-origin request header setting type.  
  The valid values are as follows:
  + **delete**
  + **set**

* `value` - (Optional, String) Specifies the back-to-origin request header parameter value.

<a name="cdn_configuration_rule_actions_http_response_header"></a>
The `http_response_header` block supports:

* `name` - (Required, String) Specifies the HTTP response header parameter name.

* `action` - (Required, String) Specifies the operation type of setting HTTP response header.  
  The valid values are as follows:
  + **set**
  + **delete**

* `value` - (Optional, String) Specifies the HTTP response header parameter value.

<a name="cdn_configuration_rule_actions_access_control"></a>
The `access_control` block supports:

* `type` - (Required, String) Specifies the access control type.  
  The valid values are as follows:
  + **block**
  + **trust**

<a name="cdn_configuration_rule_actions_request_limit_rule"></a>
The `request_limit_rule` block supports:

* `limit_rate_after` - (Required, Int) Specifies the rate limit condition.

* `limit_rate_value` - (Required, Int) Specifies the rate limit value.

<a name="cdn_configuration_rule_actions_origin_request_url_rewrite"></a>
The `origin_request_url_rewrite` block supports:

* `rewrite_type` - (Required, String) Specifies the rewrite type.  
  The valid values are as follows:
  + **simple**
  + **wildcard**
  + **regex**

* `target_url` - (Required, String) Specifies the target URL.

* `source_url` - (Optional, String) Specifies the source URL to be rewritten.

<a name="cdn_configuration_rule_actions_cache_rule"></a>
The `cache_rule` block supports:

* `ttl` - (Required, Int) Specifies the cache expiration time.

* `ttl_unit` - (Required, String) Specifies the cache expiration time unit.  
  The valid values are as follows:
  + **s**
  + **m**
  + **h**
  + **d**

* `follow_origin` - (Required, String) Specifies the cache expiration time source.  
  The valid values are as follows:
  + **off**
  + **on**
  + **min_ttl**

* `force_cache` - (Optional, String) Whether to enable forced caching.  
  The valid values are as follows:
  + **on**
  + **off**

<a name="cdn_configuration_rule_actions_request_url_rewrite"></a>
The `request_url_rewrite` block supports:

* `redirect_url` - (Required, String) Specifies the redirect URL.

* `execution_mode` - (Required, String) Specifies the execution mode.  
  The valid values are as follows:
  + **redirect**
  + **break**

* `redirect_status_code` - (Optional, Int) Specifies the redirect status code.  
  The valid values are as follows:
  + `301`
  + `302`
  + `303`
  + `307`

* `redirect_host` - (Optional, String) Specifies the redirect host.

<a name="cdn_configuration_rule_actions_browser_cache_rule"></a>
The `browser_cache_rule` block supports:

* `cache_type` - (Required, String) Specifies the cache effective type.  
  The valid values are as follows:
  + **follow_origin**
  + **ttl**
  + **never**

* `ttl` - (Optional, Int) Specifies the cache expiration time.

* `ttl_unit` - (Optional, String) Specifies the cache expiration time unit.  
  The valid values are as follows:
  + **s**
  + **m**
  + **h**
  + **d**

<a name="cdn_configuration_rule_actions_error_code_cache"></a>
The `error_code_cache` block supports:

* `code` - (Required, Int) Specifies the error code to be cached.

* `ttl` - (Required, Int) Specifies the error code cache time.

<a name="cdn_configuration_rule_actions_origin_range"></a>
The `origin_range` block supports:

* `status` - (Required, String) Specifies the origin range status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The CDN configuration rule can be imported using the format `<domain_name>/<id>` or `<domain_name>/<name>`, e.g.

```bash
$ terraform import huaweicloud_cdn_configuration_rule.test <domain_name>/<id>
```

or

```bash
$ terraform import huaweicloud_cdn_configuration_rule.test <domain_name>/<name>
```
