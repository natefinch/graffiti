package tags

import "strings"

// commonInitialisms, taken from
// https://github.com/golang/lint/blob/3d26dc39376c307203d3a221bada26816b3073cf/lint.go#L482
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

// TitleToCamel converts a given string to camel case
func TitleToCamel(s string) string {
	var result string
	if initialism := startsWithInitialism(s); initialism != "" {
		result += strings.ToLower(initialism)
	}

	if result == "" {
		return strings.ToLower(s[:1]) + s[1:]
	}

	rest := len(result)
	return result + s[rest:]

}

// startsWithInitialism returns the initialism if the given string begins with it
// taken from
// https://github.com/serenize/snaker/blob/33e5726d116cc1ee16fa506c3e4fed6a553e4cc8/snaker.go#L69
func startsWithInitialism(s string) string {
	var initialism string
	// the longest initialism is 5 char, the shortest 2
	for i := 1; i <= 5; i++ {
		if len(s) > i-1 && commonInitialisms[s[:i]] {
			initialism = s[:i]
		}
	}
	return initialism
}
