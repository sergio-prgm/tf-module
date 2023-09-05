package util

import (
	"errors"
	"strings"
)

// NormalizePath modifies the string given to it to have consistent paths
// across different OSs
func NormalizePath(p string) string {
	return strings.ReplaceAll(p, "\\", "/")
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
