# csp-report-extractor

Extracts CSP report urls from the `report-uri` part in the headers.

Writes results to `csp-found.txt`.

You can provide either a URL or a file of urls.

How to run
---

```
go run csp.go "https://xxx.xxx.xxx.xxx/#/signin"

https://xxx.xxx.xxx.xxx/vizql/csp-report
```
