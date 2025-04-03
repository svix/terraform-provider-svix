resource "svix_operational_webhooks_endpoint" "example_endpoint" {
  url          = "https://example.com"
  description  = "example description"
  rate_limit   = 1
  filter_types = ["background_task.finished"]
  metadata = jsonencode({
    key1 = "foo"
    key2 = "bar"
  })
}
