package utils

import "strings"

func RemoveBetweenChars(s, start, end string) string {
	startIndex := strings.Index(s, start)
	endIndex := strings.LastIndex(s, end)

	if startIndex != -1 && endIndex != -1 && startIndex < endIndex {
		return s[:startIndex] + s[endIndex+1:]
	}

	// Return the original string if the start or end character is not found, or if startIndex >= endIndex
	return s
}
