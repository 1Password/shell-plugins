package sdk

import (
	"net/url"
)

func URL(urlStr string) *url.URL {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return parsed
}
