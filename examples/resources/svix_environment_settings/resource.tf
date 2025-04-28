resource "svix_environment" "example_environment" {
  name = "Staging env"
  type = "development"
}

resource "svix_environment_settings" "example_environment_settings" {
  environment_id = svix_environment.example_environment.id
  color_palette_dark = {
    background_hover     = "#1A202C"
    background_primary   = "#1A202C"
    background_secondary = "#171923"
    button_primary       = "#4299E1"
    interactive_accent   = "#4299E1"
    navigation_accent    = "#4299E1"
    primary              = "#3182CE"
    text_danger          = "#FC8181"
    text_primary         = "#FFFFFF"
  }
  color_palette_light = {
    background_hover     = "#EDF2F7"
    background_primary   = "#F8F9FD"
    background_secondary = "#FFFFFF"
    button_primary       = "#3182CE"
    interactive_accent   = "#3182CE"
    navigation_accent    = "#3182CE"
    primary              = "#3182CE"
    text_danger          = "#E53E3E"
    text_primary         = "#1A202C"
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
  disable_endpoint_on_failure   = false
  display_name                  = "My company"
  enable_channels               = false
  enable_endpoint_mtls_config   = false # Requires Enterprise plan
  enable_endpoint_oauth_config  = false # Requires Enterprise plan
  enable_integration_management = true
  enable_message_stream         = false # Requires Pro or Enterprise plan
  enable_transformations        = false
  enforce_https                 = true
  event_catalog_published       = false
  read_only                     = false
  require_endpoint_channel      = false
  whitelabel_headers            = false # Requires Pro or Enterprise plan
  wipe_successful_payload       = false # Requires Pro or Enterprise plan

  # Advanced settings
  channels_strings_override = {
    channels_one  = "channel"
    channels_many = "channels"
    channels_help = "Channels are an extra dimension of filtering messages orthogonal to event types. They are case-sensitive and only messages with the corresponding channel will be sent to this endpoint."
  }
}
