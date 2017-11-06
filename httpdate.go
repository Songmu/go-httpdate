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
		`)?\s*` + // optional clock
		`([-+]?\d{2,4})?\s*` + // 7. offset
		`(\w+)?` + // 8. ASCII representation of timezone.
		`\s*$`)
	ctimeAndAsctimeReg = regexp.MustCompile(`^(\w{1,3})\s+` + // 1. month
		`(\d\d?)\s+` + // 2. day
		`(\d\d?):(\d\d)` + // 3,4. hour:min
		`(?::(\d\d))?\s+` + // 5. optional seconds
		`(?:([A-Za-z]+)\s+)?` + // 6. optional timezone
		`(\d+)` + // 7. year
		`\s*$`)
	unixLsReg = regexp.MustCompile(`^(\w{3})\s+` + // 1. month
		`(\d\d?)\s+` + // 2. day
		`(?:` +
		`(\d{4})|` + // 3. year(optional)
		`(\d{1,2}):(\d{2})` + // 4,5. hour:min
		`(?::(\d{2}))?` + // 6 optional seconds
		`)\s*$`)
	iso8601Reg = regexp.MustCompile(`^(\d{4})` + // 1. year
		`[-\/]?` +
		`(\d\d?)` + // 2. numerical month
		`[-\/]?` +
		`(\d\d?)` + // 3. day
		`(?:` +
		`(?:\s+|[-:Tt])` + // separator before clock
		`(\d\d?):?(\d\d)` + // 4,5. hour:min
		`(?::?(\d\d)(?:\.(\d+))?)?` + // 6,7. optional seconds and fractional
		`)?\s*` + // optional clock
		`([-+]?\d\d?:?(:?\d\d)?|Z|z)?` + // 8. offset (Z is "zero meridian", i.e. GMT)
		`\s*$`)
	winDirReg = regexp.MustCompile(`^(\d{2})-` + // 1. mumerical month
		`(\d{2})-` + // 2. day
		`(\d{2})\s+` + // 3. year
		`(\d\d?):(\d\d)([APap][Mm])` + // 4,5,6. hour:min AM/PM
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
		maybeAMPM := strings.ToUpper(matches[8])
		if maybeAMPM != "AM" && maybeAMPM != "PM" {
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

	if matches := unixLsReg.FindStringSubmatch(timeStr); len(matches) > 0 {
		l := loc
		if l == nil {
			l = time.UTC
		}
		y := a2i(matches[3])
		if matches[3] == "" {
			y = time.Now().Year()
		}
		return time.Date(
			y,
			shortMonth2Month[matches[1]],
			a2i(matches[2]),
			a2i(matches[4]),
			a2i(matches[5]),
			a2i(matches[6]),
			0,
			l,
		), nil
	}

	if matches := iso8601Reg.FindStringSubmatch(timeStr); len(matches) > 0 {
		var l *time.Location
		if strings.ToLower(matches[8]) == "z" {
			l = time.UTC
		} else if matches[8] != "" {
			l = time.FixedZone("", fourDigits2offset(matches[8]))
		} else {
			l = loc
		}
		if l == nil {
			// default timezone
			l = time.UTC
		}
		d := time.Date(
			a2i(matches[1]),
			time.Month(a2i(matches[2])),
			a2i(matches[3]),
			a2i(matches[4]),
			a2i(matches[5]),
			a2i(matches[6]),
			0, // XXX needs care fraction
			l,
		)
		if loc != nil {
			d = d.In(loc)
		}
		return d, nil
	}

	if matches := winDirReg.FindStringSubmatch(timeStr); len(matches) > 0 {
		l := loc
		if l == nil {
			l = time.UTC
		}
		hr := a2i(matches[4])
		switch strings.ToUpper(matches[6]) {
		case "AM":
			if hr == 12 {
				hr = 0
			}
		case "PM":
			if hr != 12 {
				hr += 12
			}
		}
		return time.Date(
			a2i(matches[3])+1900,
			time.Month(a2i(matches[1])),
			a2i(matches[2]),
			hr,
			a2i(matches[5]),
			0,
			0,
			l,
		), nil
	}

	return time.Time{}, fmt.Errorf("not implemented")
}
