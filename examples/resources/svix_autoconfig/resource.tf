# Provider credentials are not required for AutoConfig.
provider "svix" {}

resource "svix_autoconfig" "example" {
  token = "auto_v1_..."
  url   = "https://api.us.example.com/webhooks/acme"
  filter_types = [
    "invoice.paid",
    "user.created",
  ]
}
