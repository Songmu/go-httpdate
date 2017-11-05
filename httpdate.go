package httpdate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// Thu, 03 Feb 1994 00:00:00 GMT
	fastReg        = regexp.MustCompile(`^[SMTWF][a-z][a-z], (\d\d) ([JFMAJSOND][a-z][a-z]) (\d\d\d\d) (\d\d):(\d\d):(\d\d) GMT$`)
	uselessWdayReg = regexp.MustCompile(`^(?i)(?:Sun|Mon|Tue|Wed|Thu|Fri|Sat)[a-z]*,?\s*`)
	mostFormatReg  = regexp.MustCompile(`^(\d\d?)` + // day
		`(?:\s+|[-\/])` +
		`(\w+)` + // month
		`(?:\s+|[-\/])` +
		`(\d+)` + // year
		`(?:` +
		`(?:\s+|:)` + // separator before clock
		`(\d\d?):(\d\d)` + // hour:min
		`(?::(\d\d))?` + // optional seconds
		`)?` + // optional clock
		`\s*` +
		`([-+]?\d{2,4}|[A-Za-z]+)?` + // timezone
		`\s*` +
		`(?:\(\w+\)|\w{3,})?` + // ASCII representation of timezone.
		`\s*$`)
)

var shortMonth2Month = map[string]time.Month{
	"Jan": time.January,
	"Feb": time.February,
	"Mar": time.March,
	"Apr": time.April,
	"May": time.May,
	"Jun": time.June,
	"Jul": time.July,
	"Aug": time.August,
	"Sep": time.September,
	"Oct": time.October,
	"Nov": time.November,
	"Dec": time.December,
}

func a2i(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func Str2Time(timeStr string, loc *time.Location) (time.Time, error) {
	if matches := fastReg.FindStringSubmatch(timeStr); len(matches) > 0 {
		return time.Date(
			a2i(matches[3]),
			shortMonth2Month[matches[2]],
			a2i(matches[1]),
			a2i(matches[4]),
			a2i(matches[5]),
			a2i(matches[6]),
			0,
			loc,
		), nil
	}
	timeStr = strings.TrimSpace(timeStr)
	timeStr = uselessWdayReg.ReplaceAllString(timeStr, "")

	return time.Time{}, fmt.Errorf("not implemented")
}
