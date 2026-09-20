package mutation

import (
	"fmt"
	"net/url"
)

func QueryParam(rawURL, name, value string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse URL: %w", err)
	}

	query := u.Query()

	if _, exists := query[name]; !exists {
		return "", fmt.Errorf("query parameter %q not found", name)
	}

	query.Set(name, value)
	u.RawQuery = query.Encode()

	return u.String(), nil
}
