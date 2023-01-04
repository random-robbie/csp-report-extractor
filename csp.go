package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a URL or file as an argument")
		return
	}

	input := os.Args[1]

	// Check if the argument is a file
	if _, err := os.Stat(input); err == nil {
		// Open the file
		f, err := os.Open(input)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		// Read the file line by line
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			// Set the URL variable as the scanned line
			url := scanner.Text()

			// Send the URL to the grabber function
			grabber(url)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	} else {
		// Set the URL variable as the input argument
		url := input

		// Send the URL to the grabber function
		grabber(url)
	}
}

func grabber(url2 string) {

	u, err := url.Parse(url2)
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
		if reportURI != "" {
			// Open the file in append mode
			f, err := os.OpenFile("csp-found.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()

			// Write the string to the file with a new line after it
			_, err = fmt.Fprintln(f, reportURL)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("Successfully wrote string to file")
		} else {
			fmt.Println("reportURI string is empty")
		}
	}
}
