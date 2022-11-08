package swag

import (
	"net/url"
	"strings"
)

// pathJoin joins the given path elements into a single path.
func pathJoin(left, right string) string {
	if left == "" {
		if right == "" {
			return ""
		}

		return "/" + strings.TrimPrefix(right, "/")
	}

	r, err := url.JoinPath(left, right)

	if len(left) > 0 && left[0] != '/' {
		return "/" + r
	}

	_ = err
	return r
}
