package main

import (
	"flag"
	"fmt"

	// "github.com/sergio-prgm/tf-module/utils"
	"log"
	"os"
	"strings"

	"github.com/sergio-prgm/tf-module/pkg/inout"
	"github.com/sergio-prgm/tf-module/pkg/scf"
	"github.com/sergio-prgm/tf-module/pkg/util"
)

type Modules struct {
	Name      string   `yaml:"name"`
	Resources []string `yaml:"resources"`
}

type F struct {
	Modules []Modules `yaml:"modules"`
	Confg   []string  `yaml:"config"`
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

	configModules := inout.ReadConfig(conf)
	// err = yaml.Unmarshal(conf, &configModules)

	// if err != nil {
	// 	log.Fatal()
	// }

	for i := 0; i < len(configModules.Modules); i++ {
		fmt.Printf("\nmodule: %s\nresources: %v\n", configModules.Modules[i].Name, configModules.Modules[i].Resources)
	}
	parsedBlocks := inout.ReadTf(allTf)

	// fmt.Printf("Providers length: %d\n", len(result.providers))
	// fmt.Printf("Providers: %v\n", result.providers)
	// fmt.Printf("Modules length: %d\n", len(parsedBlocks.modules))
	// fmt.Printf("Modules: %v\n", parsedBlocks.modules)
	resourceMap := inout.CreateVars(parsedBlocks.Resources, configModules.Modules)
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

	scf.CreateFolders(configModules)
	err = scf.CreateFiles(parsedBlocks, varsContent, tfvarsContent, configModules)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Print(util.EmphasizeStr("Emphasize str\n", util.Blue, util.Normal))
	// fmt.Print(util.EmphasizeStr("Emphasize str\n", util.Blue, util.Bold))
}
