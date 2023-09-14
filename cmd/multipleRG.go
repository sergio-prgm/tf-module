package cmd

import (
	"fmt"
	"os"

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

}

func runMultipleRG(cmd *cobra.Command, args []string) {
	src := util.NormalizePath(rsrc)
	dirPath := src + "/___Combined_Resource_Groups___"

	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		// directory exists
		err = os.RemoveAll(dirPath)
		if err != nil {
			fmt.Println("Failed to delete folder:", err)
		} else {
			fmt.Println("Folder deleted successfully.")
		}
	}

	if err := os.Mkdir(dirPath, 0755); err != nil {
		fmt.Println("Error:", err)
	}

	json, terra := inout.ReadMultipleResourceGroups(src, dirPath)

	src = src + "/___Combined_Resource_Groups___/"
	inout.WriteToFile(json, src+"aztfexportResourceMapping.json", "Sucecefully combined the json files")
	inout.WriteToFile(terra, src+"main.tf", "Sucecefully combined the main.tf files")
}
