// Package scf contains functions whose purpose
// is the general scaffolding of the project
package scf

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sergio-prgm/tf-module/pkg/inout"
	"github.com/sergio-prgm/tf-module/pkg/util"
	"golang.org/x/exp/slices"
)

// createMainFiles creates all the files that are general to the
// terraform project and not iniside of the "modules" folder, i.e.:
// main.tf, terraform.tfvars, variables.tf
func CreateMainFiles(mainContent string, varsContent string, tfvarsContent string) error {
	err := os.WriteFile("./output/main.tf",
		[]byte(mainContent),
		os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating main.tf:\n%v", err)
	}

	err = os.WriteFile("./output/terraform.tfvars",
		[]byte(tfvarsContent),
		os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating terraform.tfvars:\n%v", err)
	}

	err = os.WriteFile("./output/variables.tf",
		[]byte(varsContent),
		os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating variables.tf:\n%v", err)
	}

	fmt.Print("\noutput/main.tf created...")
	fmt.Print("\n")

	return nil
}

func CreateModuleFiles(filePath string, content string, variables string) error {
	err := os.WriteFile(filePath+"main.tf",
		[]byte(content),
		os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating %s:\n%v", filePath+"main.tf", err)
	} else {
		fmt.Printf("\n%s created...", filePath+"main.tf")
	}

	_, err = os.Create(filePath + "output.tf")
	if err != nil {
		log.Fatalf("Error creating %s:\n%v", filePath+"output.tf", err)
	} else {
		fmt.Printf("\n%s created...", filePath+"output.tf")
	}

	err = os.WriteFile(filePath+"variables.tf",
		[]byte(variables),
		os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating %s:\n%v", filePath+"variables.tf", err)
	} else {
		fmt.Printf("\n%s created...", filePath+"variables.tf")
	}
	fmt.Println()
	return nil
}

// createFiles creates the module files containing the resources
// specified in the yaml config file
func CreateFiles(parsedBlocks inout.ParsedTf, resourceMap map[string]inout.VarsContents, configModules inout.F) error {
	fmt.Print(util.EmphasizeStr("\nCreating files...", util.Green, util.Bold))

	modulesBlocks := ""

	// CreateMain files
	for i, v := range configModules.Modules {
		resourceCall := ""
		for _, r := range v.Resources {
			cleanResource := strings.Replace(r, "azurerm_", "", 1) + "s"
			resourceCall += fmt.Sprintf("\t%s = var.%s\n", cleanResource, cleanResource)
		}

		modulesBlocks += fmt.Sprintf(
			"module \"%s\" {\n\tsource = \"./Modules/%s\"\n%s}\n",
			v.Name,
			v.Name,
			resourceCall,
		)

		if i != len(configModules.Modules)-1 {
			modulesBlocks += "\n"
		}
	}

	mainContent := strings.Join(parsedBlocks.Providers, "\n\n") + "\n\n" + modulesBlocks

	tfvarsContent := "// Automatically generated variables\n// Should be changed\n"
	varsContent := "// Automatically generated variables\n// Should be changed"
	for name, resource := range resourceMap {
		encodedVar, err := json.MarshalIndent(resource, " ", "  ")
		if err != nil {
			fmt.Println(err)
		}

		tfvarsContent += fmt.Sprintf("%s = %s\n", name, string(encodedVar))
		varsContent += fmt.Sprintf("\n\nvariable \"%s\" { type = any }", name)
	}

	CreateMainFiles(mainContent, varsContent, tfvarsContent)

	// CreateModuleFiles
	for _, v := range configModules.Modules {
		filePath := fmt.Sprintf("./output/Modules/%s/", v.Name)
		variables := ""
		content := ""
		for resourceName, resource := range resourceMap {
			cleanResource := resourceName[:len(resourceName)-1]

			if slices.Contains(v.Resources, "azurerm_"+cleanResource) {
				newVar := fmt.Sprintf("variable %s { type = any }\n", strings.Replace(resourceName, "azurerm_", "", 1))
				if !strings.Contains(variables, newVar) {
					variables += newVar
				}

				content += fmt.Sprintf("resource \"azurerm_%s\" \"res_%s\" {\n", cleanResource, resourceName)
				content += fmt.Sprintf("\tfor_each = {for k, v in var.%s : k => v}\n", resourceName)

				for _, resourceList := range resource {
					for attribute := range resourceList {
						// if length == 1 obviar la linea de try
						// repasar esto porque no tiene buena pinta
						fmt.Println(attribute)
						attributeString := fmt.Sprintf("\t%[1]s = each.value.%[1]s\n", attribute)
						tryString := fmt.Sprintf("\t%[1]s = try(each.value.%[1]s, null)\n", attribute)
						if !strings.Contains(content, attribute) {
							content += tryString
						} else {
							content = strings.Replace(content, tryString, attributeString, 1)
						}
					}
				}
				content += "}\n\n"
			}
		}
		CreateModuleFiles(filePath, content, variables)
	}
	return nil
}

// createFolders creates all the necessary folders with the information outlined
// in the yaml config file
func CreateFolders(config inout.F) {
	fmt.Print(util.EmphasizeStr("\nCreating folders...", util.Green, util.Bold))
	_, err := os.Stat("output")

	if os.IsNotExist(err) {

		fmt.Print("\nCreating folders...")

		err = os.Mkdir("output", os.ModePerm)
		if err != nil {
			log.Fatalf("\nError creating output dir:\n%v", err)
		}
		err = os.MkdirAll("output/Modules", os.ModePerm)
		if err != nil {
			log.Fatalf("\nError creating Modules dir:\n%v", err)
		}
	} else {
		fmt.Print("\n'output' folder already exists.")
	}

	for _, v := range config.Modules {
		fmt.Printf("\nCreating %s folder", v.Name)
		path := fmt.Sprintf("output/Modules/%s", v.Name)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Print("\n")
}
