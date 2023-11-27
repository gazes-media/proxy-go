package internal

import (
	"io"
	"net/url"
	"regexp"
)

// The function `processM3U8` takes an `io.Reader` containing an M3U8 file, replaces all URLs with a
// modified version, and returns the modified file as a byte slice.
func ProcessM3U8(body io.Reader) ([]byte, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	regex := regexp.MustCompile(`https?://\S+`)
	prefix := []byte("https://proxy.ketsuna.com/?url=")

	result := regex.ReplaceAllFunc(bodyBytes, func(match []byte) []byte {
		encodedURL := url.QueryEscape(string(match))
		return append(prefix, []byte(encodedURL)...)
	})

	return result, nil
}
