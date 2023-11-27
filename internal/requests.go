package internal

import (
	"fmt"
	"io"
	"net/http"
	"slices"
)

// The function `GetAndHandleGETRequest` sends a GET request to a specified URL and returns the
// response.
func GetAndHandleGETRequest(w http.ResponseWriter, queryURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header = defaultHeaders

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// The function `ForwardResponse` reads the response body from an HTTP response, copies the headers
// from the response to the writer, sets the "Access-Control-Allow-Origin" header to "*", writes the
// response body to the writer, and handles any errors that occur.
func ForwardResponse(w http.ResponseWriter, resp *http.Response) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	for key, values := range resp.Header {
		if slices.Contains([]string{"Access-Control-Allow-Origin", "Server", "Content-Encoding"}, key) {
			continue
		}

		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.StatusCode)

	if _, err := w.Write(body); err != nil {
		fmt.Println("Error writing response:", err)
	}
}
