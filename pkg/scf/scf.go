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

func backendVariables(backend inout.BackendConf, is_backend bool, outputResource string, entry_point string, data string) string {
	backend_string := ""
	if backend.Container_name != "" && backend.Key_prefix != "" && backend.Resource_group_name != "" && backend.Storage_account_name != "" {
		if is_backend {
			backend_string = "terraform {\n"
			backend_string += "\tbackend \"azurerm\" {\n"
			backend_string += "\t\tresource_group_name = \"" + backend.Resource_group_name + "\"\n"
			backend_string += "\t\tstorage_account_name = \"" + backend.Storage_account_name + "\"\n"
			backend_string += "\t\tcontainer_name = \"" + backend.Container_name + "\"\n"
			backend_string += "\t\tkey = \"" + backend.Key_prefix + entry_point + ".tfstate\"\n"
			backend_string += "\t}\n"
			backend_string += "}\n"
		} else {
			if !strings.Contains(data, "data \"terraform_remote_state\" \""+outputResource+"\" {\n") {
				backend_string += "data \"terraform_remote_state\" \"" + outputResource + "\" {\n"
				backend_string += "\tbackend = \"azurerm\"\n"
				backend_string += "\tconfig = {\n"
				backend_string += "\t\tresource_group_name = \"" + backend.Resource_group_name + "\"\n"
				backend_string += "\t\tstorage_account_name = \"" + backend.Storage_account_name + "\"\n"
				backend_string += "\t\tcontainer_name = \"" + backend.Container_name + "\"\n"
				backend_string += "\t\tkey = \"" + backend.Key_prefix + "." + entry_point + ".tfstate\"\n"
				backend_string += "\t}\n"
				backend_string += "}\n"
			}
		}
	} else {
		if !is_backend {
			if !strings.Contains(data, "data \"terraform_remote_state\" \""+outputResource+"\" {\n") {
				backend_string += "data \"terraform_remote_state\" \"" + outputResource + "\" {\n"
				backend_string += "\tbackend = \"local\"\n"
				backend_string += "\tconfig = {\n"
				backend_string += "\tpath = \"./../" + entry_point + "/terraform.tfstate\"\n"
				backend_string += "\t}\n"
				backend_string += "}\n"
			}
		}
	}
	return backend_string
}

