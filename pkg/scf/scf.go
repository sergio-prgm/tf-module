// Package scf contains functions whose purpose
// is the general scaffolding of the project
package scf

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sergio-prgm/tf-module/pkg/gen"
	"github.com/sergio-prgm/tf-module/pkg/inout"
	"github.com/sergio-prgm/tf-module/pkg/util"
	"golang.org/x/exp/slices"
)

// CreateMainFiles creates all the files that are general to the
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

// CreateModuleFiles creates all the files that are general to the modules
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

// CreateFiles creates the module files containing the resources
// specified in the yaml config file
func CreateFiles(parsedBlocks inout.ParsedTf, resourceMap map[string]gen.VarsContents, configModules inout.YamlMapping) error {
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
	GenerateModulesFiles(configModules, resourceMap)
	return nil
}

// GenerateModulesFiles generates the content of the main.tf file for each module
// with is respective resources and variables
func GenerateModulesFiles(configModules inout.YamlMapping, resourceMap map[string]gen.VarsContents) {
	var keys_array []string
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
				block_content := ""
				keys_array = []string{}
				blockContents := make(map[string]string)
				for _, resourceList := range resource {
					for key, value := range resourceList {
						switch v := value.(type) {
						case []interface{}:
							for _, second_val := range v {
								if innerMap, ok := second_val.(map[string]interface{}); ok {
									if key != "tags" {
										if len(innerMap) == 0 {
											appendToBlock(blockContents, key, "", "")
										}
										for innerKey := range innerMap {
											line := fmt.Sprintf("\t\t\t%s = try(%s.value[\"%s\"], null)\n", innerKey, key, innerKey)
											appendToBlock(blockContents, key, innerKey, line)
										}
									} else {
										keys_array, block_content = addBasicModuleField(keys_array, block_content, key)
									}
								} else {
									keys_array, block_content = addBasicModuleField(keys_array, block_content, key)
								}
							}
						default:
							keys_array, block_content = addBasicModuleField(keys_array, block_content, key)
						}
					}
				}
				content += block_content
				combinedBlockContents := make(map[string]string)
				for fullKey, line := range blockContents {
					blockKey := strings.Split(fullKey, "-")[0]
					combinedBlockContents[blockKey] += line
				}

				// Now, wrap each block content in its outer structure
				for blockKey, blockContent := range combinedBlockContents {
					fullBlock := fmt.Sprintf("\tdynamic \"%s\" {\n\t\tfor_each = try(each.value.%s, [])\n\t\tcontent {\n%s\t\t}\n\t}\n", blockKey, blockKey, blockContent)
					content += fullBlock
				}
				content += "}\n\n"
			}
		}
		CreateModuleFiles(filePath, content, variables)
	}
}

// addBasicModuleField
// adds the basic fields on a resource (e.g name = try(each.value.name, null))
func addBasicModuleField(keys_array []string, block_content string, key string) ([]string, string) {
	attributeString := fmt.Sprintf("\t%s = each.value.%s\n", key, key)
	tryString := fmt.Sprintf("\t%s = try(each.value.%s, null)\n", key, key)
	if !stringExists(keys_array, key) {
		keys_array = append(keys_array, key)
		block_content += tryString
	} else {
		block_content = strings.Replace(block_content, tryString, attributeString, 1)
	}
	return keys_array, block_content
}

// stringExists
// it checks if a respective string exists in an array of strings
func stringExists(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// appendToBlock
// it appends an block key to a map of strings to keep track of wich keys already exists
func appendToBlock(blockMap map[string]string, blockKey, innerKey, line string) {
	fullKey := blockKey + "-" + innerKey // unique key for each inner attribute
	if _, exists := blockMap[fullKey]; !exists {
		blockMap[fullKey] = line
	}
}

// createFolders creates all the necessary folders with the information outlined
// in the yaml config file
func CreateFolders(config inout.YamlMapping) {
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
