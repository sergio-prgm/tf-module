package main

import (
	"flag"
	"fmt"

	"log"

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

	parsedBlocks := inout.ReadTfFiles(src)
	resourcesMapping := inout.JsonParser(src + "aztfexportResourceMapping.json")
	csv_resources := inout.ParseCSV(yml + "module_map.csv")

	mapped_yaml := inout.GenerateModuleYaml(resourcesMapping, csv_resources)
	inout.WriteYaml(yml+"tfmodule.yaml", mapped_yaml)
	configModules := inout.ReadConfig(yml)

	resourceMap := inout.CreateVars(parsedBlocks.Resources, configModules.Modules)

	scf.CreateFolders(configModules)
	err := scf.CreateFiles(parsedBlocks, resourceMap, configModules)
	if err != nil {
		log.Fatal(err)
	}
}
