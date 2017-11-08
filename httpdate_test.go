package httpdate

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestStr2Time(t *testing.T) {
	expect := time.Date(1994, time.February, 3, 0, 0, 0, 0, time.UTC)
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "RFC1123",
			input: "Thu, 03 Feb 1994 00:00:00 GMT",
		},
		{
			name:  "old rfc850 HTTP format",
			input: "Thursday, 03-Feb-94 00:00:00 GMT",
		},
		{
			name:  "broken rfc850 HTTP format",
			input: "Thursday, 03-Feb-1994 00:00:00 GMT",
		},
		{
			name:  "common logfile format1",
			input: "03/Feb/1994:00:00:00 0000",
		},
		{
			name:  "common logfile format2",
			input: "03/Feb/1994:01:00:00 +0100",
		},
		{
			name:  "common logfile format1",
			input: "02/Feb/1994:23:00:00 -0100",
		},
		{
			name:  "HTTP format (no weekday)",
			input: "03 Feb 1994 00:00:00 GMT",
		},
		{
			name:  "old rfc850 (no weekday)",
			input: "03-Feb-94 00:00:00 GMT",
		},
		{
			name:  "broken rfc850 (no weekday)",
			input: "03-Feb-1994 00:00:00 GMT",
		},
		{
			name:  "broken rfc850 (no weekday, no seconds)",
			input: "03-Feb-1994 00:00 GMT",
		},
		{
			name:  "VMS dir listing format",
			input: "03-Feb-1994 00:00",
		},
		{
			name:  "old rfc850 HTTP format (no weekday, no time)",
			input: "03-Feb-94",
		},
		{
			name:  "broken rfc850 HTTP format (no weekday, no time)",
			input: "03-Feb-1994",
		},
		{
			name:  "proposed new HTTP format (no weekday, no time)",
			input: "03 Feb 1994",
		},
		{
			name:  "common logfile format (no time, no offset)",
			input: "03/Feb/1994",
		},
		{
			name:  "A few tests with extra space at various places 1",
			input: "  03/Feb/1994      ",
		},
		{
			name:  "A few tests with extra space at various places 2",
			input: "  03   Feb   1994  0:00  ",
		},
		{
			name:  "Tests a commonly used (faulty?) date format of php cms systems",
			input: "Thu, 03 Feb 1994 00:00:00 +0000 GMT",
		},
		{
			name:  "ctime format",
			input: "Thu Feb  3 00:00:00 GMT 1994",
		},
		{
			name:  "same as ctime, except no TZ",
			input: "Thu Feb  3 00:00:00 1994",
		},
		{
			name:  "Unix 'ls -l' format",
			input: "Feb  3 1994",
		},
		{
			name:  "ISO 8601 formats 1",
			input: "1994-02-03 00:00:00 +0000",
		},
		{
			name:  "ISO 8601 formats 2",
			input: "1994-02-03",
		},
		{
			name:  "ISO 8601 formats 3",
			input: "19940203",
		},
		{
			name:  "ISO 8601 formats 4",
			input: "1994-02-03T00:00:00+0000",
		},
		{
			name:  "ISO 8601 formats 5",
			input: "1994-02-02T23:00:00-0100",
		},
		{
			name:  "ISO 8601 formats 6",
			input: "1994-02-02T23:00:00-01:00",
		},
		{
			name:  "ISO 8601 formats 7",
			input: "1994-02-03T00:00:00 Z",
		},
		{
			name:  "ISO 8601 formats 8",
			input: "19940203T000000Z",
		},
		{
			name:  "ISO 8601 formats 9",
			input: "199402030000",
		},
		{
			name:  "Windows 'dir' format",
			input: "02-03-94  12:00AM",
		},
	}
	for _, tc := range testCases {
		out, err := Str2Time(tc.input, time.UTC)
		if err != nil {
			t.Errorf("%s error should be nil but: %s", tc.name, err)
		}
		out = out.In(time.UTC)
		if !reflect.DeepEqual(out, expect) {
			t.Errorf("Parse failed(%s):\n out:  %+v\n want: %+v", tc.name, out, expect)
		}
	}
}

