package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Modules struct {
	Name      string   `yaml:"name"`
	Resources []string `yaml:"resources"`
}

type F struct {
	Modules []Modules `yaml:"modules"`
}

type parsedTf struct {
	providers []string
	modules   []string
}

func StartsWith(line, comp string) (bool, error) {
	if len(line) < len(comp) {
		return false, errors.New("Chars out ouf range!")
	} else {
		for i := 0; i < len(comp); i++ {
			if comp[i] != line[i] {
				return false, nil
			}
		}
		return true, nil
	}
}

func readTf(raw []byte) parsedTf {
	file := string(raw[:])
	fileLines := strings.Split(file, "\n")

	isProv := false
	isModule := false
	isBlock := false

	var rawProv []string
	var rawModules []string

	var currentBlock string
	for i := 0; i < len(fileLines); i++ {
		// if len(fileLines[i]) > 8 && fileLines[i][0:8] == "resource" {
		// 	isModule = true
		// 	isBlock = true
		// }
		if !isBlock {
			fmt.Print("\naint block")
			if fileLines[i][:8] == "resource" {
				isModule = true
				isBlock = true
			}
			// if fileLines[i][:8] == "terrafor" {
			// 	isBlock = true
			// 	isProv = true
			// }
		}
		if fileLines[i] == "}" && isBlock {
			if isModule {
				currentBlock += fileLines[i]
				rawModules = append(rawModules, currentBlock)
				isModule = false
			} else if isProv {
				currentBlock += fileLines[i]
				rawProv = append(rawProv, currentBlock)
				isProv = false
			} else {
				currentBlock = ""
				isBlock = false
				fmt.Printf("\n\n%d:\n%v", len(rawModules), rawModules[len(rawModules)-1])
			}
		}
		if isBlock {
			currentBlock += fileLines[i] + "\n"
		}
	}
	return parsedTf{
		modules:   rawModules,
		providers: rawProv,
	}

	// TODO save an isModule variable that resets when first char is }
}

func main() {
	configFile := "./example/tfmodule.yaml"
	// tfFile := "./example/main.tf"
	// file := "./easy.yaml"
	conf, err := os.ReadFile(configFile)

	if err != nil {
		log.Fatalf("ERROR: %s doesn't elist", configFile)
	} else {
		fmt.Printf("Reading modules from %s\n", configFile)
	}

	// tf, err := os.ReadFile(tfFile)
	//
	// if err != nil {
	// 	log.Fatalf("ERROR: %s doesn't elist", tfFile)
	// } else {
	// 	fmt.Printf("Reading terraform main from %s\n", tfFile)
	// }

	module := F{}
	err = yaml.Unmarshal(conf, &module)

	if err != nil {
		log.Fatal()
	}

	for i := 0; i < len(module.Modules); i++ {
		fmt.Printf("\nmodule: %s\nresources: %v\n", module.Modules[i].Name, module.Modules[i].Resources)
	}
	// fmt.Print(tf)
	// readTf(tf)
	// fmt.Printf("%v", parsed.modules[0])
}
