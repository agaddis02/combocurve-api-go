package pagination

import (
	"net/http"
	"regexp"
)

func GetNextPageURL(responseHeaders http.Header) string {
	linkHeader := responseHeaders.Get("Link")

	re := regexp.MustCompile(`<([^<>]+)>;rel="([^"]+)"`)

	matches := re.FindAllStringSubmatch(linkHeader, -1)

	for _, match := range matches {
		if len(match) == 3 && match[2] == "next" {
			return match[1]
		}
	}

	return ""
}
