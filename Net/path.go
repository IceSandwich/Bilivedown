package Net

import (
	"fmt"
	"strings"
)

func GetDomain(url string) string {
	hs := strings.Index(url, "://")
	he := strings.Index(url[hs+3:], "/")
	if he == -1 {
		he = len(url) - hs - 3
	}
	return url[hs+3 : hs+he+3]
}

func GetProtocol(url string) string {
	hs := strings.Index(url, "://")
	return url[:hs]
}

func CombUrl(protocol string, domain string, url string) string {
	return fmt.Sprintf("%s://%s/%s", protocol, domain, url)
}
