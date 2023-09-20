package gen

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/sergio-prgm/tf-module/pkg/inout"
)

type VarsContents []map[string]interface{}

var Not_Found_resources []inout.UnmappedOutputs
var Found_resources []inout.UnmappedOutputs

// CreateVars creates a structured map[resource_name]contents{}
// to use in tfvars, variables, modules, etc.
func CreateVars(rawResources []string, modules []inout.Modules, mapped_imports []inout.Imports) map[string]VarsContents {
	var vars map[string]VarsContents = make(map[string]VarsContents)
	for _, v := range modules {
		for _, resource := range rawResources {
			resoureceArray := strings.Split(resource, "\n")
			rawResourceName := strings.Split(resource, "\"")[1]
			resourceName := strings.Replace(rawResourceName, "azurerm_", "", 1) + "s"

			blockContent := strings.Join(resoureceArray[1:len(resoureceArray)-1], "\n")

			if slices.Contains(v.Resources, rawResourceName) {
				// fmt.Println("\nRaw block content:")
				// fmt.Println(blockContent)
				newResource := ParseResource(blockContent, mapped_imports, resourceName, modules)
				// fmt.Printf("\n%v\n", newResource)
				vars[resourceName] = append(vars[resourceName], newResource)
			}
		}
	}
	fmt.Println("")
	return vars
}

// ParseResource creates a structured map[string]interface{}
// Its used to parse the values on the main.tf of aztfexport to generate the terraform.tfvars
func ParseResource(rawResource string, mapped_imports []inout.Imports, resourceName string, modules []inout.Modules) map[string]interface{} {
	var resource map[string]interface{}
	content_to_add := ""
	content := ""
	//separa tudo por linhas
	stringArr := strings.Split(rawResource, "\n")
	//percorre linha a linha
	i := 0
	last_var := ""
	v := ""
	for i < len(stringArr) {
		v = stringArr[i]
		//split da linha por espacos
		splittedStr := strings.Split(strings.TrimSpace(v), " ")
		//Dentro de um block
		if slices.Contains(splittedStr, "{") && !slices.Contains(splittedStr, "=") {
			last_var = splittedStr[0]
			content += "\"" + splittedStr[0] + "\" : [\n"
			content += "{\n"
			i, content = insideBracket(stringArr, i, content, mapped_imports, resourceName, modules)
			still_first_string := true
			i += 1
			for still_first_string && i < len(stringArr) {
				v = stringArr[i]
				splittedStr := strings.Split(strings.TrimSpace(v), " ")
				// ainda dentro do bloco
				if last_var == splittedStr[0] {
					content += "\n{\n"
					i, content = insideBracket(stringArr, i, content, mapped_imports, resourceName, modules)
					i++
				} else {
					if content[len(content)-1] == ',' {
						content = content[:len(content)-1]
						content += "\n"
					}
					content += "],\n"
					still_first_string = false
					i--
				}
			}
		} else if slices.Contains(splittedStr, "{") {
			content += "\"" + splittedStr[0] + "\" : {\n"
			i, content = insideBracket(stringArr, i, content, mapped_imports, resourceName, modules)
		} else {
			v = stringArr[i]
			splittedStr := strings.Split(strings.TrimSpace(v), " ")
			content_to_add = addContent(splittedStr, mapped_imports, ",", resourceName, modules)
			content += content_to_add
		}
		i++
	}

	if countChar(content, '[') != countChar(content, ']') {
		content = content[:len(content)-1]
		content += "\n],\n"
	}
	if content[len(content)-1] == '\n' {
		content = content[:len(content)-1]
	}
	if content[len(content)-1] == ',' {
		content = content[:len(content)-1]
		content += "\n"
	}

	jsonedString := "{" + strings.ReplaceAll(content, " = ", ":") + "\n}"
	err := json.Unmarshal([]byte(jsonedString), &resource)
	if err != nil {
		fmt.Println(jsonedString)
		fmt.Println("Here", err)
	}

	return resource
}

// countChar returnes an int
// It takes a string and a char and count the number of times that char appears in the string
func countChar(s string, char rune) int {
	count := 0
	for _, c := range s {
		if c == char {
			count++
		}
	}
	return count
}

func addUnmappedResource(resources []inout.UnmappedOutputs, unfound_resource inout.UnmappedOutputs) []inout.UnmappedOutputs {
	exists := false
	for _, resource := range resources {
		if resource.ResourceName == unfound_resource.ResourceName && resource.ResourceVariable == unfound_resource.ResourceVariable {
			exists = true
			continue
		}
	}
	if !exists {
		resources = append(resources, unfound_resource)
	}
	return resources
}

