package main

import (
	"fmt"
	// "github.com/sergio-prgm/tf-module/utils"
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
				// fmt.Print("\nStart of resource\n")
				isModule = true
				isBlock = true
			} else if firstWord == "terraform" || firstWord == "provider" {
				// fmt.Print("\nStart of provider/tf\n")
				isBlock = true
				isProv = true
			} else {
				// fmt.Print("\nBlank space\n")
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

func createModuleFiles(parsedBlocks parsedTf, configModules F) error {
	fmt.Print(util.EmphasizeStr("\nCreating files...", util.Green, util.Bold))

	modulesBlocks := ""

	for i, v := range configModules.Modules {

		modulesBlocks += fmt.Sprintf(
			"module \"%s\" {\n\tresource = \"./Modules/%s\"\n}\n",
			v.Name,
			v.Name,
		)

		if i != len(configModules.Modules)-1 {
			modulesBlocks += "\n"
		}
	}

	mainContent := strings.Join(parsedBlocks.providers, "\n\n") + "\n\n" + modulesBlocks

	err := os.WriteFile("./output/main.tf",
		[]byte(mainContent),
		os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating main.tf:\n%v", err)
	}

	fmt.Print("\noutput/main.tf created...")
	fmt.Print("\n")

	for _, v := range configModules.Modules {
		filePath := fmt.Sprintf("./output/Modules/%s/", v.Name)
		content := ""
		for _, module := range parsedBlocks.modules {
			resourceName := strings.Split(module, "\"")[1]
			if slices.Contains(v.Resources, resourceName) {
				if content == "" {
					content = module
				} else {
					content = content + "\n\n" + module
				}
			}
		}
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

		_, err = os.Create(filePath + "variables.tf")
		if err != nil {
			log.Fatalf("Error creating %s:\n%v", filePath+"variables.tf", err)
		} else {
			fmt.Printf("\n%s created...", filePath+"variables.tf")
		}
		fmt.Println()

	}

	return nil
}

func createFolders(config F) {
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

// validateModules checks whether the information in the config file matches the
// contents of the main.tf and prompts the user the information
func validateModules(configFile F, parsedFile parsedTf) bool {
	return false
}

// updateConfig allows the user to modify the contents of the config file to
// accommodate late changes or forgotten modules/resources/etc.
func updateConfig() error {
	return nil
}

func main() {
	// TODO ask for the source files path in a flag [otherwise default to ./] (yaml & main)
	configFile := "./example/tfmodule.yaml"
	tfFile := "./example/main.tf"
	conf, err := os.ReadFile(configFile)

	if err != nil {
		log.Fatalf("ERROR: %s doesn't exist", configFile)
	} else {
		fmt.Printf("Reading modules from %s\n", util.EmphasizeStr(configFile, util.Yellow, util.Normal))
	}

	tf, err := os.ReadFile(tfFile)

	if err != nil {
		log.Fatalf("ERROR: %s doesn't exist", tfFile)
	} else {
		fmt.Printf("Reading terraform main from %s\n", util.EmphasizeStr(tfFile, util.Yellow, util.Normal))
	}

	configModules := F{}
	err = yaml.Unmarshal(conf, &configModules)

	if err != nil {
		log.Fatal()
	}

	for i := 0; i < len(configModules.Modules); i++ {
		fmt.Printf("\nmodule: %s\nresources: %v\n", configModules.Modules[i].Name, configModules.Modules[i].Resources)
	}
	parsedBlocks := readTf(tf)

	// fmt.Printf("Providers length: %d\n", len(result.providers))
	// fmt.Printf("Providers: %v\n", result.providers)
	// fmt.Printf("Modules length: %d\n", len(result.modules))
	// fmt.Printf("Modules: %v\n", result.modules)
	createFolders(configModules)
	err = createModuleFiles(parsedBlocks, configModules)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Print(util.EmphasizeStr("Emphasize str\n", util.Blue, util.Normal))
	// fmt.Print(util.EmphasizeStr("Emphasize str\n", util.Blue, util.Bold))
}
