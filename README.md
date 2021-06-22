dnsbl-exporter
================
DNSBL status Prometheus metric exporter 

Configuration
-------------
### config.yaml
```
# blacklist domain name, will be queried with 4.3.2.1.dnsbl.domain.example
blacklist: dnsbl.domain.example

# IP address results to monitor, value is a label for Prometheus
listCodes:
  127.0.0.2: "OPEN PROXY"
  127.0.0.3: "SPAM"
  
# IPs to validate (will resolve A record IPs for hostnames)
addresses:
  - my.hostname.test
  - 127.0.0.2
```


### Resources
https://tools.ietf.org/html/rfc5782
