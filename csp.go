package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a URL as an argument")
		os.Exit(1)
	}

	u, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Get the Content-Security-Policy-Report-Only header, if present
	headerValue := resp.Header.Get("Content-Security-Policy-Report-Only")

	// If the header is not present, get the Content-Security-Policy header
	if headerValue == "" {
		headerValue = resp.Header.Get("Content-Security-Policy")
	}

	var reportURI string

	parts := strings.Split(headerValue, ";")
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), " ", 2)
		if len(kv) != 2 {
			continue
		}

		if kv[0] == "report-uri" {
			reportURI = kv[1]
		}
	}

	if reportURI != "" {
		reportURL, err := url.Parse(reportURI)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if !reportURL.IsAbs() {
			reportURL = u.ResolveReference(reportURL)
		}

		fmt.Println(reportURL)
	}
}
