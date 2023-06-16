package utils

import "errors"


func StartsWith(line, comp string) (bool, error) {
	if len(line) < len(comp) {
		return false, errors.New("Chars out ouf range!")
	} else {
		for i := 0; i < len(comp); i++ {
			if comp[i] != line[i] {
				return false, nil
			}
		}
		return true, nil
	}
}
