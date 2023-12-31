package util

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
)

// NormalizePath modifies the string given to it to have consistent paths
// across different OSs
func NormalizePath(p string) string {
	path := strings.ReplaceAll(p, "\\", "/")
	if path[len(path)-1] != '/' {
		path += "/"
	}
	return path
}

func StartsWith(line, comp string) (bool, error) {
	if len(line) < len(comp) {
		return false, errors.New("chars out ouf range")
	} else {
		for i := 0; i < len(comp); i++ {
			if comp[i] != line[i] {
				return false, nil
			}
		}
		return true, nil
	}
}

func FirstWordIs(line, comp string) (bool, error) {
	if len(line) < len(comp) {
		return false, errors.New("chars out ouf range")
	} else {
		firstWord := strings.Split(line, " ")[0]
		return firstWord == comp, nil
	}
}

func CheckTerraformVersion() {
	cmdVersion := exec.Command("terraform", "--version")

	// Run the command and capture its output
	output, err := cmdVersion.Output()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
	}

	// Use a regular expression to extract the version number
	re := regexp.MustCompile(`Terraform v([\d\.]+)`)
	matches := re.FindSubmatch(output)
	if matches == nil {
		fmt.Println("Failed to find version in output")
	}

	// Extract version
	versionStr := string(matches[1])

	version, err := semver.NewVersion(versionStr)
	if err != nil {
		fmt.Printf("Failed to parse version: %s\n", err)
	}

	// Check if version is greater than or equal to 1.5
	constraint, _ := semver.NewConstraint(">= 1.5")
	if constraint.Check(version) {
		fmt.Print(EmphasizeStr(fmt.Sprintf("Terraform Version %s is greater than or equal to 1.5\n\n", versionStr), Green, Bold))
	} else {
		fmt.Print(EmphasizeStr(fmt.Sprintf("Terraform Version %s and needs to be at least 1.5\n\n", versionStr), Red, Bold))
	}
}
