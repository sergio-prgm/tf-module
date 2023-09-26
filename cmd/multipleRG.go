package cmd

import (
	"github.com/sergio-prgm/tf-module/pkg/inout"
	"github.com/sergio-prgm/tf-module/pkg/util"
	"github.com/spf13/cobra"
)

// multipleRGCmd represents the multipleRG command
var multipleRGCmd = &cobra.Command{
	Use:   "multipleRG",
	Short: "A brief description of your command",
	Long:  `A longer description`,
	Run:   runMultipleRG,
}

func init() {
	rootCmd.AddCommand(multipleRGCmd)
	multipleRGCmd.PersistentFlags().StringVar(&rsrc_rg, "rg", "./resource_groups", "The folder or path where the resources groups folders are located")
}

func runMultipleRG(cmd *cobra.Command, args []string) {
	util.CheckTerraformVersion()
	src := util.NormalizePath(rsrc_rg)
	dirPath := src + "/___Combined_Resource_Groups___"

	json, terra := inout.ReadMultipleResourceGroups(src, dirPath)

	src = src + "/___Combined_Resource_Groups___/"
	inout.WriteToFile(json, src+"aztfexportResourceMapping.json", "Sucecefully combined the json files")
	inout.WriteToFile(terra, src+"main.tf", "Sucecefully combined the main.tf files")
}
