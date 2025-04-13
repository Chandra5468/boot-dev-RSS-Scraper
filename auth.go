package main

import (
	"errors"
	"net/http"
	"strings"
)

// Gets API key from headers of a http request
// Ex : Authorization : Bearer api_key_here
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("no authentication info found in headers")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2 { // splitting with Bearer and key
		return "", errors.New("malformed header for authorization")
	}

	if vals[0] != "Bearer" {
		return "", errors.New("malformed first part of authorization header")
	}
	return vals[1], nil
}
