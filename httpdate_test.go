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
