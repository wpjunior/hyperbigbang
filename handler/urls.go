package handler

import (
	"errors"
	"strings"
)

var errInvalidURL = errors.New("Invalid URL")

func convertPathToURL(path string) (string, error) {
	size := len(path)
	if strings.HasPrefix(path, "/http/") {
		return "http://" + path[6:size], nil
	} else if strings.HasPrefix(path, "/https/") {
		return "https://" + path[7:size], nil
	}
	return "", errInvalidURL
}
