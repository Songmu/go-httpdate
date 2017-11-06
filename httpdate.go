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
	mostFormatReg  = regexp.MustCompile(`^(\d\d?)` + // 1. day
		`(?:\s+|[-\/])` +
		`(\w+)` + // 2. month
		`(?:\s+|[-\/])` +
		`(\d+)` + // 3. year
		`(?:` +
		`(?:\s+|:)` + // separator before clock
		`(\d\d?):(\d\d)` + // 4. 5. hour:min
		`(?::(\d\d))?` + // 6. optional seconds
		`)?` + // optional clock
		`\s*` +
		`([-+]?\d{2,4})?` + // 7. offset
		`\s*` +
		`(\w+)?` + // 8. ASCII representation of timezone.
		`\s*$`)
	ctimeAndAsctimeReg = regexp.MustCompile(`^(\w{1,3})` + // 1. month
		`\s+` +
		`(\d\d?)` + // 2. day
		`\s+` +
		`(\d\d?):(\d\d)` + // 3,4. hour:min
		`(?::(\d\d))?` + // 5. optional seconds
		`\s+` +
		`(?:([A-Za-z]+)\s+)?` + // 6. optional timezone
		`(\d+)` + // 7. year
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

var fourDigitsReg = regexp.MustCompile(`^([-+])?(\d\d?):?(\d\d)?$`)

func fourDigits2offset(str string) int {
	if matches := fourDigitsReg.FindStringSubmatch(str); len(matches) > 0 {
		hour := a2i(matches[2])
		min := a2i(matches[3])
		offset := hour*60*60 + min*60
		if matches[1] == "-" {
			offset *= -1
		}
		return offset
	}
	return 0
}

// Str2Time detect date format from string and parse it
func Str2Time(timeStr string, loc *time.Location) (time.Time, error) {
	// no time zone is detected and loc is nil, UTC location is used (time.Local is better?)
	if matches := fastReg.FindStringSubmatch(timeStr); len(matches) > 0 {
		d := time.Date(
			a2i(matches[3]),
			shortMonth2Month[matches[2]],
			a2i(matches[1]),
			a2i(matches[4]),
			a2i(matches[5]),
			a2i(matches[6]),
			0,
			time.UTC,
		)
		if loc != nil {
			d = d.In(loc)
		}
		return d, nil
	}
	timeStr = strings.TrimSpace(timeStr)
	timeStr = uselessWdayReg.ReplaceAllString(timeStr, "")

	if matches := mostFormatReg.FindStringSubmatch(timeStr); len(matches) > 0 {
		maybeAMPM := strings.ToLower(matches[8])
		if maybeAMPM != "am" || maybeAMPM != "pm" {
			var l *time.Location
			if matches[8] != "" {
				l2, err := time.LoadLocation(matches[8])
				if err == nil {
					l = l2
				}
			}
			if l == nil && matches[7] != "" {
				l = time.FixedZone(matches[8], fourDigits2offset(matches[7]))
			}
			if l == nil {
				l = loc
				if l == nil {
					l = time.UTC
				}
			}

			y := a2i(matches[3])
			if y < 100 {
				y += 1900
			}
			d := time.Date(
				y,
				shortMonth2Month[matches[2]],
				a2i(matches[1]),
				a2i(matches[4]),
				a2i(matches[5]),
				a2i(matches[6]),
				0,
				l,
			)
			if loc != nil {
				d = d.In(loc)
			}
			return d, nil
		}
	}

	if matches := ctimeAndAsctimeReg.FindStringSubmatch(timeStr); len(matches) > 0 {
		var l *time.Location
		if matches[6] != "" {
			l2, err := time.LoadLocation(matches[6])
			if err == nil {
				l = l2
			}
		}
		if l == nil {
			l = loc
			if l == nil {
				l = time.UTC
			}
		}

		y := a2i(matches[7])
		if y < 100 {
			y += 1900
		}

		d := time.Date(
			y,
			shortMonth2Month[matches[1]],
			a2i(matches[2]),
			a2i(matches[3]),
			a2i(matches[4]),
			a2i(matches[5]),
			0,
			l,
		)
		if loc != nil {
			d = d.In(loc)
		}
		return d, nil
	}

	return time.Time{}, fmt.Errorf("not implemented")
}
