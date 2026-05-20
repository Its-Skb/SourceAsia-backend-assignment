package validators

import "testing"

func TestValidateURLs_ValidURLs(t *testing.T) {

	urls := []string{
		"https://example.com/image.jpg",
		"http://example.com/video.mp4",
	}

	err := ValidateURLs(urls)

	if err != nil {
		t.Errorf("expected valid URLs, got error: %v", err)
	}
}

func TestValidateURLs_InvalidScheme(t *testing.T) {

	urls := []string{
		"ftp://example.com/file.jpg",
	}

	err := ValidateURLs(urls)

	if err == nil {
		t.Errorf("expected invalid scheme error")
	}
}

func TestValidateURLs_TooManyURLs(t *testing.T) {

	var urls []string

	for i := 0; i < 21; i++ {
		urls = append(urls, "https://example.com/image.jpg")
	}

	err := ValidateURLs(urls)

	if err == nil {
		t.Errorf("expected max URLs validation error")
	}
}
