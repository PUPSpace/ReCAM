package utils

import (
	"errors"

	"github.com/gosimple/slug"
	"github.com/jaevor/go-nanoid"
)

// GenerateToken func for giving token access to a route.
func GenerateToken() (string, error) {
	randomWords, err := nanoid.CustomASCII("123456789ABCDEFGHIJKLMNPQRSTUVWXYZ", 9)
	if err != nil {
		return "", err
	}

	token := randomWords()
	return token, nil
}

// GenerateSlug func for parsing route name to a slig.
func GenerateSlug(routeName string) (string, error) {
	if routeName == "" {
		return "", errors.New("unable to generate slug from an empty string")
	}

	slugStr := slug.Make(routeName)
	return slugStr, nil
}
