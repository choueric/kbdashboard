package main

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/choueric/clog"
)

const (
	CRED    = "\x1b[31;1m"
	CGREEN  = "\x1b[32;1m"
	CEND    = "\x1b[0;m"
	VERSION = "0.1"
)

func printTitle(format string, v ...interface{}) {
	fmt.Printf("%s%s%s\n", CGREEN, fmt.Sprintf(format, v...), CEND)
}

// if OutputDir and ModInstallDir is relative, change it to absolute
// by adding SrcDir prefix.
func fixRelativeDir(p string, pre string) string {
	if !path.IsAbs(p) {
		p = path.Join(pre, p)
	}
	return p
}

func isNumber(str string) bool {
	if m, _ := regexp.MatchString("^[0-9]+$", str); !m {
		return false
	} else {
		return true
	}
}

func checkFileExsit(p string) bool {
	_, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		return false
	} else if err != nil {
		clog.Fatal(err)
	}

	return true
}
