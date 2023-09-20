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
	Short: "It terrafys the existing infraestructure in terraform code",
	Long: `
Transfer the infrastructure present in a resource group to separate Terraform code by modules and with all the variables in a .tfvars file.

Requires the --conf flag, which should be the path to the folder with the configuration file (tfmodule.yaml), 
and the --src flag, which should be the path to the folder with the resource group code generated by the aztfexport command`,
	Run: runTerrafy,
}

func init() {
	rootCmd.AddCommand(terrafyCmd)
	terrafyCmd.PersistentFlags().BoolVar(&rg, "rg", false, "Cambiar (default false)")
}

func runTerrafy(cmd *cobra.Command, args []string) {
	src := util.NormalizePath(rsrc)
	yml := util.NormalizePath(ryml)
	if rg {
		src += "___Combined_Resource_Groups___/"
	}

	fmt.Print(util.EmphasizeStr(fmt.Sprintf("Reading config in %s\n", yml), util.Yellow, util.Normal))
	fmt.Print(util.EmphasizeStr(fmt.Sprintf("Reading terraform code in %s\n", src), util.Yellow, util.Normal))

	parsedBlocks := inout.ReadTfFiles(src)
	resourcesMapping := inout.JsonParser(src + "aztfexportResourceMapping.json")

	configModules := inout.ReadConfig(yml)
	scf.CreateFolders(configModules)
	/////
	_, imports_mapping := gen.GenerateImports(resourcesMapping, configModules)

	resourceMap := gen.CreateVars(parsedBlocks.Resources, configModules.Modules, imports_mapping)

	err := scf.CreateFiles(parsedBlocks, resourceMap, configModules)
	if err != nil {
		log.Fatal(err)
	}

}
