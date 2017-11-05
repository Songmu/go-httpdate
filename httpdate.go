package httpdate

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	// Thu, 03 Feb 1994 00:00:00 GMT
	fastReg = regexp.MustCompile(`^[SMTWF][a-z][a-z], (\d\d) ([JFMAJSOND][a-z][a-z]) (\d\d\d\d) (\d\d):(\d\d):(\d\d) GMT$`)
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
	return time.Time{}, fmt.Errorf("not implemented")
}
