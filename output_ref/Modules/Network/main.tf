resource "azurerm_virtual_network" "res_virtual_network" {
  for_each            = {for k,v in var.virtual_networks : k => v}
  address_space       = each.value.address_space
  location            = each.value.location
  name                = each.value.name
  resource_group_name = each.value.resource_group_name
}

resource "azurerm_subnet" "res_subnet" {
  for_each            = {for k,v in var.subnets : k => v}
  address_prefixes     = each.value.address_prefixes
  name                 = each.value.name
  resource_group_name  = each.value.resource_group_name
  virtual_network_name = each.value.virtual_network_name
}

# resource "azurerm_virtual_network" "res-1" {
#   address_space       = ["10.2.0.0/16"]
#   location            = "eastus"
#   name                = "VM1-vnet"
#   resource_group_name = "defender-prueba"
# }

# resource "azurerm_subnet" "res-2" {
#   address_prefixes     = ["10.2.0.0/24"]
#   name                 = "default"
#   resource_group_name  = "defender-prueba"
#   virtual_network_name = "VM1-vnet"
# }

# resource "azurerm_virtual_network" "res-3" {
#   address_space       = ["10.3.0.0/16"]
#   location            = "northeurope"
#   name                = "vm2-vnet"
#   resource_group_name = "defender-prueba"
# }

# resource "azurerm_subnet" "res-4" {
#   address_prefixes     = ["10.3.0.0/24"]
#   name                 = "default"
#   resource_group_name  = "defender-prueba"
#   virtual_network_name = "vm2-vnet"
# }
