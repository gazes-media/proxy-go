package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlQuery := r.URL.Query()["url"]
		if len(urlQuery) > 0 {
			oldUrl := urlQuery[0]
			fmt.Println("url: " + oldUrl)
			defaultHeaders := http.Header{}
			defaultHeaders.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
			defaultHeaders.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
			defaultHeaders.Add("Accept-Language", "en-US,en;q=0.9,fr;q=0.8")
			defaultHeaders.Add("Connection", "keep-alive")
			defaultHeaders.Add("Upgrade-Insecure-Requests", "1")
			defaultHeaders.Add("Cache-Control", "max-age=0")
			defaultHeaders.Add("TE", "Trailers")
			// check if the oldUrl contains a valid website
			if !strings.Contains(oldUrl, "http") || !strings.Contains(oldUrl, "https") {
				returnErr(w)
				// except any possible to the next step
				return
			}
			client := &http.Client{}
			req, err := http.NewRequest("GET", oldUrl, nil)
			if err != nil {
				fmt.Printf("%s", err)
				returnErr(w)
			}
			if strings.Contains(oldUrl, "https://scansmangas.me") {
				defaultHeaders.Add("Authority", "scansmangas.me")
				defaultHeaders.Add("Referer", "https://manga-scan.me/")
			}
			req.Header = defaultHeaders
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("%s", err)
				returnErr(w)
			}
			defer resp.Body.Close()
			for k, v := range resp.Header {
				w.Header().Set(k, v[0])
			}
			w.Header().Del("Content-Encoding")
			w.Header().Del("Server")
			w.Header().Del("Access-Control-Allow-Origin")
			w.Header().Set("Access-Control-Allow-Origin", "*")

			w.WriteHeader(resp.StatusCode)
			bodyToPrint, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("%s", err)
				returnErr(w)
			}

			if strings.Contains(oldUrl, "m3u8") {
				// convert body to string
				bodyString := string(bodyToPrint)
				// split body by new line
				bodySplit := strings.Split(bodyString, "\n")
				bodyToFill := ""
				// loop through each line
				for _, line := range bodySplit {
					// check if line contains http or https
					if strings.Contains(line, "http") || strings.Contains(line, "https") {
						// replace each oldUrl with the proxy oldUrl
						urlToPut := url.QueryEscape(line)
						line = strings.Replace(line, line, "https://proxy.ketsuna.com/?url="+urlToPut, 1)
					}
					// add the line to the body
					bodyToFill += line + "\n"
				}
				// convert body to byte
				bodyToPrint = []byte(bodyToFill)
			}

			w.Write([]byte(bodyToPrint))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Website unavailable"}`))
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "2545"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func returnErr(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"error": "Website unavailable"}`))
	return w
}
