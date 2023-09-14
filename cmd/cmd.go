package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rsrc string
	ryml string
	rg   bool
)

// In cmd/root.go

func init() {
	rootCmd.PersistentFlags().StringVar(&rsrc, "src", "./src", "The folder or path where the aztfexport files are located")
	rootCmd.PersistentFlags().StringVar(&ryml, "conf", "./conf", "The folder or path where the yaml config file is located")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "tfmodule",
	Short: "A tool to bring your azure infraestructure to terraform code",
}
