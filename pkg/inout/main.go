// Package inout contains functions that deal with io
// operations, reading and writing the contents of the
// files generated by [scf]
package inout

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sergio-prgm/tf-module/pkg/util"
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

func ReadConfig(fileName string) F {
	filePath := fmt.Sprintf("%s/tfmodule.yaml", strings.TrimSuffix(fileName, "/"))
	conf, err := os.ReadFile(filePath)

	if err != nil {
		log.Fatalf("ERROR: %s doesn't exist", filePath)
	} else {
		fmt.Printf("Reading modules from %s\n", util.EmphasizeStr(filePath, util.Yellow, util.Normal))
	}

	return ParseConfig(conf)
}

// ParseConfig parses the contents of the provided yaml file
// into a usable struct
func ParseConfig(conf []byte) F {
	configModules := F{}
	err := yaml.Unmarshal(conf, &configModules)

	if err != nil {
		log.Fatal()
	}
	for i := 0; i < len(configModules.Modules); i++ {
		fmt.Printf("\nmodule: %s\nresources: %v\n", configModules.Modules[i].Name, configModules.Modules[i].Resources)
	}
	return configModules
}

func ReadTfFiles(fileName string) ParsedTf {
	tfFile := fmt.Sprintf("%s/main.tf", strings.TrimSuffix(fileName, "/"))
	tfFiles, err := os.ReadDir(strings.TrimSuffix(fileName, "/"))
	allTf := []byte("")

	for i := len(tfFiles) - 1; i >= 0; i-- {
		v := tfFiles[i]
		if strings.HasSuffix(v.Name(), ".tf") {
			fmt.Println(v.Name())
			currentTfFile, err := os.ReadFile(fileName + v.Name())
			if err != nil {
				log.Fatal(err)
			}
			allTf = append(allTf, currentTfFile...)
		}
	}

	if err != nil {
		log.Fatalf("ERROR: %s doesn't exist", tfFile)
	} else {
		fmt.Printf("Reading terraform main from %s\n", util.EmphasizeStr(tfFile, util.Yellow, util.Normal))
	}

	return ReadTf(allTf)
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
// func CreateVars(rawResources []string, modules []Modules) map[string][]string {
func CreateVars(rawResources []string, modules []Modules) map[string]VarsContents {
	// var vars map[string][]map[string]interface{}
	// var vars map[string][]string = make(map[string][]string)
	var vars map[string]VarsContents = make(map[string]VarsContents)

	for _, v := range modules {
		for _, resource := range rawResources {
			resoureceArray := strings.Split(resource, "\n")
			rawResourceName := strings.Split(resource, "\"")[1]
			resourceName := strings.Replace(rawResourceName, "azurerm_", "", 1) + "s"

			blockContent := strings.Join(resoureceArray[1:len(resoureceArray)-1], "\n")

			if slices.Contains(v.Resources, rawResourceName) {
				fmt.Println("\nRaw block content:")
				fmt.Println(blockContent)
				newResource := ParseResource(blockContent)
				// fmt.Printf("\n%v\n", newResource)
				vars[resourceName] = append(vars[resourceName], newResource)
			}
		}
	}
	return vars
}

// ParseResource converts the contents of a resource block into a map
func ParseResource(rawResource string) map[string]interface{} {
	var resource map[string]interface{}
	stringArr := strings.Split(rawResource, "\n")
	quotedString := ""
	includesInnerBlock := false

	for i, v := range stringArr {
		splittedStr := strings.Split(strings.TrimSpace(v), " ")
		if slices.Contains(splittedStr, "{") {
			if !strings.Contains(v, "=") {
				splittedStr = slices.Insert(splittedStr, 1, "=")
			}
			includesInnerBlock = true
		}
		fmt.Println(includesInnerBlock)

		if splittedStr[0] == "}" {
			quotedString += "}"
			includesInnerBlock = false
		} else {
			quotedString += fmt.Sprintf("\"%s\" %s", splittedStr[0], strings.Join(splittedStr[1:], " "))
		}

		if includesInnerBlock {
			if !slices.Contains(splittedStr, "{") && strings.TrimSpace(stringArr[i+1]) != "}" {
				quotedString += ",\n"
			} else {
				quotedString += "\n"
			}
		} else if i != len(stringArr)-1 {
			quotedString += ",\n"
		}

		// if includesInnerBlock {
		// 	quotedString += "\n"
		// }
	}

	fmt.Println(quotedString)
	jsonedString := "{" + strings.ReplaceAll(quotedString, "=", ":") + "\n}"
	err := json.Unmarshal([]byte(jsonedString), &resource)
	if err != nil {
		fmt.Println("Here", err)
	}

	return resource
}