// CreateMainFiles creates all the files that are general to the
// terraform project and not iniside of the "modules" folder, i.e.:
// main.tf, terraform.tfvars, variables.tf
func CreateMainFiles(mainContent string, varsContent string, tfvarsContent string, backend inout.BackendConf) error {
	backend_string := backendVariables(backend, true, "", "", "")
	if backend_string != "" {
		mainContent = strings.Replace(mainContent, "backend \"local\" {}", "", 1)
	}
	err := os.WriteFile("./output/main.tf",
		[]byte(mainContent),
		os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating main.tf:\n%v", err)
	}

	err = os.WriteFile("./output/backend.tf",
		[]byte(backend_string),
		os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating backend.tf:\n%v", err)
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
func CreateModuleFiles(filePath string, content string, variables string, outputs string) error {
	err := os.WriteFile(filePath+"main.tf",
		[]byte(content),
		os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating %s:\n%v", filePath+"main.tf", err)
	} else {
		fmt.Printf("\n%s created...", filePath+"main.tf")
	}

	err = os.WriteFile(filePath+"output.tf",
		[]byte(outputs),
		os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating %s:\n%v", filePath+"output.tf", err)
	} else {
		fmt.Printf("\n%s created...", filePath+"output.tf")
	}
	variables += "variable common { type = any }\n"
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
func CreateFiles(parsedBlocks inout.ParsedTf, resourceMap map[string]gen.VarsContents, configModules inout.YamlMapping, ep bool) error {
	fmt.Print(util.EmphasizeStr("\nCreating files...", util.Green, util.Bold))
	outputs := GenerateModulesFiles(configModules, resourceMap)

	modulesBlocks := ""
	entry_point_mapp := make(map[string]string)
	output_entry_point := make(map[string]string)
	data_entry_point := make(map[string]string)
	// CreateMain files
	for i, v := range configModules.Modules {
		resourceCall := ""
		//Inside each module ref on main.tf
		for _, r := range v.Resources {
			cleanResource := strings.Replace(r, "azurerm_", "", 1) + "s"
			resourceCall += fmt.Sprintf("\t%s = var.%s\n", cleanResource, cleanResource)
		}

		for _, output := range outputs {
			if v.Name == output.OuputModule {
				entry_point := ""
				for _, val_module := range configModules.Modules {
					if val_module.Name == output.OuptutModuleRef {
						entry_point = val_module.EntryPoint
						break
					}
				}
				data_entry_point[v.EntryPoint] += backendVariables(configModules.Backend, false, output.OputputResource, entry_point, data_entry_point[v.EntryPoint])
				if ep {
					string_to_add := "\t" + output.OputputResource + " = data.terraform_remote_state." + output.OputputResource + ".outputs." + output.OputputResource + "\n"
					if !strings.Contains(resourceCall, string_to_add) {
						resourceCall += string_to_add
					}
				} else {
					string_to_add := "\t" + output.OputputResource + " = module." + output.OuptutModuleRef + "." + output.OputputResource + "\n"
					if !strings.Contains(resourceCall, string_to_add) {
						resourceCall += string_to_add
					}
				}
			}
			if v.Name == output.OuptutModuleRef {
				if !strings.Contains(output_entry_point[v.EntryPoint], "output \""+output.OputputResource+"\" {\n") {
					output_entry_point[v.EntryPoint] += "output \"" + output.OputputResource + "\" {\n"
					output_entry_point[v.EntryPoint] += "\t value = module." + output.OuptutModuleRef + "." + output.OputputResource + "\n"
					output_entry_point[v.EntryPoint] += "}\n\n"
				}
			}
		}
		entry_point_mapp[v.EntryPoint] += fmt.Sprintf(
			"module \"%s\" {\n\tsource = \"./../../Modules/%s\"\n common = var.common\n%s}\n",
			v.Name,
			v.Name,
			resourceCall,
		)
		modulesBlocks += fmt.Sprintf(
			"module \"%s\" {\n\tsource = \"./Modules/%s\"\n common = var.common\n%s}\n",
			v.Name,
			v.Name,
			resourceCall,
		)

		if i != len(configModules.Modules)-1 {
			modulesBlocks += "\n"
		}
	}

	mainContent := strings.Join(parsedBlocks.Providers, "\n\n") + "\n\n" + modulesBlocks

	common_tfvars := "// Automatically generated variables\n// Should be changed\n"
	varsContent := "// Automatically generated variables\n// Should be changed"
	varsContent += "\n\nvariable \"common\" { type = any }"

	common_tfvars += "common = {\n"
	for _, common := range configModules.CommonVars {
		common_tfvars += "\t\"" + common.Name + "\" : [\n"
		for _, val := range common.Value {
			common_tfvars += "\t\t\"" + val + "\",\n"
		}

		common_tfvars += "\t]\n"
	}
	common_tfvars += "}\n"
	tfvarsContent := common_tfvars

	tfvars_entry_point := make(map[string]string)
	vars_entry_point := make(map[string]string)
	for _, module := range configModules.Modules {
		tfvarsContent += "\n// Start of the variables for the module " + module.Name + "\n"
		for _, resource_name := range module.Resources {
			for name, resource := range resourceMap {
				name_compare := name[:len(name)-1]
				resource_name_compare := strings.Replace(resource_name, "azurerm_", "", 1)
				if name_compare == resource_name_compare {
					encodedVar, err := json.MarshalIndent(resource, " ", "  ")
					if err != nil {
						fmt.Println(err)
					}
					tfvars_entry_point[module.EntryPoint] += fmt.Sprintf("%s = %s\n\n", name, string(encodedVar))
					tfvarsContent += fmt.Sprintf("%s = %s\n", name, string(encodedVar))
					if !strings.Contains(varsContent, fmt.Sprintf("\n\nvariable \"%s\" { type = any }", name)) {
						vars_entry_point[module.EntryPoint] += fmt.Sprintf("\n\nvariable \"%s\" { type = any }", name)
						varsContent += fmt.Sprintf("\n\nvariable \"%s\" { type = any }", name)
					}
				}

			}
		}
		tfvarsContent += "// End of the variables for the module " + module.Name + "\n"
	}
	if ep {
		CreateMainFilesEntryPoints(configModules.Backend, entry_point_mapp, vars_entry_point, tfvars_entry_point, common_tfvars+"\n", parsedBlocks, output_entry_point, data_entry_point)
	} else {
		CreateMainFiles(mainContent, varsContent, tfvarsContent, configModules.Backend)
	}
	return nil
}

func CreateMainFilesEntryPoints(backend inout.BackendConf, entry_point_mapp map[string]string, vars_entry_point map[string]string, tfvars_entry_point map[string]string, common_tfvars string, parsedBlocks inout.ParsedTf, output_entry_point map[string]string, data_entry_point map[string]string) error {
	backend_string := ""
	common_tfvars_message := "// Automatically generated variables\n// Should be changed\n"
	common_vars := "// Automatically generated variables\n// Should be changed"
	common_vars += "\n\nvariable \"common\" { type = any }"

	err := os.WriteFile("./output/EntryPoints/0_CommonVars/terraform.tfvars",
		[]byte(common_tfvars),
		os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating main.tf:\n%v", err)
	}

	for key, value := range entry_point_mapp {
		mainContent := strings.Join(parsedBlocks.Providers, "\n\n") + "\n\n" + value
		if backend.Container_name != "" && backend.Key_prefix != "" && backend.Resource_group_name != "" && backend.Storage_account_name != "" {
			mainContent = strings.Replace(mainContent, "backend \"local\" {}", "", 1)
		}

		err := os.WriteFile("./output/EntryPoints/"+key+"/main.tf",
			[]byte(mainContent),
			os.ModePerm)

		if err != nil {
			log.Fatalf("Error creating main.tf:\n%v", err)
		}

		err = os.WriteFile("./output/EntryPoints/"+key+"/output.tf",
			[]byte(output_entry_point[key]),
			os.ModePerm)

		if err != nil {
			log.Fatalf("Error creating output.tf:\n%v", err)
		}

		err = os.WriteFile("./output/EntryPoints/"+key+"/data.tf",
			[]byte(data_entry_point[key]),
			os.ModePerm)

		if err != nil {
			log.Fatalf("Error creating data.tf:\n%v", err)
		}

		err = os.WriteFile("./output/EntryPoints/"+key+"/terraform.tfvars",
			[]byte(common_tfvars_message+tfvars_entry_point[key]),
			os.ModePerm)

		if err != nil {
			log.Fatalf("Error creating terraform.tfvars:\n%v", err)
		}

		err = os.WriteFile("./output/EntryPoints/"+key+"/variables.tf",
			[]byte(common_vars+vars_entry_point[key]),
			os.ModePerm)

		if err != nil {
			log.Fatalf("Error creating variables.tf:\n%v", err)
		}

		backend_string = backendVariables(backend, true, "", "."+key, "")

		err = os.WriteFile("./output/EntryPoints/"+key+"/backend.tf",
			[]byte(backend_string),
			os.ModePerm)

		if err != nil {
			log.Fatalf("Error creating backend.tf:\n%v", err)
		}

		fmt.Print("\n./output/EntryPoints/" + key + "/main.tf created...")
		fmt.Print("\n")
	}

	return nil
}

// GenerateModulesFiles generates the content of the main.tf file for each module
// with is respective resources and variables
func GenerateModulesFiles(configModules inout.YamlMapping, resourceMap map[string]gen.VarsContents) []inout.Outputs {
	var keys_array []string
	var blockInnerKey []inout.BlockInnerKey
	var outputs []inout.Outputs
	content_mapp := make(map[string]string)
	variables_mapp := make(map[string]string)
	outputs_mapp := make(map[string]string)
	for _, v := range configModules.Modules {
		//filePath := fmt.Sprintf("./output/Modules/%s/", v.Name)
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
				//meter tudo dentro de uma funcao e torna la recursiva?
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
										for innerKey, inner_value := range innerMap {
											//deveria ser aqui o bloco dentro de bloco
											if second_block, ok := inner_value.(map[string]interface{}); ok {
												blockInnerKey, outputs = addBlockInsideBlock(key, innerKey, second_block, blockInnerKey, configModules, cleanResource, outputs)
												appendToBlock(blockContents, key, "", "")
											} else {
												access_variable := fmt.Sprintf("%s.value[\"%s\"]", key, innerKey)
												line := ""
												line, outputs = change_resource_id_reference(innerKey, configModules, cleanResource, "\t\t\t", access_variable, outputs)
												appendToBlock(blockContents, key, innerKey, line)
											}
										}
									} else {
										keys_array, block_content, outputs = addBasicModuleField(keys_array, block_content, key, configModules, cleanResource, outputs)
									}
								} else {
									keys_array, block_content, outputs = addBasicModuleField(keys_array, block_content, key, configModules, cleanResource, outputs)
								}
							}
						default:
							keys_array, block_content, outputs = addBasicModuleField(keys_array, block_content, key, configModules, cleanResource, outputs)
						}
					}
				}
				content += block_content
				combinedBlockContents := make(map[string]string)
				for fullKey, line := range blockContents {
					blockKey := strings.Split(fullKey, "-")[0]
					combinedBlockContents[blockKey] += line
				}
				visitedBlocks := make(map[string]string)
				for _, value := range blockInnerKey {
					visitedBlocks[value.InnerKey+" "+value.MainKey] += value.Line
				}
				// Now, wrap each block content in its outer structure
				for blockKey, blockContent := range combinedBlockContents {
					for key, value := range visitedBlocks {
						fmt.Println()
						keys := strings.Split(key, " ")

						if keys[1] == blockKey {
							blockContent += fmt.Sprintf("\t\t\tdynamic \"%s\" {\n\t\t\t\tfor_each = try(%s.value[\"%s\"], []) == [] ? [] : [1]\n\t\t\t\tcontent {\n", keys[0], blockKey, keys[0])
							blockContent += value
							blockContent += "\t\t\t\t}\n\t\t\t}\n"
						}

					}
					fullBlock := fmt.Sprintf("\tdynamic \"%s\" {\n\t\tfor_each = try(each.value.%s, [])\n\t\tcontent {\n%s\t\t}\n\t}\n", blockKey, blockKey, blockContent)
					content += fullBlock
				}
				content += "}\n\n"
			}
		}
		content_mapp[v.Name] = content
		variables_mapp[v.Name] = variables
		//CreateModuleFiles(filePath, content, variables)
	}
	for _, output := range outputs {
		if !strings.Contains(outputs_mapp[output.OuptutModuleRef], "output \""+output.OputputResource+"\" {\n") {
			outputs_mapp[output.OuptutModuleRef] += "output \"" + output.OputputResource + "\" {\n"
			outputs_mapp[output.OuptutModuleRef] += "\t value = azurerm_" + output.OputputResource + ".res_" + output.OputputResource + "s\n}\n"
			variables_mapp[output.OuputModule] += "variable " + output.OputputResource + "{ type = any }\n"
		}

	}
	for _, v := range configModules.Modules {
		filePath := fmt.Sprintf("./output/Modules/%s/", v.Name)
		CreateModuleFiles(filePath, content_mapp[v.Name], variables_mapp[v.Name], outputs_mapp[v.Name])
	}
	return outputs
}

