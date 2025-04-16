resource "svix_environment" "example_environment" {
  name = "Staging env"
  type = "development"
}


resource "svix_event_type_openapi_import" "event_types" {
  environment_id = svix_environment.example_environment.id
  spec_raw       = file("${path.module}/openapi.json")
}
