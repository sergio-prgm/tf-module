resource "azurerm_resource_group" "res-0" {
  for_each = { for k, v in var.resource_groups : k => v }
  location = each.value.location
  name     = each.value.name
  tags     = each.value.tags
}
# resource "azurerm_resource_group" "res-0" {
#   location = "eastus"
#   name     = "defender-prueba"
#   tags = {
#     responsable = "Sergio Esteve"
#   }
# }
