---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "svix_operational_webhooks_endpoint Resource - svix"
subcategory: ""
description: |-
  
---

# svix_operational_webhooks_endpoint (Resource)



## Example Usage

```terraform
resource "svix_operational_webhooks_endpoint" "example_endpoint" {
  url = "https://example.com"
  metadata = jsonencode({
    key1 = "foo"
    key2 = "bar"
  })
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filter_types` (List of String)
- `url` (String)

### Optional

- `description` (String)
- `disabled` (Boolean)
- `metadata` (String)
- `rate_limit` (Number)
- `secret` (String, Sensitive)
- `uid` (String)

### Read-Only

- `created_at` (String)
- `id` (String) The ID of this resource.
- `updated_at` (String)
