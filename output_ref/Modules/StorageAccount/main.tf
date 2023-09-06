resource "azurerm_storage_account" "res-5" {
  for_each                        = { for k, v in var.storage_accounts : k => v }
  account_replication_type        = each.value.account_replication_type
  account_tier                    = each.value.account_tier
  allow_nested_items_to_be_public = each.value.allow_nested_items_to_be_public
  location                        = each.value.location
  name                            = each.value.name
  resource_group_name             = each.value.resource_group_name
  tags                            = each.value.tags
}

# resource "azurerm_storage_account" "res-10" {
#   account_kind                    = "Storage"
#   account_replication_type        = "LRS"
#   account_tier                    = "Standard"
#   default_to_oauth_authentication = true
#   location                        = "eastus"
#   name                            = "defenderpruebaa10c"
#   resource_group_name             = "defender-prueba"
# }

resource "azurerm_storage_container" "res-12" {
  for_each             = { for k, v in var.storage_containers : k => v }
  name                 = each.value.name
  storage_account_name = each.value.storage_account_name
}

# resource "azurerm_storage_container" "res-13" {
#   name                 = "azure-webjobs-secrets"
#   storage_account_name = "defenderpruebaa10c"
# }

# resource "azurerm_storage_container" "res-14" {
#   name                 = "scm-releases"
#   storage_account_name = "defenderpruebaa10c"
# }
