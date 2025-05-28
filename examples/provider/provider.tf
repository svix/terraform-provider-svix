terraform {
  required_providers {
    svix = {
      source  = "registry.terraform.io/svix/svix"
      version = "0.2.3"
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
