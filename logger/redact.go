package logger

import "regexp"

const ipAddrRegex = `\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`

type redactor struct {
	ipAddrRegex *regexp.Regexp
}

func NewRedactor() *redactor {
	return &redactor{regexp.MustCompile(ipAddrRegex)}
}

func (r *redactor) Redact(message string) string {

	if r.ipAddrRegex.MatchString(message) {
		return r.ipAddrRegex.ReplaceAllString(message, "xxx.xxx.xxx.xxx")
	}

	return message
}