func TestStr2Time_time(t *testing.T) {
	expect := time.Date(1994, time.February, 3, 4, 56, 1, 0, time.UTC)
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "RFC1123",
			input: "Thu, 03 Feb 1994 04:56:01 GMT",
		},
		{
			name:  "old rfc850 HTTP format",
			input: "Thursday, 03-Feb-94 04:56:01 GMT",
		},
		{
			name:  "broken rfc850 HTTP format",
			input: "Thursday, 03-Feb-1994 04:56:01 GMT",
		},
		{
			name:  "common logfile format1",
			input: "03/Feb/1994:04:56:01 0000",
		},
		{
			name:  "common logfile format2",
			input: "03/Feb/1994:05:56:01 +0100",
		},
		{
			name:  "common logfile format1",
			input: "03/Feb/1994:03:56:01 -0100",
		},
		{
			name:  "HTTP format (no weekday)",
			input: "03 Feb 1994 04:56:01 GMT",
		},
		{
			name:  "old rfc850 (no weekday)",
			input: "03-Feb-94 04:56:01 GMT",
		},
		{
			name:  "broken rfc850 (no weekday)",
			input: "03-Feb-1994 04:56:01 GMT",
		},
		{
			name:  "Tests a commonly used (faulty?) date format of php cms systems",
			input: "Thu, 03 Feb 1994 04:56:01 +0000 GMT",
		},
		{
			name:  "ctime format",
			input: "Thu Feb  3 04:56:01 GMT 1994",
		},
		{
			name:  "same as ctime, except no TZ",
			input: "Thu Feb  3 04:56:01 1994",
		},
		{
			name:  "ISO 8601 formats 1",
			input: "1994-02-03 04:56:01 +0000",
		},
		{
			name:  "ISO 8601 formats 4",
			input: "1994-02-03T04:56:01+0000",
		},
		{
			name:  "ISO 8601 formats 5",
			input: "1994-02-03T03:56:01-0100",
		},
		{
			name:  "ISO 8601 formats 6",
			input: "1994-02-03T03:56:01-01:00",
		},
		{
			name:  "ISO 8601 formats 7",
			input: "1994-02-03T04:56:01 Z",
		},
		{
			name:  "ISO 8601 formats 8",
			input: "19940203T045601Z",
		},
	}
	for _, tc := range testCases {
		out, err := Str2Time(tc.input, time.UTC)
		if err != nil {
			t.Errorf("%s error should be nil but: %s", tc.name, err)
		}
		out = out.In(time.UTC)
		if !reflect.DeepEqual(out, expect) {
			t.Errorf("Parse failed(%s):\n out:  %+v\n want: %+v", tc.name, out, expect)
		}
	}
}

func TestStr2Time_noSec(t *testing.T) {
	y := time.Now().Year()
	expect := time.Date(y, time.February, 3, 16, 56, 0, 0, time.UTC)
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "broken rfc850 (no weekday, no seconds)",
			input: fmt.Sprintf("03-Feb-%d 16:56 GMT", y),
		},
		{
			name:  "VMS dir listing format",
			input: fmt.Sprintf("03-Feb-%d 16:56", y),
		},
		{
			name:  "A few tests with extra space at various places 2",
			input: fmt.Sprintf("  03   Feb   %d  16:56  ", y),
		},
		{
			name:  "Unix 'ls -l' format",
			input: "Feb  3 16:56",
		},
		{
			name:  "ISO 8601 formats 9",
			input: fmt.Sprintf("%d02031656", y),
		},
		{
			name:  "Windows 'dir' format",
			input: fmt.Sprintf("02-03-%d   4:56PM", y-2000),
		},
	}
	for _, tc := range testCases {
		out, err := Str2Time(tc.input, time.UTC)
		if err != nil {
			t.Errorf("%s error should be nil but: %s", tc.name, err)
		}
		out = out.In(time.UTC)
		if !reflect.DeepEqual(out, expect) {
			t.Errorf("Parse failed(%s/%s):\n out:  %+v\n want: %+v", tc.name, tc.input, out, expect)
		}
	}
}

func TestStr2Time_nsec(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  time.Time
	}{
		{
			name:  "Nine 9 nsec",
			input: "2006-01-02T15:04:05.999999999+00:00",
			want:  time.Date(2006, 1, 2, 15, 4, 5, 999999999, time.UTC),
		},
		{
			name:  "Third place after the decimal point",
			input: "2006-01-02T15:04:05.012Z",
			want:  time.Date(2006, 1, 2, 15, 4, 5, 12000000, time.UTC),
		},
	}
	for _, tc := range testCases {
		out, err := Str2Time(tc.input, time.UTC)
		if err != nil {
			t.Errorf("%s error should be nil but: %s", tc.name, err)
		}
		out = out.In(time.UTC)
		if !reflect.DeepEqual(out, tc.want) {
			t.Errorf("Parse failed(%s/%s):\n out:  %+v\n want: %+v", tc.name, tc.input, out, tc.want)
		}
	}
}
