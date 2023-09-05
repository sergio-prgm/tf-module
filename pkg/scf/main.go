// Package scf contains functions whose purpose
// is the general scaffolding of the project
package scf

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sergio-prgm/tf-module/pkg/inout"
	"github.com/sergio-prgm/tf-module/pkg/util"
	"golang.org/x/exp/slices"
)

// createMainFiles creates all the files that are general to the
// terraform project and not iniside of the "modules" folder, i.e.:
// main.tf, tf.vars, variables.tf
func CreateMainFiles(mainContent string, varsContent string, tfvarsContent string) error {
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

func CreateModuleFiles(filePath string, content string, variables string) error {
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

	err = os.WriteFile(filePath+"variables.tf",
		[]byte(variables),
		os.ModePerm)
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
func CreateFiles(parsedBlocks inout.ParsedTf, varsContent string, tfvarsContent string, configModules inout.F) error {
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

	mainContent := strings.Join(parsedBlocks.Providers, "\n\n") + "\n\n" + modulesBlocks
	CreateMainFiles(mainContent, varsContent, tfvarsContent)

	// use in createVars
	for _, v := range configModules.Modules {
		filePath := fmt.Sprintf("./output/Modules/%s/", v.Name)
		variables := ""
		content := ""
		for _, resource := range parsedBlocks.Resources {
			resourceName := strings.Split(resource, "\"")[1]
			if slices.Contains(v.Resources, resourceName) {
				newVar := fmt.Sprintf("variable %s { type = list(any) }\n", strings.Replace(resourceName, "azurerm_", "", 1)+"s")
				if !strings.Contains(variables, newVar) {
					variables += newVar
				}
				if content == "" {
					content = resource
				} else {
					content = content + "\n\n" + resource
				}
			}
		}
		// call createModuleFiles
		CreateModuleFiles(filePath, content, variables)
	}
	return nil
}

// createFolders creates all the necessary folders with the information outlined
// in the yaml config file
func CreateFolders(config inout.F) {
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
