resource "azurerm_resource_group" "res-0" {
  location = "eastus"
  name     = "defender-prueba"
  tags = {
    responsable = "Sergio Esteve"
  }
}
resource "azurerm_virtual_network" "res-1" {
  address_space       = ["10.2.0.0/16"]
  location            = "eastus"
  name                = "VM1-vnet"
  resource_group_name = "defender-prueba"
  depends_on = [
    azurerm_resource_group.res-0,
  ]
}
resource "azurerm_subnet" "res-2" {
  address_prefixes     = ["10.2.0.0/24"]
  name                 = "default"
  resource_group_name  = "defender-prueba"
  virtual_network_name = "VM1-vnet"
  depends_on = [
    azurerm_virtual_network.res-1,
  ]
}
resource "azurerm_virtual_network" "res-3" {
  address_space       = ["10.3.0.0/16"]
  location            = "northeurope"
  name                = "vm2-vnet"
  resource_group_name = "defender-prueba"
  depends_on = [
    azurerm_resource_group.res-0,
  ]
}
resource "azurerm_subnet" "res-4" {
  address_prefixes     = ["10.3.0.0/24"]
  name                 = "default"
  resource_group_name  = "defender-prueba"
  virtual_network_name = "vm2-vnet"
  depends_on = [
    azurerm_virtual_network.res-3,
  ]
}
resource "azurerm_storage_account" "res-5" {
  account_replication_type        = "LRS"
  account_tier                    = "Standard"
  allow_nested_items_to_be_public = false
  location                        = "westeurope"
  name                            = "storageaccname1"
  resource_group_name             = "defender-prueba"
  tags = {
    ms-resource-usage = "azure-cloud-shell"
  }
  depends_on = [
    azurerm_resource_group.res-0,
  ]
}
resource "azurerm_storage_account" "res-10" {
  account_kind                    = "Storage"
  account_replication_type        = "LRS"
  account_tier                    = "Standard"
  default_to_oauth_authentication = true
  location                        = "eastus"
  name                            = "storageacc2"
  resource_group_name             = "defender-prueba"
  depends_on = [
    azurerm_resource_group.res-0,
  ]
}
resource "azurerm_storage_container" "res-12" {
  name                 = "azure-webjobs-hosts"
  storage_account_name = "defenderpruebaa10c"
}
resource "azurerm_storage_container" "res-13" {
  name                 = "azure-webjobs-secrets"
  storage_account_name = "defenderpruebaa10c"
}
resource "azurerm_storage_container" "res-14" {
  name                 = "scm-releases"
  storage_account_name = "defenderpruebaa10c"
}
