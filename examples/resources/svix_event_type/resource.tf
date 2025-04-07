resource "svix_event_type" "example_event_type" {
  name        = "invoice.paid"
  description = "An invoice was paid by a user"
  schemas = jsonencode({
    "1" = {
      description = "An invoice was paid by a user"
      properties = {
        invoiceId = {
          description = "The invoice id"
          type        = "string"
        },
        userId = {
          description = "The user id"
          type        = "string"
        }
      }
      required = ["invoiceId", "userId"]
      title    = "Invoice Paid Event"
      type     = "object"
    }
  })
}
