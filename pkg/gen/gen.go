package gen

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/sergio-prgm/tf-module/pkg/inout"
)

type VarsContents []map[string]interface{}

// CreateVars creates a structured map[resource_name]contents{}
// to use in tfvars, variables, modules, etc.
func CreateVars(rawResources []string, modules []inout.Modules) map[string]VarsContents {
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
				newResource := ParseResource(blockContent)
				// fmt.Printf("\n%v\n", newResource)
				vars[resourceName] = append(vars[resourceName], newResource)
			}
		}
	}
	return vars
}

// ParseResource creates a structured map[string]interface{}
// Its used to parse the values on the main.tf of aztfexport to generate the terraform.tfvars
func ParseResource(rawResource string) map[string]interface{} {
	var resource map[string]interface{}
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
			content += "\"" + splittedStr[0] + "\" = [\n"
			content += "{\n"
			i, content = insideBracket(stringArr, i, content)
			still_first_string := true
			i += 1
			for still_first_string && i < len(stringArr) {
				v = stringArr[i]
				splittedStr := strings.Split(strings.TrimSpace(v), " ")
				// ainda dentro do bloco
				if last_var == splittedStr[0] {
					content += "\n{\n"
					i, content = insideBracket(stringArr, i, content)
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
			content += "\"" + splittedStr[0] + "\" = {\n"
			i, content = insideBracket(stringArr, i, content)
		} else {
			v = stringArr[i]
			splittedStr := strings.Split(strings.TrimSpace(v), " ")
			content += fmt.Sprintf("\"%s\" %s", splittedStr[0], strings.Join(splittedStr[1:], " "))
			content += ",\n"
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

	jsonedString := "{" + strings.ReplaceAll(content, "=", ":") + "\n}"
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

// insideBracket returnes an int and a string
// Its used to parse the values inside brackets and its nested values
func insideBracket(stringArr []string, i int, content string) (int, string) {
	bracket_count := 1
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
			content += fmt.Sprintf("\"%s\" %s", splittedStr[0], strings.Join(splittedStr[1:], " "))
			content += ","
		}
	}
	return i, content
}

// GenerateImports returnes an string
// It generates the imports blocks for the imports.tf file
func GenerateImports(resources []inout.Resource, modules inout.YamlMapping) string {
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
			formattedResourceType := fmt.Sprintf("module.%s.%s.%ss[\"%d\"]", moduleName, resource.ResourceType, resource_type, index)
			output.WriteString(fmt.Sprintf("import {\n  to = %s\n  id = \"%s\"\n}\n\n", formattedResourceType, resource.ResourceID))
		} else {
			otherOutput.WriteString(fmt.Sprintf("%s\n", resource.ResourceType))
		}
	}

	finalString := output.String()
	path := "./output/imports.tf"
	success := "\nAll the imports where generated correctly!\nFile: " + path
	inout.WriteToFile(finalString, path, success)
	path = "./output/unmapped_resources.txt"
	success = "\nUnmapped Resources:\n" + otherOutput.String() + "\nAll the unmapped_resources where write correctly!\nFile: " + path
	inout.WriteToFile(otherOutput.String(), path, success)

	fmt.Println("Data written to files successfully!")

	return finalString
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
func GenerateModuleYaml(resourcesMapping []inout.Resource, modules_map []inout.ModuleResource) []inout.ModuleResource {
	var resources []inout.ModuleResource
	for _, resource := range resourcesMapping {
		for _, mapped_resource := range modules_map {
			if resource.ResourceType == mapped_resource.ResourceType {
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
