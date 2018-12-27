# kong_consumer
A representation of Kong's [consumer object](https://docs.konghq.com/0.14.x/admin-api/#consumer-object)

### Example usage

```hcl
resource "kong_consumer" "crocodile-hunter" {
  username = "steve-irwin"
}

```

### Fields Reference
The following fields are supported:


| field     | explanation     | type      | default     | required                         |
| :-------- | :-------------- | :-------- | :---------- | :------------------------------- |
|`custom_id`|A unique identifier representing a user or service of your API. It can be used to map to existing users in your database. |`string`| N/A| N|
|`username`|A unique username representing a consumer of the API. |`string`| N/A| N|




### Import
Existing Kong consumers can also be imported into Terraform state
Example:
 ```bash
  terraform import kong_consumer.crocodile-hunter <consumer UUID>
```
