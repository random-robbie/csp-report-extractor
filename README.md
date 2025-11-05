# csp-report-extractor

A tool to extract Content Security Policy (CSP) report URLs from HTTP headers. It identifies reporting endpoints configured in `Content-Security-Policy` and `Content-Security-Policy-Report-Only` headers, specifically the `report-uri` directive.

## Features

- ✅ Extracts `report-uri` directives from CSP headers
- ✅ Detects `report-to` directives (modern CSP reporting)
- ✅ Supports both single URL and bulk file processing
- ✅ Handles relative and absolute URLs
- ✅ Skips SSL/TLS verification for testing environments
- ✅ Configurable output file
- ✅ Verbose mode for debugging
- ✅ Clean, minimal output by default

## Installation

### From source
```bash
go install -v github.com/random-robbie/csp-report-extractor@latest
```

### Build locally
```bash
git clone https://github.com/random-robbie/csp-report-extractor.git
cd csp-report-extractor
go build -o csp-report-extractor
```

## Usage

### Basic usage with a single URL
```bash
csp-report-extractor "https://example.com"
```

### Process multiple URLs from a file
```bash
csp-report-extractor urls.txt
```

### Verbose output with custom output file
```bash
csp-report-extractor -v -o results.txt urls.txt
```

### Command-line options
```
  -v           Enable verbose output
  -o string    Output file for CSP report URLs (default "csp-found.txt")
```

## Examples

### Single URL
```bash
$ csp-report-extractor "https://example.com"
https://example.com/csp-report
```

### File input
Create a file `urls.txt`:
```
https://example.com
https://test.example.org
https://app.example.net
# This is a comment and will be ignored
https://another-site.com
```

Run the tool:
```bash
$ csp-report-extractor urls.txt
https://example.com/csp-report
https://test.example.org/api/csp-violations
https://app.example.net/security/csp
https://another-site.com/report
```

### Verbose mode
```bash
$ csp-report-extractor -v "https://example.com"
Fetching URL: https://example.com
Response status: 200 OK
Found Content-Security-Policy header: default-src 'self'; report-uri /csp-report
https://example.com/csp-report
```

## Output

Results are saved to `csp-found.txt` by default (or the file specified with `-o`). Each discovered CSP report URL is written on a new line.

## Supported CSP Directives

- **report-uri** - Legacy CSP reporting directive (fully supported)
- **report-to** - Modern CSP reporting API (detected but requires Reporting-Endpoints header for full resolution)

## Headers Checked

The tool examines the following HTTP response headers:
- `Content-Security-Policy-Report-Only`
- `Content-Security-Policy`

## Use Cases

- Security assessment and penetration testing
- Bug bounty reconnaissance
- CSP policy auditing
- Monitoring CSP reporting endpoints

## Notes

- The tool skips SSL/TLS certificate verification to work with self-signed certificates in testing environments
- Empty lines and lines starting with `#` in input files are ignored
- Duplicate URLs are automatically removed from results
- Errors are logged but don't stop processing of remaining URLs

## Hosting

You can run this tool on a VPS for continuous monitoring.

[![DigitalOcean Referral Badge](https://web-platforms.sfo2.cdn.digitaloceanspaces.com/WWW/Badge%203.svg)](https://www.digitalocean.com/?refcode=e22bbff5f6f1&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)

Get $200 free credit for 60 days when you sign up with a payment method.

## License

See [LICENSE](LICENSE) file for details.
