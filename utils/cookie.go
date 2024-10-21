package utils

import (
	"net/http"
	"strings"

	"github.com/samber/lo"
)

func CookieMergen(cookies ...[]*http.Cookie) []*http.Cookie {
	jar := map[string]*http.Cookie{}
	for _, cookies2 := range cookies {
		for _, co := range cookies2 {
			jar[co.Name] = co
		}
	}
	return lo.Values(jar)
}

func CookieToString(cookies []*http.Cookie) string {
	return strings.Join(lo.Map(cookies, func(it *http.Cookie, index int) string {
		return it.Name + "=" + it.Value
	}), "; ")
}
