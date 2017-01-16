package main

import (
	"fmt"
	"io"
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

func printCmd(cmd string, profile string) {
	fmt.Printf("execute command %s'%s'%s for %s[%s]%s\n", CGREEN, cmd, CEND,
		CGREEN, profile, CEND)
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

func checkError(err error) {
	if err != nil {
		clog.Fatal(err)
	}
}

func printDefOption(cmd string) {
	fmt.Printf("    %s*%s This is the default option for %s'%s'%s command.\n",
		CRED, CEND, CGREEN, cmd, CEND)
}

func copyFileContents(src, dst string) (err error) {
	fmt.Printf("copy %s'%s'%s -> %s'%s'%s\n", CGREEN, src, CEND,
		CGREEN, dst, CEND)
	in, err := os.Open(src)
	checkError(err)
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	checkError(err)
	defer out.Close()

	_, err = io.Copy(out, in)
	checkError(err)

	err = out.Sync()
	return
}
