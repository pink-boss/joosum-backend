package util

import "strings"

// EnsureHTTPPrefix 는 URL 앞에 http:// 혹은 https:// 가 없으면 https:// 를 붙여 반환합니다.
func EnsureHTTPPrefix(url string) string {
	if url == "" {
		return url
	}
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	return "https://" + url
}
