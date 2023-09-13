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
	Short: "A brief description",
	Long:  `a longer description`,
	Run:   runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) {
	yml := util.NormalizePath(ryml)
	src := util.NormalizePath(rsrc)

	resourcesMapping := inout.JsonParser(src + "aztfexportResourceMapping.json")
	configModules := inout.ReadConfig(yml)
	csvResources := gen.CheckResources(resourcesMapping, configModules)
	inout.WriteToCsv(csvResources, yml+"existing_resources.csv")
}
