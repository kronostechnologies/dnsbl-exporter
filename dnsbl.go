package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var (
	resolverUrl = url.URL{
		Scheme:   "https",
		Host:     "cloudflare-dns.com",
		Path:     "dns-query",
		RawQuery: "type=A&do=false&cd=false",
	}
	ipRegex      = regexp.MustCompile(`^([0-9]+)\.([0-9]+)\.([0-9]+)\.([0-9]+)$`)
	reverseIpSub = "$4.$3.$2.$1.%s"
	metrics      []*Metric
	nextCheck    time.Time
)

type stringSet map[string]*struct{}

func (set stringSet) Add(v string) {
	set[v] = nil
}

func (set stringSet) ToList() []string {
	var l []string
	for k, _ := range set {
		l = append(l, k)
	}
	return l
}

type DnsResponseAnswer struct {
	Name string
	Type int
	Data string
}

type DnsResponse struct {
	Status int
	Answer []*DnsResponseAnswer
}

type Metric struct {
	Hostname  string
	IpAddress string
	Value     uint16
	Lists     []string
	ListCount int
}

func getIps(hostname string) []string {
	query := resolverUrl.Query()
	query.Set("name", hostname)
	resolverUrl.RawQuery = query.Encode()

	requestContext, _ := context.WithTimeout(context.Background(), time.Second*1)
	request, requestError := http.NewRequest("GET", resolverUrl.String(), nil)
	if requestError != nil {
		panic(requestError)
	}

	request.Header.Add("Accept", "application/dns-json")
	request.WithContext(requestContext)

	response, responseError := http.DefaultClient.Do(request)
	if responseError != nil {
		panic(responseError)
	}

	out, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		panic(readError)
	}
	closeError := response.Body.Close()
	if closeError != nil {
		panic(closeError)
	}

	dnsResponse := DnsResponse{}
	jsonError := json.Unmarshal(out, &dnsResponse)
	if jsonError != nil {
		panic(jsonError)
	}

	var ipAddresses []string

	if dnsResponse.Status == 0 {
		for _, result := range dnsResponse.Answer {
			if result.Type == 1 {
				ipAddresses = append(ipAddresses, result.Data)
			}
		}
	}

	return ipAddresses
}

func getMetrics(addresses []string) []*Metric {

	if updateCheck(time.Duration(config.Interval) * time.Second) {
		metrics = metrics[:0]

		for _, address := range addresses {
			var ips []string

			if ipRegex.MatchString(address) {
				ips = []string{address}
			} else {
				ips = getIps(address)
			}

			for _, ip := range ips {
				blHostFormat := ipRegex.ReplaceAllString(ip, reverseIpSub)
				blResults := getIps(fmt.Sprintf(blHostFormat, config.Blacklist))

				var ipListCount = 0
				var ipListNames = &stringSet{}

				for _, result := range blResults {
					ipListCount++

					if name, ok := config.ListCodes[result]; ok {
						ipListNames.Add(name)
					} else {
						ipListNames.Add(result)
					}
				}

				metrics = append(metrics, &Metric{
					Hostname:  address,
					IpAddress: ip,
					Lists:     ipListNames.ToList(),
					ListCount: ipListCount,
				})
			}
		}

	}

	return metrics
}

func updateCheck(interval time.Duration) bool {
	now := time.Now()
	if nextCheck.Before(now) {
		nextCheck = now.Add(interval)
		return true
	}
	return false
}
