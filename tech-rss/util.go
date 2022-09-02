package main

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

func contains[T comparable](all []*T, target T) bool {
	for _, item := range all {
		if *item == target {
			return true
		}
	}
	return false
}

func getURLs(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	urls := strings.Split(strings.TrimSpace(string(data)), "\n")
	for i, u := range urls {
		urls[i] = strings.TrimRight(u, "/")
	}
	return urls, nil
}
