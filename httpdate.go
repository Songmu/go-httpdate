package httpdate

import (
	"fmt"
	"math"
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
		`(?:\s*|[-:Tt])` + // separator before clock
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

var shortMon2Mon = map[string]time.Month{
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

var offsetStrReg = regexp.MustCompile(`^([-+])?(\d\d?):?(\d\d)?$`)

func offsetStr2offset(str string) int {
	if m := offsetStrReg.FindStringSubmatch(str); len(m) > 0 {
		hour := a2i(m[2])
		min := a2i(m[3])
		offset := hour*60*60 + min*60
		if m[1] == "-" {
			offset *= -1
		}
		return offset
	}
	return 0
}

// Str2Time detect date format from string and parse it
func Str2Time(origTimeStr string, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = time.Local
	}
	if m := fastReg.FindStringSubmatch(origTimeStr); len(m) > 0 {
		return time.Date(a2i(m[3]), shortMon2Mon[m[2]], a2i(m[1]), a2i(m[4]), a2i(m[5]), a2i(m[6]), 0, time.UTC), nil
	}
	timeStr := strings.TrimSpace(origTimeStr)
	timeStr = uselessWdayReg.ReplaceAllString(timeStr, "")

	adjustYear := func(str string) int {
		y := a2i(str)
		switch {
		case y >= 100:
			return y
		case y >= 69: // Unix time starts Dec 31 1969 in some time zones
			return y + 1900
		}
		return y + 2000
	}

	loadLoc := func(str string) *time.Location {
		switch strings.ToUpper(str) {
		case "Z":
			return time.UTC
		case "UT":
			str = "UTC"
		}
		l, _ := time.LoadLocation(str)
		return l
	}

	if m := mostFormatReg.FindStringSubmatch(timeStr); len(m) > 0 {
		switch strings.ToUpper(m[8]) {
		case "AM", "PM":
			// nop and through the next check
		default:
			var l *time.Location
			if m[8] != "" {
				l = loadLoc(m[8])
			}
			if l == nil && m[7] != "" {
				l = time.FixedZone(m[8], offsetStr2offset(m[7]))
			}
			if l == nil {
				l = loc
			}
			return time.Date(adjustYear(m[3]), shortMon2Mon[m[2]], a2i(m[1]), a2i(m[4]), a2i(m[5]), a2i(m[6]), 0, l), nil
		}
	}

	if m := ctimeAndAsctimeReg.FindStringSubmatch(timeStr); len(m) > 0 {
		var l *time.Location
		if m[6] != "" {
			l = loadLoc(m[6])
		}
		if l == nil {
			l = loc
		}
		return time.Date(adjustYear(m[7]), shortMon2Mon[m[1]], a2i(m[2]), a2i(m[3]), a2i(m[4]), a2i(m[5]), 0, l), nil
	}

	if m := unixLsReg.FindStringSubmatch(timeStr); len(m) > 0 {
		y := a2i(m[3])
		if m[3] == "" {
			y = time.Now().Year()
		}
		return time.Date(y, shortMon2Mon[m[1]], a2i(m[2]), a2i(m[4]), a2i(m[5]), a2i(m[6]), 0, loc), nil
	}

	if m := iso8601Reg.FindStringSubmatch(timeStr); len(m) > 0 {
		var l *time.Location
		if strings.ToLower(m[8]) == "z" {
			l = time.UTC
		} else if m[8] != "" {
			l = time.FixedZone("", offsetStr2offset(m[8]))
		} else {
			l = loc
		}
		nsec := 0
		fracStr := m[7]
		if fracStr != "" {
			of := 9 - len(fracStr)
			if of <= 0 {
				nsec = a2i(fracStr[0:9])
			} else {
				nsec = a2i(fracStr) * int(math.Pow(10, float64(of)))
			}
		}
		return time.Date(a2i(m[1]), time.Month(a2i(m[2])), a2i(m[3]), a2i(m[4]), a2i(m[5]), a2i(m[6]), nsec, l), nil
	}

	if m := winDirReg.FindStringSubmatch(timeStr); len(m) > 0 {
		l := loc
		hr := a2i(m[4])
		ampm := strings.ToUpper(m[6])
		if ampm == "AM" && hr == 12 {
			hr = 0
		} else if ampm == "PM" && hr != 12 {
			hr += 12
		}
		return time.Date(adjustYear(m[3]), time.Month(a2i(m[1])), a2i(m[2]), hr, a2i(m[5]), 0, 0, l), nil
	}

	return time.Time{}, fmt.Errorf("parsing time %q: parse failed", origTimeStr)
}

// Time2Str returns RFC1123 formatted date
func Time2Str(t time.Time) string {
	return t.In(time.FixedZone("GMT", 0)).Format(time.RFC1123)
}
