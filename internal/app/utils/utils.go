package utils

import "strings"

const (
	EmptyDate = "0001-01-01 00:00:00 +0000 UTC"
)

func ExtractObjectNameFromUrl(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}
