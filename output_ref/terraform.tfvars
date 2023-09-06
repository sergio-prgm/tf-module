// Automatically generated variables
// Should be changed
subnets = [
  {
    address_prefixes : ["10.2.0.0/24"]
    name : "default"
    resource_group_name : "defender-prueba"
    virtual_network_name : "VM1-vnet"
    tags : {}
  },
  {
    address_prefixes : ["10.3.0.0/24"]
    name : "default"
    resource_group_name : "defender-prueba"
    virtual_network_name : "vm2-vnet"
    tags: {}
  }
]

resource_groups = [
  {
    location : "eastus"
    name : "defender-prueba"
    tags : {
      responsable : "Sergio Esteve"
    }
  }
]

storage_accounts = [
  {
    account_kind : "Storage"
    account_replication_type : "LRS"
    account_tier : "Standard"
    allow_nested_items_to_be_public : false
    default_to_oauth_authentication : null
    location : "westeurope"
    name : "cloudshellstacc1"
    resource_group_name : "defender-prueba"
    tags : {
      ms-resource-usage : "azure-cloud-shell"
    }
  },
  {
    account_kind : "Storage"
    account_replication_type : "LRS"
    account_tier : "Standard"
    allow_nested_items_to_be_public : null
    default_to_oauth_authentication : true
    location : "eastus"
    name : "defenderpruebaa10c"
    resource_group_name : "defender-prueba"
    # Añadido para hacer pruebas
    tags: {
      # ms-resource-usage : "azure-cloud-shell"
    }
  }
]

storage_containers = [
  {
    name : "azure-webjobs-hosts"
    storage_account_name : "defenderpruebaa10c"
  },
  {
    name : "azure-webjobs-secrets"
    storage_account_name : "defenderpruebaa10c"
  },
  {
    name : "scm-releases"
    storage_account_name : "defenderpruebaa10c"
  }
]

virtual_networks = [
  {
    address_space : ["10.2.0.0/16"]
    location : "eastus"
    name : "VM1-vnet"
    resource_group_name : "defender-prueba"
  },
  {
    address_space : ["10.3.0.0/16"]
    location : "northeurope"
    name : "vm2-vnet"
    resource_group_name : "defender-prueba"
  }
]

