terraform {
  required_providers {
    svix = {
      source  = "registry.terraform.io/svix/svix"
      version = "0.3.4"
    }
  }
}

# Configuration-based authentication (required for most resources)
provider "svix" {
  server_url = "https://api.eu.svix.com"
  token      = "***"
}

# Credentials can also be set via SVIX_TOKEN and SVIX_SERVER_URL.
# They are optional when only using svix_autoconfig (auth is in the AutoConfig token).
provider "svix" {}
