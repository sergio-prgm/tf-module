# TFModule | Automatize deployable IaC in azure

## Purpose

Following the steps of [azexport](https://github.com/Azure/aztfexport) and given that it isn't built for production/reading code,
TFModule is a tool that let's you define a module structure and creates the necessary files to deploy your existing infrastructure.

## How to use it 

1. Make sure tfexport is already installed in your system. 

2. Create a tfmodule.yaml file where the module structure is defined (example file
can be seen [here](./example/tfmodule.yaml) 

3. Run the following command specifying the resource group that you want to convert
to Terraform.

```sh
tfmodule desired-rg
```
