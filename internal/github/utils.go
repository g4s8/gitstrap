package github

import (
	"errors"
	"net/http"

	gh "github.com/google/go-github/v38/github"
)

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	ghErr := new(gh.ErrorResponse)
	if errors.As(err, &ghErr) && ghErr.Response.StatusCode == http.StatusNotFound {
		return true
	}
	return false
}
