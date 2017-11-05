package httpdate

import (
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
	}
	for _, tc := range testCases {
		out, err := Str2Time(tc.input, time.UTC)
		if err != nil {
			t.Errorf("%s error should be nil but: %s", tc.name, err)
		}

		if !reflect.DeepEqual(out, expect) {
			t.Errorf("Parse failed(%s):\n out:  %+v\n want: %+v", tc.name, out, expect)
		}
	}
}
