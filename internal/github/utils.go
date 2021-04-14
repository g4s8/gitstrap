package github

import (
	"net/http"

	gh "github.com/google/go-github/v33/github"
)

// resolveResponseByErr determines the result from a GitHub API response.
// In some cases it is necessary to interpret the "404" response
// as "false" condition e.g. response on get repositoty request.
// Also several GitHub API methods return boolean responses indicated by the HTTP
// status code in the response (true indicated by a 204, false indicated by a
// 404). This helper function will determine that result and hide the 404
// error if present. Any other error will be returned through as-is.
func resolveResponseByErr(err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	if err, ok := err.(*gh.ErrorResponse); ok && err.Response.StatusCode == http.StatusNotFound {
		// Simply false. In this one case, we do not pass the error through.
		return false, nil
	}

	// some other real error occurred
	return false, err
}