func findResourceInYaml(resourceName string, modules []inout.Modules) bool {
	for _, module := range modules {
		for _, resource := range module.Resources {
			if resource == "azurerm_"+resourceName {
				return true
			}
		}
	}
	return false
}

func addContent(splittedStr []string, mapped_imports []inout.Imports, end_string string, resourceName string, modules []inout.Modules) string {
	string_to_join := ""
	content := ""
	its_id := false
	if strings.Contains(splittedStr[0], "_id") && !strings.Contains(splittedStr[0], "_ids") {
		its_id = true
		found_resource := findResourceInYaml(strings.Replace(splittedStr[0], "_id", "", 1), modules)
		found_id := false
		for _, mapped_import := range mapped_imports {
			temp_string := fmt.Sprintf("\"%s\" %s", splittedStr[0], strings.Join(splittedStr[1:], " "))
			tempo_strings := strings.Split(temp_string, " = ")
			if "\""+mapped_import.Resource_id+"\"" == tempo_strings[1] {
				string_to_join = "\"" + strconv.FormatInt(int64(mapped_import.Resource_key), 10) + "\""
				found_id = true
				continue
			}
		}
		if !found_id || !found_resource {
			string_to_join = ""
			unfound_resource := inout.UnmappedOutputs{
				ResourceName:     resourceName,
				ResourceVariable: splittedStr[0],
			}
			Not_Found_resources = addUnmappedResource(Not_Found_resources, unfound_resource)
		}
		if found_id && found_resource {
			string_to_join = ""
			found_resource := inout.UnmappedOutputs{
				ResourceName:     resourceName,
				ResourceVariable: splittedStr[0],
			}
			Found_resources = addUnmappedResource(Found_resources, found_resource)
		}
		if !found_resource {
			its_id = false
		}
	} else if strings.Contains(splittedStr[0], "_ids") {
		its_id = true
		resources_temp_string := fmt.Sprintf("\"%s\" %s", splittedStr[0], strings.Join(splittedStr[1:], " "))
		resources_tempo_strings := strings.Split(resources_temp_string, " = ")
		resources_tempo_strings[1] = resources_tempo_strings[1][1 : len(resources_tempo_strings[1])-1]
		resources_ids := strings.Split(resources_tempo_strings[1], ",")
		found_resource := findResourceInYaml(strings.Replace(splittedStr[0], "_ids", "", 1), modules)
		for _, resource_id := range resources_ids {
			found_id := false
			for _, mapped_import := range mapped_imports {
				if "\""+mapped_import.Resource_id+"\"" == resource_id {
					string_to_join += "\"" + strconv.FormatInt(int64(mapped_import.Resource_key), 10) + "\", "
					found_id = true
				}
			}
			if !found_id {
				string_to_join += resource_id + ", "
				unfound_resource := inout.UnmappedOutputs{
					ResourceName:     resourceName,
					ResourceVariable: splittedStr[0],
				}
				Not_Found_resources = addUnmappedResource(Not_Found_resources, unfound_resource)
			}
			if !found_resource {
				string_to_join = ""
				unfound_resource := inout.UnmappedOutputs{
					ResourceName:     resourceName,
					ResourceVariable: splittedStr[0],
				}
				Not_Found_resources = addUnmappedResource(Not_Found_resources, unfound_resource)
				continue
			}
			if found_id && found_resource {
				string_to_join = ""
				found_resource := inout.UnmappedOutputs{
					ResourceName:     resourceName,
					ResourceVariable: splittedStr[0],
				}
				Found_resources = addUnmappedResource(Found_resources, found_resource)
			}
		}
		if string_to_join != "" {
			string_to_join = "[" + string_to_join[:len(string_to_join)-2] + "]"

		}
	}
	if string_to_join == "" {
		if its_id && !strings.Contains(splittedStr[0], "_ids") {
			content += fmt.Sprintf("\"%s\" %s", splittedStr[0]+"__full__", strings.Join(splittedStr[1:], " "))
		} else {
			content += fmt.Sprintf("\"%s\" %s", splittedStr[0], strings.Join(splittedStr[1:], " "))
		}
	} else {
		content += fmt.Sprintf("\"%s\" %s", splittedStr[0], " : "+string_to_join)
	}
	content += end_string
	return content
}

// insideBracket returnes an int and a string
// Its used to parse the values inside brackets and its nested values
func insideBracket(stringArr []string, i int, content string, mapped_imports []inout.Imports, resourceName string, modules []inout.Modules) (int, string) {
	bracket_count := 1
	content_to_add := ""
	for bracket_count > 0 {
		i++
		v := stringArr[i]
		splittedStr := strings.Split(strings.TrimSpace(v), " ")
		if slices.Contains(splittedStr, "{") {
			content += "\"" + splittedStr[0] + "\" : {\n"
			bracket_count += 1
		} else if slices.Contains(splittedStr, "}") {
			if content[len(content)-1] == ',' {
				content = content[:len(content)-1]
				content += "\n"
			}
			content += "},"
			bracket_count -= 1
		} else {
			content += "\n"
			content_to_add = addContent(splittedStr, mapped_imports, ",", resourceName, modules)
			content += content_to_add
		}
	}
	return i, content
}

