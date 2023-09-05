package main

import (
	"flag"
	"fmt"

	"log"
	"strings"

	"github.com/sergio-prgm/tf-module/pkg/inout"
	"github.com/sergio-prgm/tf-module/pkg/scf"
	"github.com/sergio-prgm/tf-module/pkg/util"
)

func main() {
	// Maybe extract this towards inout too
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

	configModules := inout.ReadConfig(yml)
	parsedBlocks := inout.ReadTfFiles(src)

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

	scf.CreateFolders(configModules)
	err := scf.CreateFiles(parsedBlocks, varsContent, tfvarsContent, configModules)
	if err != nil {
		log.Fatal(err)
	}
}
