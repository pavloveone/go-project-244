package code

import (
	"os"
	"strings"
)

func ParseByPaths(paths []string) (string, error) {
	var s strings.Builder
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		s.Write(data)
	}
	return s.String(), nil
}
