resource "svix_environment" "example_environment" {
  name = "Staging env"
  type = "development"
}

resource "svix_ingest_source_resource" "example_ingest_source" {
  environment_id = svix_environment.example_environment.id
  type           = "cron"
  name           = "example cron source"
  uid            = "example-cron-source"
  config = jsonencode({
    schedule = "5 4 * * *" # At 04:05
    payload  = "Some payload"
  })
}

