/*
 * Copyright (C) 2016 Eric Chou <zhssmail@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
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
	CRED   = "\x1b[31;1m"
	CGREEN = "\x1b[32;1m"
	CEND   = "\x1b[0;m"
)

func printTitle(format string, def bool, v ...interface{}) {
	if def {
		fmt.Printf("%s%s%s %s*%s\n", CGREEN, fmt.Sprintf(format, v...), CEND,
			CRED, CEND)
	} else {
		fmt.Printf("%s%s%s\n", CGREEN, fmt.Sprintf(format, v...), CEND)
	}
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

func wrap(f func([]string, *Config) int, args []string, data interface{}) (int, error) {
	config := data.(*Config)
	ret := f(args, config)
	return ret, nil
}
