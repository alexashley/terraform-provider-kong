# [kong_service](https://github.com/alexashley/terraform-provider-kong/tree/master/kong/provider/resource_kong_service.go)
A representation of Kong's [service object](https://docs.konghq.com/0.14.x/admin-api/#service-object)

### Example usage

~~~
resource "kong_service" "mockbin" {
  name  = "mockbin"
  url   = "https://mockbin.org/request"
}
~~~

### Fields Reference
The following fields are supported:


| field     | explanation     | type      | default     | required                         |
| :-------- | :-------------- | :-------- | :---------- | :------------------------------- |
|`name`|The service name. |`string`| | Y|
|`url`|The url for the service. It encapsulates protocol, host, port and path |`string`| | Y|
|`connect_timeout`|Time in milliseconds to connect to the upstream server. |`int`| 60000| N|
|`read_timeout`|Time in milliseconds between two read operations to the upstream server. |`int`| 60000| N|
|`retries`|Number of times Kong will try to proxy if there's an error. |`int`| 5| N|
|`write_timeout`|Time in milliseconds between two successive write operations to the upstream server. |`int`| 60000| N|


### Computed Fields
The following computed attributes are also available:

| field     | explanation     | type    |
|-----------|-----------------|---------|
|`created_at`|Unix timestamp representing the time the service was created. |int|
|`updated_at`|Unix timestamp for the last time the service was updated. |int|

### Import
Existing Kong services can be imported:
```bash
terraform import kong_service.name-of-service-to-import <service UUID>
```

[GitHub](https://github.com/alexashley/terraform-provider-kong)
