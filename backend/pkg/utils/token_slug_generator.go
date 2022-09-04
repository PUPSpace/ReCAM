package utils

import (
	// "log"
	"errors"

	"github.com/gosimple/slug"
	"github.com/jaevor/go-nanoid"
)

//GenerateToken func for giving token access to a route.
func GenerateToken() (string, error) {
	randomWords, err := nanoid.CustomASCII("123456789ABCDEFGHIJKLMNPQRSTUVWXYZ", 9)
	if err != nil {
		return "", err
	}

	token := randomWords()
	return token, nil
}

//GenerateSlug func for parsing route name to a slig.
func GenerateSlug(routeName string) (string, error) {
	if len(routeName) == 0 {
		return "", errors.New("unable to generate slug from an empty string")
	}

	slug := slug.Make(routeName)
	return slug, nil
}
