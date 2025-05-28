resource "svix_environment" "example_environment" {
  name = "Staging env"
  type = "development"
}

resource "svix_environment_settings" "example_environment_settings" {
  environment_id = svix_environment.example_environment.id
  whitelabel_settings = {
    display_name    = "My Company"
    logo_url        = "https://www.example.com/static/logo.png"
    base_font_size  = 16
    font_family     = "Custom"
    font_family_url = "https://fonts.gstatic.com/s/librebaskerville.woff2"
    border_radius = {
      button = "full"
      card   = "lg"
      input  = "none"
    }
    color_palette_dark = {
      surface_hover      = "#1A202C"
      background         = "#1A202C"
      surface_secondary  = "#171923"
      button_primary     = "#4299E1"
      interactive_accent = "#4299E1"
      navigation_accent  = "#4299E1"
      primary            = "#3182CE"
      text_danger        = "#FC8181"
      text_primary       = "#FFFFFF"
    }
    color_palette_light = {
      surface_hover      = "#EDF2F7"
      background         = "#F8F9FD"
      surface_secondary  = "#FFFFFF"
      button_primary     = "#3182CE"
      interactive_accent = "#3182CE"
      navigation_accent  = "#3182CE"
      primary            = "#3182CE"
      text_danger        = "#E53E3E"
      text_primary       = "#1A202C"
    }
    # Advanced settings
    channels_strings_override = {
      channels_one  = "channel"
      channels_many = "channels"
      channels_help = "Channels are an extra dimension of filtering messages orthogonal to event types. They are case-sensitive and only messages with the corresponding channel will be sent to this endpoint."
    }
  }
  disable_endpoint_on_failure           = false
  enable_advanced_endpoint_types        = false # Requires Pro or Enterprise plan
  enable_channels                       = false
  enable_endpoint_mtls_config           = false # Requires Enterprise plan
  enable_endpoint_oauth_config          = false # Requires Enterprise plan
  enable_transformations                = false
  enforce_https                         = true
  event_catalog_published               = false
  require_endpoint_channels             = false
  require_endpoint_event_types          = false
  whitelabel_headers                    = false # Requires Pro or Enterprise plan
  delete_payload_on_successful_delivery = false # Requires Pro or Enterprise plan
}
