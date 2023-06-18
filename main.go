package main

import (
	"fmt"
	// "github.com/sergio-prgm/tf-module/utils"
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

		if !isBlock {

			firstWord := strings.Split(fileLines[i], " ")[0]

			if firstWord == "resource" {
				fmt.Print("\nStart of resource\n")
				isModule = true
				isBlock = true
			} else if firstWord == "terraform" || firstWord == "provider" {
				fmt.Print("\nStart of provider/tf\n")
				isBlock = true
				isProv = true
			} else {
				fmt.Print("\nBlank space\n")
				currentBlock = ""
				isBlock = false
			}
		}
		if fileLines[i] == "}" && isBlock {
			if isModule {
				currentBlock += fileLines[i]
				rawModules = append(rawModules, currentBlock)
				isModule = false
				isBlock = false
			} else if isProv {
				currentBlock += fileLines[i]
				rawProv = append(rawProv, currentBlock)
				isProv = false
				isBlock = false
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
}

func SaveModules(parsed parsedTf) error {
	_, err := os.Stat("output")
	if os.IsNotExist(err) {

		fmt.Print("Creating folder...")

		err = os.Mkdir("output", os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating dir:\n%v", err)
		}
	}

	err = os.WriteFile("./output/main.tf",
		[]byte(strings.Join(parsed.providers, "\n\n")),
		os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating main.tf:\n%v", err)
	}

	fmt.Print("\noutput/main.tf created...")
	fmt.Print("\n")
	return nil
}

func createModuleFiles(config F) {
	fmt.Print("\nRunning createModuleFiles\n")
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

func main() {
	configFile := "./example/tfmodule.yaml"
	tfFile := "./example/main.tf"
	conf, err := os.ReadFile(configFile)

	if err != nil {
		log.Fatalf("ERROR: %s doesn't exist", configFile)
	} else {
		fmt.Printf("Reading modules from %s\n", configFile)
	}

	tf, err := os.ReadFile(tfFile)

	if err != nil {
		log.Fatalf("ERROR: %s doesn't exist", tfFile)
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
	result := readTf(tf)

	// fmt.Printf("Providers length: %d\n", len(result.providers))
	// fmt.Printf("Providers: %v\n", result.providers)
	// fmt.Printf("Modules length: %d\n", len(result.modules))
	// fmt.Printf("Modules: %v\n", result.modules)
	createModuleFiles(module)
	err = SaveModules(result)
	if err != nil {
		log.Fatal(err)
	}
}
