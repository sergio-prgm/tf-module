package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Modules struct {
	Name      string   `yaml:"name"`
	Resources []string `yaml:"resources"`
}

type F struct {
	Modules []Modules `yaml:"modules"`
}

func main() {
	file := "./example/tfmodule.yaml"
	// file := "./easy.yaml"
	f, err := os.ReadFile(file)

	if err != nil {
		log.Fatalf("ERROR: %s doesn't elist", file)
	} else {
		fmt.Printf("Reading modules from %s\n", file)
	}

	module := F{}
	err = yaml.Unmarshal(f, &module)

	if err != nil {
		log.Fatal()
	}

	for i := 0; i < len(module.Modules); i++ {
		fmt.Printf("\nmodule: %s\nresources: %v\n", module.Modules[i].Name, module.Modules[i].Resources)
	}
}
