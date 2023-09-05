terraform {
  backend "local" {}
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "3.56.0"
    }
  }
}

provider "azurerm" {
  features {
  }
}

module "ResourceGroup" {
  source          = "./Modules/ResourceGroup"
  resource_groups = var.resource_groups
}

module "StorageAccount" {
  source             = "./Modules/StorageAccount"
  storage_accounts   = var.storage_accounts
  storage_containers = var.storage_containers
}

module "Network" {
  source           = "./Modules/Network"
  virtual_networks = var.virtual_networks
  subnets          = var.subnets
}
