dnsbl-exporter
================
DNSBL status Prometheus metric exporter 

Configuration
-------------
### config.yaml
```
blacklist: dnsbl.domain.example
listCodes:
  127.0.0.2: "OPEN PROXY"
  127.0.0.3: "SPAM"
addresses:
  - my.hostname.test
  - 127.0.0.2
```

### Resources
https://tools.ietf.org/html/rfc5782