package httpdate

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"
)

const (
	exitcodeOK = iota
	exitCodeParseFlagErr
	exitCodeErr
)

// Run the cli
func Run(argv []string) int {
	return (&cli{os.Stdout, os.Stderr}).run(argv)
}

type cli struct {
	outStream, errStream io.Writer
}

var numOnlyReg = regexp.MustCompile(`^\d+$`)

func (cl *cli) run(argv []string) int {
	fs := flag.NewFlagSet("httpdate", flag.ContinueOnError)
	fs.SetOutput(cl.errStream)
	fs.Usage = func() {
		fmt.Fprintf(cl.errStream, "Usage of httpdate version %s (rev: %s):\n", Version, revision)
		fs.PrintDefaults()
	}
	var str2time = fs.Bool("s", false, "force str2time")
	if err := fs.Parse(argv); err != nil {
		if err == flag.ErrHelp {
			return exitcodeOK
		}
		return exitCodeParseFlagErr
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return exitCodeParseFlagErr
	}

	arg := fs.Args()[0]
	if !*str2time && !numOnlyReg.MatchString(arg) {
		*str2time = true
	}

	if *str2time {
		t, err := Str2Time(arg, nil)
		if err != nil {
			fmt.Fprintln(cl.errStream, err)
			return exitCodeErr
		}
		fmt.Fprintf(cl.outStream, "%d\n", t.Unix())
	} else {
		s := Time2Str(time.Unix(int64(a2i(arg)), 0))
		fmt.Fprintln(cl.outStream, s)
	}
	return exitcodeOK
}
