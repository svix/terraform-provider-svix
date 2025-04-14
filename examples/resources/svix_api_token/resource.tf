resource "svix_environment" "example_environment" {
  name = "Staging env"
  type = "development"
}

resource "svix_api_token" "example_token" {
  environment_id = svix_environment.example_environment.id
  name           = "App token 1"
  scopes         = ["application:Read"]
}

