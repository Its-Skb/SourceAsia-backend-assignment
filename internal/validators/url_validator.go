package validators

import (
	"errors"
	"net/url"
	"strings"
)

const (
	MaxURLLength      = 2048
	MaxURLsPerRequest = 20
)

func ValidateURLs(urls []string) error {

	// Validate max URLs limit
	if len(urls) > MaxURLsPerRequest {
		return errors.New("maximum 20 URLs allowed per request")
	}

	for _, rawURL := range urls {

		// Trim spaces
		rawURL = strings.TrimSpace(rawURL)

		// Empty URL check
		if rawURL == "" {
			return errors.New("URL cannot be empty")
		}

		// Length validation
		if len(rawURL) > MaxURLLength {
			return errors.New("URL exceeds maximum length of 2048 characters")
		}

		// Parse URL
		parsedURL, err := url.Parse(rawURL)

		if err != nil {
			return errors.New("invalid URL format")
		}

		// Validate scheme
		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			return errors.New("URL must start with http:// or https://")
		}

		// Host validation
		if parsedURL.Host == "" {
			return errors.New("invalid URL host")
		}
	}

	return nil
}
