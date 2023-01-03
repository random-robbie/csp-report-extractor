package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
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

	resp, err := client.Get(u.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	reportURI := resp.Header.Get("Content-Security-Policy-Report-Only")
	if reportURI == "" {
		reportURI = resp.Header.Get("Content-Security-Policy")
	}

	// Extract the report-uri value from the header
	re := regexp.MustCompile(`report-uri (.*?);`)
	matches := re.FindStringSubmatch(reportURI)
	if len(matches) > 1 {
		reportURL, err := url.Parse(matches[1])
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
