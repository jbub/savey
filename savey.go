package savey

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "January 2, 2006"
)

var (
	regexID = regexp.MustCompile("\\d+")
)

// ParseID parses int64 from provided string.
func ParseID(text string) (int64, error) {
	return strconv.ParseInt(regexID.FindString(CleanString(text)), 10, 64)
}

// ParseDate parses time.Time from provided string.
func ParseDate(date string) (time.Time, error) {
	return time.Parse(dateFormat, date)
}

// CleanString strips leading and ending whitespace.
func CleanString(text string) string {
	return strings.TrimSpace(text)
}
