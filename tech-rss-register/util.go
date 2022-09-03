package main

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

func getURLs(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	urls := strings.Split(strings.TrimSpace(string(data)), "\n")
	result := make([]string, 0)
	for _, u := range urls {
		if strings.HasPrefix(u, "#") {
			continue
		}
		result = append(result, strings.TrimRight(u, "/"))
	}
	return result, nil
}
