resource "svix_environment" "example_environment" {
  name = "Staging env"
  type = "development"
}

resource "svix_operational_webhooks_endpoint" "example_endpoint" {
  environment_id = svix_environment.example_environment.id
  url            = "https://example.com"
  description    = "example description"
  rate_limit     = 1
  filter_types   = ["background_task.finished"]
  metadata = jsonencode({
    key1 = "foo"
    key2 = "bar"
  })
}
