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
)

const (
	cRED         = "\x1b[31;1m"
	cGREEN       = "\x1b[32;1m"
	cYELLOW      = "\x1b[33;1m"
	cEND         = "\x1b[0;m"
	level1Indent = "  "
	level2Indent = "      "
	level3Indent = "          "
)

func defMark() string {
	if gConfig.Color {
		return cRED + "*" + cEND
	} else {
		return "*"
	}
}

func cWrap(color string, str string) string {
	if gConfig == nil {
		return str
	}
	if gConfig.Color {
		return color + str + cEND
	} else {
		return str
	}
}

func cmdTitle(format string, def bool, v ...interface{}) {
	if def {
		fmt.Printf("%s %s\n", defMark(), cWrap(cGREEN, fmt.Sprintf(format, v...)))
	} else {
		fmt.Println(level1Indent + cWrap(cGREEN, fmt.Sprintf(format, v...)))
	}
}

func cmdInfo(format string, v ...interface{}) {
	fmt.Printf(level2Indent + fmt.Sprintf(format, v...))
}

func subcmdTitle(format string, def bool, v ...interface{}) {
	if def {
		fmt.Printf("    %s %s\n", defMark(), cWrap(cGREEN, fmt.Sprintf(format, v...)))
	} else {
		fmt.Println(level2Indent + cWrap(cGREEN, fmt.Sprintf(format, v...)))
	}
}

func subcmdInfo(format string, v ...interface{}) {
	fmt.Printf(level3Indent + fmt.Sprintf(format, v...))
}

func printCmd(cmd string, profile string) {
	fmt.Printf("execute command '%s' for [%s]\n", cWrap(cGREEN, cmd),
		cWrap(cGREEN, profile))
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
	}
	return true
}

func checkFileExsit(p string) (bool, error) {
	_, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// check if exist, if not, create one
func checkDirExist(p string) error {
	exist, err := checkFileExsit(p)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	err = os.MkdirAll(p, os.ModeDir|0775)
	if err != nil {
		return err
	}
	return nil
}

func copyFileContents(src, dst string) error {
	fmt.Printf("copy '%s' -> '%s'\n", cWrap(cGREEN, src), cWrap(cGREEN, dst))
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	err = out.Sync()
	return nil
}
