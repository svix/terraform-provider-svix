resource "svix_environment" "example_environment" {
  name = "Staging env"
  type = "development"
}

resource "svix_environment_settings" "example_environment_settings" {
  environment_id = svix_environment.example_environment.id
  color_palette_dark = {
    primary = "#3182CE"
  }
  color_palette_light = {
    primary = "#3182CE"
  }
  base_font_size  = 16
  font_family     = "Custom"
  font_family_url = "https://fonts.gstatic.com/s/librebaskerville.woff2"
  logo_url        = "https://www.example.com/static/logo.png"
  theme_override = {
    border_radius = {
      button = "full"
      card   = "lg"
      input  = "none"
    }
  }
  disable_endpoint_on_failure    = false
  display_name                   = "My company"
  enable_channels                = false
  enable_endpoint_mtls_config    = false # Requires Enterprise plan
  enable_endpoint_oauth_config   = false # Requires Enterprise plan
  enable_integration_management  = true
  enable_advanced_endpoint_types = false # Requires Pro or Enterprise plan
  enable_transformations         = false
  enforce_https                  = true
  event_catalog_published        = false
  read_only                      = false
  require_endpoint_channel       = false
  whitelabel_headers             = false # Requires Pro or Enterprise plan
  wipe_successful_payload        = false # Requires Pro or Enterprise plan

  # Advanced settings
  channels_strings_override = {
    channels_one  = "channel"
    channels_many = "channels"
    channels_help = "Channels are an extra dimension of filtering messages orthogonal to event types. They are case-sensitive and only messages with the corresponding channel will be sent to this endpoint."
  }
}
