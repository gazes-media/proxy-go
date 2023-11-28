package internal

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

var client = &http.Client{}
var defaultHeaders = http.Header{
	"User-Agent":                []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"},
	"Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
	"Accept-Language":           []string{"en-US,en;q=0.9,fr;q=0.8"},
	"Connection":                []string{"keep-alive"},
	"Upgrade-Insecure-Requests": []string{"1"},
	"Cache-Control":             []string{"max-age=0"},
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	urlString := r.URL.Query().Get("url")

	if !isValidURL(urlString) {
		http.Error(w, "The URL provided is not valid", http.StatusBadRequest)
		return
	}

	if strings.HasPrefix(urlString, "https://scansmangas.me/") {
		defaultHeaders.Add("Authority", "scansmangas.me")
		defaultHeaders.Add("Referer", "https://manga-scan.me/")
	}

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	resp, err := getClientResponse(urlString, req)
	if err != nil {
		http.Error(w, err.Error(), resp.StatusCode)
		return
	}

	copyResponseHeaders(w, resp.Header)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if strings.Contains(urlString, ".m3u8") {
		proxyM3U8(w, resp.Body)
		return
	}

	copyResponseBody(w, resp.Body)

}

func copyResponseHeaders(w http.ResponseWriter, header http.Header) {
	for key, values := range header {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}
}

func copyResponseBody(w http.ResponseWriter, body io.ReadCloser) {
	defer body.Close()
	_, _ = io.Copy(w, body)
}

func getClientResponse(urlString string, req *http.Request) (*http.Response, error) {
	req.URL, _ = url.Parse(urlString)
	req.Header = defaultHeaders
	return client.Do(req)
}
