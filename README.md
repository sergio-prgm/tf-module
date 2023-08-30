# TFModule | Automatize deployable IaC in azure

## Purpose

Following the steps of [azexport](https://github.com/Azure/aztfexport) and given that it isn't built for production/reading code,
TFModule is a tool that let's you define a module structure and creates the necessary files to deploy your existing infrastructure.

## How to use it

1. Make sure `aztfexport` is already installed in your system.

2. Make a list of desired Resource Groups to import into terraform.

3. Create an empty folder where `aztfexport` will create all the files.

4. Run tfexport (once per each Resource Groups):

```sh
aztfexport rg <rg-name>
```

5. Create a tfmodule.yaml file where the module structure is defined (example file
can be seen [here](./example/tfmodule.yaml)

6. Run the following command specifying the resource group that you want to convert
to Terraform.

```sh
tfmodule desired-rg
```

## Dev Commands

<!-- ```sh
go build ./cmd/main.go && ./main -daily
``` -->

```sh
go build -o ./main ./cmd/main.go && ./main -conf ./example/ -src ./de-pr-08-30/
go build -o .\main.exe .\cmd\main.go && .\main.exe -conf .\example\ -src .\de-pr-08-30\
```

### Cross Compile

```sh
GOOS=windows GOARCH=amd64 go build -o ./bin/win/tfmodule.exe ./cmd/main.go
```

```sh
GOOS=darwin GOARCH=amd64 go build -o ./bin/mac/tfmodule ./cmd/main.go
```

## Roadmap

**Dev** -> Improve code readability, documentation, developer experience, etc.
**Use** -> Improve usability of the tool, output clarity, etc.
**Cod** -> Make changes to the code in general, add features, etc.

- [x] Read yaml config file
- [x] Read and parse main.tf and
  - [x] Modules
  - [x] Main (providers)
- [x] Read multiple tf files (providers, terraform, etc.)
- [x] ask for the source files path in a flag [otherwise default to ./] (yaml & main)
- [x] Update or delete `depends_on` (it causes problems because it links to an external resource as if it was local)
- [x] Change `resource` to `source` in main.tf module declarations
- [x] Put them into the correct modules
- [ ] **Cod** **1** Implement `for_each` in repeated resources
(priority because this affects imports and code structure)
- [ ] **Use** Research a way to make #2 easy/straightforward.
- [ ] **Dev** Improve readibility of code :)
- [ ] **Use** Add flag to just check if all modules in main.tf are represented in yaml file.
But do it anyways
- [ ] **Use** Add csv file with all the existing resources ([csv example](./example/modules.csv))
- [ ] **Cod** Generate `.tfvars` with raw variables
- [ ] **Research** Import blocks

### Extra

- [ ] **Use** Make output more comprehensible (color, verbose description - add flag for this)
- [ ] **Cod** Create the yaml file from `csv` resource output

1. Create "output" & "output/modules" folders
2. Write "main.tf", and the respective module files
3. Bla

### Create main.tf

1. ~~Write providers (from tfexport file)~~
2. ~~Write modules (from yaml config file)~~
3. *Write variables*

### Create modules

1. ~~Write resources (from tfexport file)~~
2. ~~Create *blank* "variables.tf" and *output.tf*~~

```plaintext
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
