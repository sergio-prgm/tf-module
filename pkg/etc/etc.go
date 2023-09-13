// Package etc contains functions that don't match
// any other package (general validation, manipulating config files, etc.)
package etc

import (
	"github.com/sergio-prgm/tf-module/pkg/inout"
)

// validateModules checks whether the information in the config file matches the
// contents of the main.tf and prompts the user the information
func validateModules(configFile inout.F, parsedFile inout.ParsedTf) bool {
	return false
}

// updateConfig allows the user to modify the contents of the config file to
// accommodate late changes or forgotten modules/resources/etc.
func updateConfig() error {
	return nil
}
