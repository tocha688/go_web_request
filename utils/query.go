package utils

import "net/url"

func QueryStringify(obj map[string]string) string {
	queryParams := url.Values{}
	for key, value := range obj {
		queryParams.Set(key, value)
	}
	queryString := queryParams.Encode()
	return queryString
}
