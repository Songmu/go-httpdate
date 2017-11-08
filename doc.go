// Package httpdate provides functions that deal the date formats used by the
// HTTP protocol (and then some more).  Only the two functions,
// Time2Str() and Str2Time(), are provided.
// Str2Time() is tremendous function which automatically detect various date format and returns `time.Time`.
//
// Is is golang porting of perl's HTTP::Date - https://metacpan.org/pod/HTTP::Date
package httpdate
