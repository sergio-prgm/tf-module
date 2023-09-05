package main

import (
	"encoding/json"
	"flag"
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
	resources []string
}

// readTf reads the contents of the *tf* files produced by aztfexport and
// returns a struct with the providers and resources in it
func readTf(raw []byte) parsedTf {
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
	return parsedTf{
		resources: rawResource,
		providers: rawProv,
	}
}

// createMainFiles creates all the files that are general to the
// terraform project and not iniside of the "modules" folder, i.e.:
// main.tf, tf.vars, variables.tf
func createMainFiles(mainContent string, varsContent string, tfvarsContent string) error {
	err := os.WriteFile("./output/main.tf",
		[]byte(mainContent),
		os.ModePerm)

	if err != nil {
		log.Fatalf("Error creating main.tf:\n%v", err)
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

func createModuleFiles(filePath string, content string) error {
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
	return nil
}

// createFiles creates the module files containing the resources
// specified in the yaml config file
func createFiles(parsedBlocks parsedTf, varsContent string, tfvarsContent string, configModules F) error {
	fmt.Print(util.EmphasizeStr("\nCreating files...", util.Green, util.Bold))

	modulesBlocks := ""

	for i, v := range configModules.Modules {

		modulesBlocks += fmt.Sprintf(
			"module \"%s\" {\n\tsource = \"./Modules/%s\"\n}\n",
			v.Name,
			v.Name,
		)

		if i != len(configModules.Modules)-1 {
			modulesBlocks += "\n"
		}
	}

	mainContent := strings.Join(parsedBlocks.providers, "\n\n") + "\n\n" + modulesBlocks
	createMainFiles(mainContent, varsContent, tfvarsContent)

	// use in createVars
	for _, v := range configModules.Modules {
		filePath := fmt.Sprintf("./output/Modules/%s/", v.Name)
		content := ""
		for _, resource := range parsedBlocks.resources {
			resourceName := strings.Split(resource, "\"")[1]
			if slices.Contains(v.Resources, resourceName) {
				if content == "" {
					content = resource
				} else {
					content = content + "\n\n" + resource
				}
			}
		}
		// call createModuleFiles
		createModuleFiles(filePath, content)
	}
	return nil
}

// createFolders creates all the necessary folders with the information outlined
// in the yaml config file
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

type VarsContents []map[string]interface{}

// createVars creates a structured map[resource_name]contents{}
// to use in tfvars, variables, modules, etc.
func createVars(rawResources []string, modules []Modules) map[string][]string {
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

// parseResource converts the contents of a resource block into a map
func parseResource(rawResource string) map[string]interface{} {
	var resource map[string]interface{}
	json.Unmarshal([]byte(rawResource), &resource)
	return resource
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

	rsrc := flag.String("src", "./", "The folder or path where the aztfexport files are located")
	ryml := flag.String("conf", "./", "The folder or path where the yaml config file is located")
	check := flag.Bool("validate", false, "Validate the contents of the yaml config against the terraform file")

	flag.Parse()

	src := util.NormalizePath(*rsrc)
	yml := util.NormalizePath(*ryml)

	fmt.Print(util.EmphasizeStr(fmt.Sprintf("Reading config in %s\n", yml), util.Yellow, util.Normal))
	fmt.Print(util.EmphasizeStr(fmt.Sprintf("Reading terraform code in %s\n", src), util.Yellow, util.Normal))
	if *check {
		fmt.Print(util.EmphasizeStr("A validation will be performed before creating output files\n", util.Yellow, util.Normal))
	}

	tfFile := fmt.Sprintf("%s/main.tf", strings.TrimSuffix(src, "/"))
	configFile := fmt.Sprintf("%s/tfmodule.yaml", strings.TrimSuffix(yml, "/"))
	conf, err := os.ReadFile(configFile)

	if err != nil {
		log.Fatalf("ERROR: %s doesn't exist", configFile)
	} else {
		fmt.Printf("Reading modules from %s\n", util.EmphasizeStr(configFile, util.Yellow, util.Normal))
	}

	tfFiles, err := os.ReadDir(strings.TrimSuffix(src, "/"))
	allTf := []byte("")

	// for _, v := range tfFiles {
	for i := len(tfFiles) - 1; i >= 0; i-- {
		v := tfFiles[i]
		if strings.HasSuffix(v.Name(), ".tf") {
			fmt.Println(v.Name())
			currentTfFile, err := os.ReadFile(src + v.Name())
			if err != nil {
				log.Fatal(err)
			}
			allTf = append(allTf, currentTfFile...)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	// tf, err := os.ReadFile(tfFile)

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
	parsedBlocks := readTf(allTf)

	// fmt.Printf("Providers length: %d\n", len(result.providers))
	// fmt.Printf("Providers: %v\n", result.providers)
	// fmt.Printf("Modules length: %d\n", len(parsedBlocks.modules))
	// fmt.Printf("Modules: %v\n", parsedBlocks.modules)
	resourceMap := createVars(parsedBlocks.resources, configModules.Modules)
	tfvarsContent := "// Automatically generated variables\n// Should be changed\n"
	varsContent := "// Automatically generated variables\n// Should be changed"
	for name, resource := range resourceMap {
		fmt.Printf("\n\nResource %s\n", name)
		varBlock := ""
		varsContent += fmt.Sprintf("\n\nvariable \"%s\" { type = list(any) }", name)
		for i, v := range resource {
			fmt.Println(v)
			if i == 0 {
				varBlock = name + " = [\n" + fmt.Sprintf("\t{\n%s\n\t}", strings.ReplaceAll(v, "=", ":"))
			} else {
				varBlock = varBlock + fmt.Sprintf(",\n\t{\n%s\n\t}", strings.ReplaceAll(v, "=", ":"))

			}
		}
		tfvarsContent += varBlock + "\n]\n\n"
	}
	// fmt.Printf("Modules: %v\n", parsedBlocks.resources)

	createFolders(configModules)
	err = createFiles(parsedBlocks, varsContent, tfvarsContent, configModules)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Print(util.EmphasizeStr("Emphasize str\n", util.Blue, util.Normal))
	// fmt.Print(util.EmphasizeStr("Emphasize str\n", util.Blue, util.Bold))
}
