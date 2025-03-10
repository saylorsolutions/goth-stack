package urlprefix

import (
	"github.com/saylorsolutions/x/env"
	"net/http"
	"regexp"
	"strings"
)

const EnvUrlPrefix = "URL_PREFIX"

var (
	cachedPrefix   string
	cleansePattern = regexp.MustCompile(`(^[\s/]+|[\s/]+$)`)
)

func init() {
	initFunc()
}

func initFunc() {
	prefix := cleansePattern.ReplaceAllString(env.Val(EnvUrlPrefix, ""), "")
	if len(prefix) == 0 {
		cachedPrefix = ""
		return
	}
	cachedPrefix = "/" + prefix
}

// Apply will apply the configured URL path prefix to the given path.
func Apply(in string) string {
	in = "/" + cleansePattern.ReplaceAllString(in, "")
	if len(cachedPrefix) == 0 {
		return in
	}
	if !strings.HasPrefix(in, cachedPrefix) {
		return cachedPrefix + in
	}
	return in
}

// Get will get the cached URL path prefix.
func Get() string {
	return cachedPrefix
}

// Group applies the configured prefix to the given handler
func Group(next http.Handler) http.Handler {
	if len(cachedPrefix) == 0 {
		return next
	}
	return http.StripPrefix(Get(), next)
}
