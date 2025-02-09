package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ScanInt() (int, error) {
	var scanned string
	_, err := fmt.Scanln(&scanned)
	if err != nil {
		return 0, err
	}

	result, err := strconv.Atoi(strings.TrimSpace(scanned))
	if err != nil {
		return 0, err
	}

	return result, nil
}
