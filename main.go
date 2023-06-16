package main

import (
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

func readTf ( raw []byte) {
	file := string(raw[:])	
	fileLines := strings.Split(file, "\n")
	for i := 0; i < len(fileLines); i++ {
		if fileLines[i] == "}" {
			fmt.Print("<<EOM>>")
		}
		fmt.Printf("%d\t- %s\n", i, fileLines[i])
	}
	// TODO save an isModule variable that resets when first char is }
}

func main() {
	configFile := "./example/tfmodule.yaml"
	tfFile := "./example/main.tf"
	// file := "./easy.yaml"
	conf, err := os.ReadFile(configFile)

	if err != nil {
		log.Fatalf("ERROR: %s doesn't elist", configFile)
	} else {
		fmt.Printf("Reading modules from %s\n", configFile)
	}

	tf, err := os.ReadFile(tfFile)

	if err != nil {
		log.Fatalf("ERROR: %s doesn't elist", tfFile)
	} else {
		fmt.Printf("Reading terraform main from %s\n", tfFile)
	}


	module := F{}
	err = yaml.Unmarshal(conf, &module)

	if err != nil {
		log.Fatal()
	}

	for i := 0; i < len(module.Modules); i++ {
		fmt.Printf("\nmodule: %s\nresources: %v\n", module.Modules[i].Name, module.Modules[i].Resources)
	}

	readTf(tf)
}
