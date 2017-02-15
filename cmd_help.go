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
	"errors"
	"fmt"
	"os"

	"github.com/choueric/cmdmux"
)

type helpMap struct {
	cmd string
	f   func(sub bool)
}

var cmdHelpMap = []helpMap{
	{"list", listUsage},
	{"choose", chooseUsage},
	{"edit", editUsage},
	{"config", configUsage},
	{"build", buildUsage},
	{"install", installUsage},
	{"make", makeUsage},
	{"version", versionUsage},
	{"completion", completionUsage},
	{"help", helpUsage},
}

var COMP_FILENAME = "kbdashboard.bash-completion"

func helpUsage(sub bool) {
	cmdTitle("help [command]", true)
	cmdInfo("Print help message for one or all commands.\n\n")
}

func helpHandler(args []string, data interface{}) (int, error) {
	if len(args) == 0 {
		fmt.Printf("Usage of '%s':\n\n", cWrap(cYELLOW, os.Args[0]))
		for _, v := range cmdHelpMap {
			v.f(false)
		}
		return 1, nil
	}

	cmd := args[0]
	var f func(bool)
	for _, v := range cmdHelpMap {
		if cmd == v.cmd {
			f = v.f
			break
		}
	}
	if f == nil {
		return 0, errors.New(fmt.Sprintf("invalid command '%s'.", cmd))
	}

	fmt.Printf("Usage of comamnd '%s':\n\n", cWrap(cYELLOW, cmd))
	f(true)

	return 1, nil
}

func completionUsage(sub bool) {
	cmdTitle("completion", false)
	cmdInfo("Generate a shell completion file '%s'.\n\n", COMP_FILENAME)
}

func completionHandler(args []string, data interface{}) (int, error) {
	file, err := os.Create(COMP_FILENAME)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	if err = cmdmux.GenerateCompletion("kbdashboard", file); err != nil {
		return 0, err
	}
	fmt.Printf("Create completion file '%s' OK.", COMP_FILENAME)

	return 0, nil
}
