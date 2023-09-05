package inout

import (
	"encoding/json"
	"log"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type Modules struct {
	Name      string   `yaml:"name"`
	Resources []string `yaml:"resources"`
}

type F struct {
	Modules []Modules `yaml:"modules"`
	Confg   []string  `yaml:"config"`
}

type ParsedTf struct {
	Providers []string
	Resources []string
}

func ReadConfig(conf []byte) F {
	configModules := F{}
	err := yaml.Unmarshal(conf, &configModules)

	if err != nil {
		log.Fatal()
	}
	return configModules
}

// ReadTf reads the contents of the *tf* files produced by aztfexport and
// returns a struct with the providers and resources in it
func ReadTf(raw []byte) ParsedTf {
	file := string(raw[:])
	fileLines := strings.Split(file, "\n")

	isProv := false
	isResource := false
	isBlock := false
	isDependsOn := false

	var rawProv []string
	var rawResource []string

	var currentBlock string

	for i := 0; i < len(fileLines); i++ {

		if !isBlock {

			firstWord := strings.Split(fileLines[i], " ")[0]

			if firstWord == "resource" {
				// fmt.Print("\nStart of resource\n")
				isResource = true
				isBlock = true
			} else if firstWord == "terraform" || firstWord == "provider" {
				// fmt.Print("\nStart of provider/tf\n")
				isBlock = true
				isProv = true
			} else {
				currentBlock = ""
				isBlock = false
			}
		}
		if fileLines[i] == "}" && isBlock {
			if isResource {
				currentBlock += fileLines[i]
				rawResource = append(rawResource, currentBlock)
				isResource = false
				isBlock = false
				currentBlock = ""
			} else if isProv {
				currentBlock += fileLines[i]
				rawProv = append(rawProv, currentBlock)
				isProv = false
				isBlock = false
				currentBlock = ""
			}
		}
		if isBlock {
			if !isDependsOn {
				// if util.FirstWordIs(fileLines[i])
				firstWordInside := strings.Split(strings.TrimSpace(fileLines[i]), " ")[0]

				if firstWordInside == "depends_on" {
					isDependsOn = true
				} else {
					currentBlock += fileLines[i] + "\n"
				}
			} else {
				firstWordInside := strings.Split(strings.TrimSpace(fileLines[i]), " ")[0]
				if firstWordInside == "]" {
					isDependsOn = false
				}
			}
		}

	}
	return ParsedTf{
		Resources: rawResource,
		Providers: rawProv,
	}
}

type VarsContents []map[string]interface{}

// CreateVars creates a structured map[resource_name]contents{}
// to use in tfvars, variables, modules, etc.
func CreateVars(rawResources []string, modules []Modules) map[string][]string {
	// var vars map[string][]map[string]interface{}
	var vars map[string][]string = make(map[string][]string)

	for _, v := range modules {
		for _, resource := range rawResources {
			resoureceArray := strings.Split(resource, "\n")
			rawResourceName := strings.Split(resource, "\"")[1]
			resourceName := strings.Replace(rawResourceName, "azurerm_", "", 1) + "s"

			blockContent := strings.Join(resoureceArray[1:len(resoureceArray)-1], "\n")

			if slices.Contains(v.Resources, rawResourceName) {
				vars[resourceName] = append(vars[resourceName], blockContent)
			}
		}
	}
	return vars
}

// ParseResource converts the contents of a resource block into a map
func ParseResource(rawResource string) map[string]interface{} {
	var resource map[string]interface{}
	json.Unmarshal([]byte(rawResource), &resource)
	return resource
}
