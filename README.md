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

## Roadmap

- [x] Read yaml config file
- [x] Read and parse main.tf and
  - [x] Modules
  - [x] Main (providers)
- [ ] Add flag to just check if all modules in main.tf are represented in yaml file.
But do it anyways
- [x] Put them into the correct modules

### Extra

- [ ] Make output more comprehensible (color, verbose description - add flag for this)
- [ ] Suggest (os create) for_each in repeated resources

1. Create "output" & "output/modules" folders
2. Write "main.tf", and the respective module files

### Create main.tf

1. ~~Write providers (from tfexport file)~~
2. ~~Write modules (from yaml config file)~~
3. *Write variables*

### Create modules

1. ~~Write resources (from tfexport file)~~
2. ~~Create *blank* "variables.tf" and *output.tf*~~

├─ output/
├─── main.tf
├─── terraform.tfvars
├─── variables.tf
└─── modules/
├───── StorageAccount/
├─────── main.tf
├─────── outputs.tf
└─────── variables.tf
├───── ResourceGroups/
├─────── main.tf
├─────── outputs.tf
└─────── variables.tf