func existsBlockInnerKey(blockInnerKey []inout.BlockInnerKey, mainkey string, key string, innerKey string) bool {
	for _, block := range blockInnerKey {
		if block.MainKey == mainkey && block.InnerKey == key && block.SecondInnerKey == innerKey {
			return true
		}
	}
	return false
}

func addBlockInsideBlock(mainkey string, key string, second_block map[string]interface{}, blockInnerKey []inout.BlockInnerKey, configModules inout.YamlMapping, cleanResource string, outputs []inout.Outputs) ([]inout.BlockInnerKey, []inout.Outputs) {
	content := ""
	for innerKey := range second_block {
		acess_variable := fmt.Sprintf("%s.value.%s.%s", mainkey, key, innerKey)
		line := ""
		line, outputs = change_resource_id_reference(innerKey, configModules, cleanResource, "\t", acess_variable, outputs)
		content += line
		if !existsBlockInnerKey(blockInnerKey, mainkey, key, innerKey) {
			blockInnerKey = append(blockInnerKey, inout.BlockInnerKey{
				MainKey:        mainkey,
				InnerKey:       key,
				SecondInnerKey: innerKey,
				Line:           line,
			})
		}
	}
	//block := fmt.Sprintf("\t\t\tdynamic \"%s\" {\n\t\t\t\tfor_each = try(%s.value[\"%s\"], []) == [] ? [] : [1]\n\t\t\t\tcontent {\n%s\t\t\t\t}\n\t\t\t}\n", key, mainkey, key, content)
	return blockInnerKey, outputs
}

