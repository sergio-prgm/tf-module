package cmd

import (
	"github.com/sergio-prgm/tf-module/pkg/gen"
	"github.com/sergio-prgm/tf-module/pkg/inout"
	"github.com/sergio-prgm/tf-module/pkg/util"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check which resources are missing mapping",
	Long: `
Compare the tfmodule.yaml file with the aztfexportResourceMapping.json and 
see which resources are missing mapping, generating a CSV with all the resources 
that are missing and with those that exist in the json, and counting them.

It requires the flags --conf with the path for thefolder containing the tfmodule.yaml and the
flag --src with the  path for the folder containing the aztfexportResourceMapping.json`,
	Run: runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.PersistentFlags().BoolVar(&rg, "rg", false, "Cambiar (default false)")
}

func runCheck(cmd *cobra.Command, args []string) {
	util.CheckTerraformVersion()
	yml := util.NormalizePath(ryml)
	src := util.NormalizePath(rsrc)
	if rg {
		src += "___Combined_Resource_Groups___/"
	}

	resourcesMapping := inout.JsonParser(src + "aztfexportResourceMapping.json")
	configModules := inout.ReadConfig(yml)
	csvResources := gen.CheckResources(resourcesMapping, configModules)
	inout.WriteToCsv(csvResources, yml+"existing_resources.csv")
}
