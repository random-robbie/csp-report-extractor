package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
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
	client := &http.Client{Transport: tr, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

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
	reportURI := ""
	for k, v := range resp.Header {
		if strings.EqualFold(k, "Content-Security-Policy-Report-Only") {
			reportURI = v[0]
			break
		}
	}

	// If the header is not present, get the Content-Security-Policy header
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
