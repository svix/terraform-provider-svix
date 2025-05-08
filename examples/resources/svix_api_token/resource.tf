resource "svix_environment" "example_environment" {
  name = "Staging env"
  type = "development"
}

resource "svix_api_token" "example_token" {
  environment_id = svix_environment.example_environment.id
  name           = "Environment token"
  scopes         = ["application:Read"]
}

