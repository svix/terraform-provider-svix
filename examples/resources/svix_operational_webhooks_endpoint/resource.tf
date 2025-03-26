resource "svix_operational_webhooks_endpoint" "example_endpoint" {
  url = "https://example.com"
  metadata = jsonencode({
    key1 = "foo"
    key2 = "bar"
  })
}
