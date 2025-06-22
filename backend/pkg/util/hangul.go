package util

import (
	"strings"
	"unicode"
)

// 초성 목록
var initials = []rune{
	'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ', 'ㅅ',
	'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ',
}

// GetInitials 는 문자열의 한글 음절을 초성 문자열로 변환합니다.
func GetInitials(s string) string {
	var result []rune
	for _, r := range s {
		if r >= 0xAC00 && r <= 0xD7A3 {
			idx := (r - 0xAC00) / 588
			result = append(result, initials[idx])
		} else {
			result = append(result, unicode.ToLower(r))
		}
	}
	return string(result)
}

// HangulMatch 는 문자열이 검색어와 일치하는지 확인합니다.
// 검색어가 한글 초성으로만 이루어진 경우 초성 비교를 수행합니다.
func HangulMatch(target, query string) bool {
	if query == "" {
		return true
	}
	t := strings.ToLower(target)
	q := strings.ToLower(query)

	if strings.Contains(t, q) {
		return true
	}

	initials := GetInitials(t)
	if strings.Contains(initials, q) {
		return true
	}

	return false
}
