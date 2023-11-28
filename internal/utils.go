package internal

import "net/url"

func isValidURL(urlString string) bool {
	_, err := url.ParseRequestURI(urlString)
	return err == nil
}
