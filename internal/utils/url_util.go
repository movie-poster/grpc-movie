package utils

import (
	"net/url"
	"strings"
)

func IsURL(content string) bool {
	_, err := url.ParseRequestURI(content)
	return err == nil && strings.HasPrefix(content, "https://")
}
