/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/sergio-prgm/tf-module/pkg/gen"
	"github.com/sergio-prgm/tf-module/pkg/inout"
	"github.com/sergio-prgm/tf-module/pkg/util"
	"github.com/spf13/cobra"
)

// moduleCmd represents the module command
var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "A brief description of your command",
	Long:  `A longer description `,
	Run:   generateYaml,
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}

func generateYaml(cmd *cobra.Command, args []string) {
	yml := util.NormalizePath(ryml)
	src := util.NormalizePath(rsrc)

	resourcesMapping := inout.JsonParser(src + "aztfexportResourceMapping.json")
	csv_resources := inout.ParseCSV(yml + "module_map.csv")
	mapped_yaml := gen.GenerateModuleYaml(resourcesMapping, csv_resources)
	inout.WriteYaml(yml+"tfmodule.yaml", mapped_yaml)

}
