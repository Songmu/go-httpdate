go-httpdate
=======

[![Build Status](https://travis-ci.org/Songmu/go-httpdate.png?branch=master)][travis]
[![Coverage Status](https://coveralls.io/repos/Songmu/go-httpdate/badge.png?branch=master)][coveralls]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/go-httpdate?status.svg)][godoc]

[travis]: https://travis-ci.org/Songmu/go-httpdate
[coveralls]: https://coveralls.io/r/Songmu/go-httpdate?branch=master
[license]: https://github.com/Songmu/go-httpdate/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/go-httpdate

## Description

Well dealing the date formats used by the HTTP protocol (and then some more).  
`Str2Time` is tremendous function which can detect various date formats from string and returns `time.Time`.

Is is golang porting of perl's [HTTP::Date](https://metacpan.org/pod/HTTP::Date) (ported only `str2time` interface)

## Synopsis

    # try to parse string and returns `time.Time` or error
    t1, err := httpdate.Str2Time("Thu, 03 Feb 1994 12:33:44 GMT", time.UTC)
    t2, err := httpdate.Str2Time("2017-11-11", time.UTC)
    t3, err := httpdate.Str2Time("Thu Nov  9 18:20:31 GMT 2017", time.UTC)
    t4, err := httpdate.Str2Time("08-Feb-94 14:15:29 GMT", time.UTC)

## Supported Format

-  "Wed, 09 Feb 1994 22:23:32 GMT"       -- HTTP format
-  "Thu Feb  3 17:03:55 GMT 1994"        -- ctime(3) format
-  "Thu Feb  3 00:00:00 1994",           -- ANSI C asctime() format
-  "Tuesday, 08-Feb-94 14:15:29 GMT"     -- old rfc850 HTTP format
-  "Tuesday, 08-Feb-1994 14:15:29 GMT"   -- broken rfc850 HTTP format
-  "03/Feb/1994:17:03:55 -0700"   -- common logfile format
-  "09 Feb 1994 22:23:32 GMT"     -- HTTP format (no weekday)
-  "08-Feb-94 14:15:29 GMT"       -- rfc850 format (no weekday)
-  "08-Feb-1994 14:15:29 GMT"     -- broken rfc850 format (no weekday)
-  "1994-02-03 14:15:29 -0100"    -- ISO 8601 format
-  "1994-02-03 14:15:29"          -- zone is optional
-  "1994-02-03"                   -- only date
-  "1994-02-03T14:15:29"          -- Use T as separator
-  "19940203T141529Z"             -- ISO 8601 compact format
-  "19940203"                     -- only date
-  "08-Feb-94"         -- old rfc850 HTTP format    (no weekday, no time)
-  "08-Feb-1994"       -- broken rfc850 HTTP format (no weekday, no time)
-  "09 Feb 1994"       -- proposed new HTTP format  (no weekday, no time)
-  "03/Feb/1994"       -- common logfile format     (no time, no offset)
-  "Feb  3  1994"      -- Unix 'ls -l' format
-  "Feb  3 17:03"      -- Unix 'ls -l' format
-  "11-15-96  03:52PM" -- Windows 'dir' format

## See Also

- [HTTP::Date](https://metacpan.org/pod/HTTP::Date)

## Author

[Songmu](https://github.com/Songmu)