// GenerateImports returnes an string
// It generates the imports blocks for the imports.tf file
func GenerateImports(resources []inout.Resource, modules inout.YamlMapping) (string, []inout.Imports, []string) {
	var unmappedResources []string
	var imports_mapping []inout.Imports
	resourceModuleMapping := make(map[string]string)
	for _, module := range modules.Modules {
		for _, resourceType := range module.Resources {
			resourceModuleMapping[resourceType] = module.Name
		}
	}

	typeCounter := make(map[string]int)

	var output, otherOutput strings.Builder
	for _, resource := range resources {
		index, exists := typeCounter[resource.ResourceType]
		if !exists {
			typeCounter[resource.ResourceType] = 1
		} else {
			typeCounter[resource.ResourceType] = index + 1
		}

		moduleName, found := resourceModuleMapping[resource.ResourceType]
		resource_type := strings.Replace(resource.ResourceType, "azurerm", "res", 1)
		if found {
			imports_mapping = append(imports_mapping, inout.Imports{
				Resource_key: index,
				Resource_id:  resource.ResourceID,
			})
			formattedResourceType := fmt.Sprintf("module.%s.%s.%ss[\"%d\"]", moduleName, resource.ResourceType, resource_type, index)
			output.WriteString(fmt.Sprintf("import {\n  to = %s\n  id = \"%s\"\n}\n\n", formattedResourceType, resource.ResourceID))
		} else {
			if !slices.Contains(unmappedResources, resource.ResourceType) {
				unmappedResources = append(unmappedResources, resource.ResourceType)
			}
			otherOutput.WriteString(fmt.Sprintf("%s\n", resource.ResourceType))
		}
	}

	finalString := output.String()
	path := "./output/imports.tf"
	success := "\nAll the imports where generated correctly!\nFile: " + path
	inout.WriteToFile(finalString, path, success)
	return finalString, imports_mapping, unmappedResources
}

// AddResource
// It adds the respective mapping of a resource to his module (if doesn't already exist)
func AddResource(resources *[]inout.ModuleResource, item inout.ModuleResource) {
	// Create a map to check for existing items.
	existing := make(map[string]bool)

	// Populate the map based on the current items.
	for _, resource := range *resources {
		key := resource.Module + "|" + resource.ResourceType
		existing[key] = true
	}

	// Check if the item exists.
	key := item.Module + "|" + item.ResourceType
	if !existing[key] {
		*resources = append(*resources, item)
	}
}

// GenerateModuleYaml returns an []inout.ModuleResource
// It generates the mapping of the resources with the respective module to generate the yaml file
// The first one is the json file, and the second argument is the csv
func GenerateModuleYaml(resourcesMapping []inout.Resource, modules_map []inout.ModuleResource) []inout.ModuleResource {
	var resources []inout.ModuleResource
	for _, resource := range resourcesMapping {
		for _, mapped_resource := range modules_map {
			if mapped_resource.Module != "" && resource.ResourceType == mapped_resource.ResourceType {
				AddResource(&resources, mapped_resource)
			}
		}
	}
	return resources
}

// CheckResources
// Checks if the resources in the resourceGroup have it's mapping on the yaml file and counts them
func CheckResources(resources []inout.Resource, mapped_resources inout.YamlMapping) []inout.CsvResources {
	resourcesCsv := []inout.CsvResources{}
	resource_exists := false
	for _, resource := range resources {
		for _, module := range mapped_resources.Modules {
			for _, yaml_resource := range module.Resources {
				if resource.ResourceType == yaml_resource {
					resource_exists = true
					resourcesCsv = addToStruct(resource.ResourceType, module.Name, resourcesCsv)
					continue
				}
			}
			if resource_exists {
				continue
			}
		}
		if !resource_exists {
			resourcesCsv = addToStruct(resource.ResourceType, "", resourcesCsv)
		}
		resource_exists = false
	}
	return resourcesCsv
}

// addToStruct
// it adds a resource to the []CsvResources structure or increments its quantity if it already exists
func addToStruct(resource string, module string, structure []inout.CsvResources) []inout.CsvResources {
	for i, csv_resource := range structure {
		if resource == csv_resource.Resource {
			structure[i].Quantity += 1
			return structure
		}
	}
	return append(structure, inout.CsvResources{
		Resource: resource,
		Module:   module,
		Quantity: 1,
	})
}
