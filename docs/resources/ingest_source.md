---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "svix_ingest_source Resource - Svix"
subcategory: ""
description: |-
  
---

# svix_ingest_source (Resource)



## Example Usage

```terraform
resource "svix_environment" "example_environment" {
  name = "Staging env"
  type = "development"
}

resource "svix_ingest_source" "example_ingest_source" {
  environment_id = svix_environment.example_environment.id
  type           = "cron"
  name           = "example cron source"
  uid            = "example-cron-source"
  config = jsonencode({
    schedule = "5 4 * * *" # At 04:05
    payload  = "Some payload"
  })
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The Id to the environment that this resource will be created in
- `name` (String)
- `type` (String) Can be one of `generic-webhook`, `cron`, `adobe-sign`, `beehiiv`, `brex`, `clerk`, `docusign`, `github`, `guesty`, `hubspot`, `incident-io`, `lithic`, `nash`, `pleo`, `replicate`, `resend`, `safebase`, `sardine`, `segment`, `shopify`, `slack`, `stripe`, `stych`, `svix`, `zoom`

### Optional

- `config` (String, Sensitive) The config may include sensitive fields(webhook signing secret for example)

Documentation for the config can be found in the [API docs](https://api.svix.com/docs#tag/Ingest-Source/operation/v1.ingest.source.create)
- `ingest_url` (String)
- `uid` (String)

### Read-Only

- `created_at` (String)
- `id` (String) The ID of this resource.
- `updated_at` (String)
