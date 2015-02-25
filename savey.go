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

// parseID parses int64 from provided string.
func parseID(text string) (int64, error) {
	return strconv.ParseInt(regexID.FindString(cleanString(text)), 10, 64)
}

// parseDate parses time.Time from provided string.
func parseDate(date string) (time.Time, error) {
	return time.Parse(dateFormat, date)
}

// cleanString strips leading and ending whitespace.
func cleanString(text string) string {
	return strings.TrimSpace(text)
}
