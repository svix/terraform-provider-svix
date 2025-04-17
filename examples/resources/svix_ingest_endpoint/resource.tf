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

resource "svix_ingest_endpoint" "example_endpoint" {
  environment_id   = svix_environment.example_environment.id
  ingest_source_id = svix_ingest_source.example_ingest_source.id
  url              = "https://example.com"
  description      = "example description"
  rate_limit       = 1
  metadata = jsonencode({
    key1 = "foo"
    key2 = "bar"
  })
}
