# csp-report-extractor

Extracts CSP report urls from the `report-uri` or `report-url` part in the headers.

Writes results to `csp-found.txt`.

You can provide either a URL or a file of urls.

Install
---
```
go install -v github.com/random-robbie/csp-report-extractor@latest
```



How to run
---

```
csp "https://xxx.xxx.xxx.xxx/#/signin"

https://xxx.xxx.xxx.xxx/vizql/csp-report
```