func change_resource_id_reference(key string, configModules inout.YamlMapping, cleanResource string, tabs string, acess_variable string, outputs []inout.Outputs) (string, []inout.Outputs) {
	main_resource := "azurerm_" + cleanResource
	id_resource := strings.Replace(key, "_ids", "", 1)
	id_resource = strings.Replace(id_resource, "_id", "", 1)
	id_resource = strings.Replace(id_resource, "__full__", "", 1)
	this_resource_module := ""
	id_resource_module := ""
	var tryString string
	key = strings.Replace(key, "__full__", "", 1)
	acess_variable = strings.Replace(acess_variable, "__full__", "", 1)
	common_vars := inout.Yaml_mapping.CommonVars
	found_common_name := false
	for _, common := range common_vars {
		found_common_name = false
		if common.Name == id_resource {
			found_common_name = true
			break
		}
	}

	if strings.Contains(key, "_id") {
		resource_key := "azurerm_" + id_resource
		for _, module := range configModules.Modules {
			for _, resource := range module.Resources {
				if resource == main_resource {
					this_resource_module = module.Name
				}

				if resource_key == resource {
					id_resource_module = module.Name
				}
			}
		}

		if this_resource_module == id_resource_module {
			tryString = tabs + key + " = "
			if strings.Contains(key, "_ids") {
				tryString += "[for id in " + acess_variable + ": startswith(id, \"/subscriptions/\") ? id : " + resource_key + ".res_" + id_resource + "s[id].id]\n"
				tryString += "[for id in " + acess_variable + ": " + resource_key + ".res_" + id_resource + "s[id].id]\n"
			} else {
				second_acess_variable := ""
				if strings.Contains(acess_variable, "\"]") {
					second_acess_variable = strings.Replace(acess_variable, "\"]", "__full__\"]", 1)
				} else {
					second_acess_variable = acess_variable + "__full__"
				}
				tryString += "try(" + resource_key + ".res_" + id_resource + "s[" + acess_variable + "].id, try(" + second_acess_variable + ", null))\n"
			}
		} else if this_resource_module == "" || id_resource_module == "" {
			//Um deles nao tem nada, nao mudar a logica que ja tava antes

			tryString = fmt.Sprintf("%s%s = try(%s, null)\n", tabs, key, acess_variable)
		} else {
			tryString = tabs + key + " = "
			if strings.Contains(key, "_ids") {
				tryString += "[for id in " + acess_variable + ": startswith(id, \"/subscriptions/\") ? id : " + "var." + id_resource + "[id].id]\n"
			} else {
				second_acess_variable := ""
				if strings.Contains(acess_variable, "\"]") {
					second_acess_variable = strings.Replace(acess_variable, "\"]", "__full__\"]", 1)
				} else {
					second_acess_variable = acess_variable + "__full__"
				}
				tryString += "try(var." + id_resource + "[" + acess_variable + "].id, try(var." + id_resource + "[" + second_acess_variable + "].id, null))\n"
			}
			/// add to outputs
			to_add := inout.Outputs{
				OuputModule:     this_resource_module,
				OputputResource: id_resource,
				OuptutModuleRef: id_resource_module,
			}
			if !existsOutput(outputs, to_add) {
				outputs = append(outputs, to_add)
			}
		}

	} else {
		if found_common_name {
			second_acess_variable := strings.Replace(acess_variable, "\"]", "__full__\"]", 1)
			tryString = fmt.Sprintf("%s%s = try(var.common.%s[%s], try(%s,null))\n", tabs, key, key, acess_variable, second_acess_variable)
		} else {
			tryString = fmt.Sprintf("%s%s = try(%s, null)\n", tabs, key, acess_variable)
		}
	}

	return tryString, outputs
}

