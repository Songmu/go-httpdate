package httpdate_test

import (
	"fmt"
	"time"

	httpdate "github.com/Songmu/go-httpdate"
)

func ExampleStr2Time() {
	t1, _ := httpdate.Str2Time("Thu, 03 Feb 1994 12:33:44 GMT", time.UTC)
	t2, _ := httpdate.Str2Time("2017-11-11", time.UTC)
	t3, _ := httpdate.Str2Time("Thu Nov  9 18:20:31 GMT 2017", time.UTC)
	t4, _ := httpdate.Str2Time("08-Feb-94 14:15:29 GMT", time.UTC)

	fmt.Println(t1)
	fmt.Println(t2)
	fmt.Println(t3)
	fmt.Println(t4)
	// Output:
	// 1994-02-03 12:33:44 +0000 GMT
	// 2017-11-11 00:00:00 +0000 UTC
	// 2017-11-09 18:20:31 +0000 GMT
	// 1994-02-08 14:15:29 +0000 GMT
}
