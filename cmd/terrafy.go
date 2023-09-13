package cmd

import (
	"fmt"
	"log"

	"github.com/sergio-prgm/tf-module/pkg/gen"
	"github.com/sergio-prgm/tf-module/pkg/inout"
	"github.com/sergio-prgm/tf-module/pkg/scf"
	"github.com/sergio-prgm/tf-module/pkg/util"
	"github.com/spf13/cobra"
)

// terrafyCmd represents the terrafy command
var terrafyCmd = &cobra.Command{
	Use:   "apply",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runTerrafy,
}

func init() {
	rootCmd.AddCommand(terrafyCmd)
}

func runTerrafy(cmd *cobra.Command, args []string) {
	src := util.NormalizePath(rsrc)
	yml := util.NormalizePath(ryml)

	fmt.Print(util.EmphasizeStr(fmt.Sprintf("Reading config in %s\n", yml), util.Yellow, util.Normal))
	fmt.Print(util.EmphasizeStr(fmt.Sprintf("Reading terraform code in %s\n", src), util.Yellow, util.Normal))

	parsedBlocks := inout.ReadTfFiles(src)
	resourcesMapping := inout.JsonParser(src + "aztfexportResourceMapping.json")
	/*
		csv_resources := inout.ParseCSV(yml + "module_map.csv")

		mapped_yaml := gen.GenerateModuleYaml(resourcesMapping, csv_resources)

		inout.WriteYaml(yml+"tfmodule.yaml", mapped_yaml)
	*/
	configModules := inout.ReadConfig(yml)

	resourceMap := gen.CreateVars(parsedBlocks.Resources, configModules.Modules)

	scf.CreateFolders(configModules)
	err := scf.CreateFiles(parsedBlocks, resourceMap, configModules)
	if err != nil {
		log.Fatal(err)
	}

	gen.GenerateImports(resourcesMapping, configModules)
}