func existsOutput(outputs []inout.Outputs, to_add inout.Outputs) bool {
	for _, val := range outputs {
		if val.OuputModule == to_add.OuputModule && val.OputputResource == to_add.OputputResource && val.OuptutModuleRef == to_add.OputputResource {
			return true
		}
	}
	return false
}

// addBasicModuleField
// adds the basic fields on a resource (e.g name = try(each.value.name, null))
func addBasicModuleField(keys_array []string, block_content string, key string, configModules inout.YamlMapping, cleanResource string, outputs []inout.Outputs) ([]string, string, []inout.Outputs) {
	tryString, outputs := change_resource_id_reference(key, configModules, cleanResource, "\t", "each.value."+key, outputs)
	key = strings.Replace(key, "__full__", "", 1)
	if !stringExists(keys_array, key) {
		keys_array = append(keys_array, key)
		block_content += tryString
	} else {
		tryString = ""
	}
	return keys_array, block_content, outputs
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
func CreateFolders(config inout.YamlMapping, ep bool) {
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
	if ep {
		fmt.Printf("\nCreating %s folder", "0_CommonVars")
		path := fmt.Sprintf("output/EntryPoints/%s", "0_CommonVars")
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range config.Modules {
			fmt.Printf("\nCreating %s folder", v.EntryPoint)
			path := fmt.Sprintf("output/EntryPoints/%s", v.EntryPoint)
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Print("\n")
	}
}
