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

## Dev Commands

```sh
go build ./cmd/main.go && ./main -daily
```

### Cross Compile

```sh
GOOS=windows GOARCH=amd64 go build -o ./bin/win/tfmodule.exe ./cmd/main.go
```

```sh
GOOS=darwin GOARCH=amd64 go build -o ./bin/mac/tfmodule ./cmd/main.go
```

## Roadmap

- [x] Read yaml config file
- [x] Read and parse main.tf and
  - [x] Modules
  - [x] Main (providers)
- [ ] Read multiple tf files (providers, terraform, etc.)
- [ ] ask for the source files path in a flag [otherwise default to ./] (yaml & main)
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

```
output
├── Modules
│   ├── Network
│   │   ├── main.tf
│   │   ├── output.tf
│   │   └── variables.tf
│   ├── ResourceGroup
│   │   ├── main.tf
│   │   ├── output.tf
│   │   └── variables.tf
│   └── StorageAccount
│       ├── main.tf
│       ├── output.tf
│       └── variables.tf
└── main.tf

5 directories, 10 files
```
