package internal

import (
	"fmt"
	"net/http"
	"strings"
)

var defaultHeaders = http.Header{
	"User-Agent":                []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"},
	"Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
	"Accept-Language":           []string{"en-US,en;q=0.9,fr;q=0.8"},
	"Connection":                []string{"keep-alive"},
	"Upgrade-Insecure-Requests": []string{"1"},
	"Cache-Control":             []string{"max-age=0"},
}

// The HandleIndex function handles the index route by making a GET request to a specified URL and
// processing the response based on the URL type.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	queryURL := r.URL.Query().Get("url")

	if strings.Contains(queryURL, "scansmangas.me") {
		defaultHeaders.Add("Authority", "scansmangas.me")
		defaultHeaders.Add("Referer", "https://manga-scan.me/")
	}

	resp, err := GetAndHandleGETRequest(w, queryURL)
	if err != nil {
		errorMessage := fmt.Sprintf("Error occurred while making GET request: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	if strings.Contains(queryURL, "m3u8") {
		modifiedM3U8, err := ProcessM3U8(resp.Body)
		if err != nil {
			http.Error(w, "Failed to parse the m3u8 file", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(modifiedM3U8)
		return
	}

	ForwardResponse(w, resp)
}
