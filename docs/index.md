---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Svix Provider"
subcategory: ""
description: |-
  
---

# Svix Provider



## Example Usage

```terraform
terraform {
  required_providers {
    svix = {
      source  = "registry.terraform.io/svix/svix"
      version = "0.3.0"
    }
  }
}

# Configuration-based authentication
provider "svix" {
  server_url = "https://api.eu.svix.com"
  token      = "***"
}

# The provider can also be configured via the SVIX_TOKEN and SVIX_SERVER_URL environment variables
provider "svix" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `server_url` (String) Svix server url
- `token` (String, Sensitive) Api token
