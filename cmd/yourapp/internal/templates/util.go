package templates

import (
	"fmt"
	"github.com/a-h/templ"
	"yourapp/feature/auth"
	"yourapp/foundation/urlprefix"
)

var csrfFormKey = auth.CSRFFormKey

func sprintf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func prefix(url string) templ.SafeURL {
	return templ.SafeURL(urlprefix.Apply(url))
}

func prefixString(url string) string {
	return urlprefix.Apply(url)
}
