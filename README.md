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
csp-report-extractor "https://xxx.xxx.xxx.xxx/#/signin"

https://xxx.xxx.xxx.xxx/vizql/csp-report
```


You can run it on a VPS.

[![DigitalOcean Referral Badge](https://web-platforms.sfo2.cdn.digitaloceanspaces.com/WWW/Badge%203.svg)](https://www.digitalocean.com/?refcode=e22bbff5f6f1&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)

You get free $200 credit for 60 days if you sign up and add a payment method.
