package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	verbose    bool
	outputFile string
)

func main() {
	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
	flag.StringVar(&outputFile, "o", "csp-found.txt", "Output file for CSP report URLs")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: csp-report-extractor [options] <URL or file>")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  csp-report-extractor https://example.com")
		fmt.Println("  csp-report-extractor urls.txt")
		fmt.Println("  csp-report-extractor -v -o output.txt urls.txt")
		os.Exit(1)
	}

	input := flag.Arg(0)

	// Check if the argument is a file
	if _, err := os.Stat(input); err == nil {
		processFile(input)
	} else {
		// Process single URL
		processURL(input)
	}
}

func processFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file %s: %v", filename, err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		urlStr := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if urlStr == "" || strings.HasPrefix(urlStr, "#") {
			continue
		}

		if verbose {
			fmt.Printf("[Line %d] Processing: %s\n", lineNum, urlStr)
		}

		processURL(urlStr)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file: %v", err)
	}
}

func processURL(urlStr string) {
	reportURLs, err := extractCSPReports(urlStr)
	if err != nil {
		log.Printf("Error processing %s: %v", urlStr, err)
		return
	}

	if len(reportURLs) == 0 {
		if verbose {
			fmt.Printf("No CSP report URLs found for: %s\n", urlStr)
		}
		return
	}

	for _, reportURL := range reportURLs {
		fmt.Println(reportURL)
		if err := appendToFile(reportURL); err != nil {
			log.Printf("Error writing to file: %v", err)
		}
	}
}

func extractCSPReports(urlStr string) ([]string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Create HTTP client with TLS skip verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   0, // No timeout for now
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	if verbose {
		fmt.Printf("Fetching URL: %s\n", u.String())
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if verbose {
		fmt.Printf("Response status: %s\n", resp.Status)
	}

	var reportURLs []string

	// Check both CSP headers
	headers := []string{
		"Content-Security-Policy-Report-Only",
		"Content-Security-Policy",
	}

	for _, headerName := range headers {
		headerValue := resp.Header.Get(headerName)
		if headerValue == "" {
			continue
		}

		if verbose {
			fmt.Printf("Found %s header: %s\n", headerName, headerValue)
		}

		// Extract report-uri and report-to directives
		urls := parseCSPReportDirectives(headerValue, u)
		reportURLs = append(reportURLs, urls...)
	}

	// Remove duplicates
	return uniqueStrings(reportURLs), nil
}

func parseCSPReportDirectives(headerValue string, baseURL *url.URL) []string {
	var reportURLs []string

	parts := strings.Split(headerValue, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)

		// Handle report-uri directive
		if strings.HasPrefix(part, "report-uri ") {
			uri := strings.TrimSpace(strings.TrimPrefix(part, "report-uri"))
			if uri != "" {
				resolvedURL := resolveURL(uri, baseURL)
				if resolvedURL != "" {
					reportURLs = append(reportURLs, resolvedURL)
				}
			}
		}

		// Handle report-to directive (modern CSP)
		if strings.HasPrefix(part, "report-to ") {
			groupName := strings.TrimSpace(strings.TrimPrefix(part, "report-to"))
			if groupName != "" && verbose {
				fmt.Printf("Found report-to group: %s (requires Reporting-Endpoints header to resolve)\n", groupName)
			}
		}
	}

	return reportURLs
}

func resolveURL(uri string, baseURL *url.URL) string {
	reportURL, err := url.Parse(uri)
	if err != nil {
		if verbose {
			log.Printf("Failed to parse report URI %s: %v", uri, err)
		}
		return ""
	}

	if !reportURL.IsAbs() {
		reportURL = baseURL.ResolveReference(reportURL)
	}

	return reportURL.String()
}

func appendToFile(line string) error {
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, line)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
