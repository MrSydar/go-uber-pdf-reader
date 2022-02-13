package invoice

import "regexp"

const (
	dayInDateRegexStr       = `\d{1,2}`
	monthNameInDateRegexStr = `[a-z]+`
	yearInDateRegexStr      = `\d{4}`
)

var dayInDateRegex = regexp.MustCompile(dayInDateRegexStr)
var monthNameInDateRegex = regexp.MustCompile(monthNameInDateRegexStr)
var yearInDateRegex = regexp.MustCompile(yearInDateRegexStr)
