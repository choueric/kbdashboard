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
	"os"

	"github.com/choueric/clog"
	"github.com/choueric/cmdmux"
)

type helpMap struct {
	cmd string
	f   func()
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

const completionFileName = "kbdashboard-completion"

func helpUsage() {
	cmdTitle("help [command]", false)
	cmdInfo("Print help message for one or all commands.\n\n")
}

func helpHandler(args []string, data interface{}) (int, error) {
	if len(args) == 0 {
		fmt.Printf("Usage of %s'%s'%s:\n\n", CYELLOW, os.Args[0], CEND)
		for _, v := range cmdHelpMap {
			v.f()
		}
		return 1, nil
	}

	cmd := args[0]
	var f func()
	for _, v := range cmdHelpMap {
		if cmd == v.cmd {
			f = v.f
			break
		}
	}
	if f == nil {
		clog.Fatalf("invalid command name '%s'\n", cmd)
	}

	fmt.Printf("Usage of comamnd %s'%s'%s:\n\n", CYELLOW, cmd, CEND)
	f()

	return 1, nil
}

func completionUsage() {
	cmdTitle("completion", false)
	cmdInfo("Generate a shell completion file '%s'.\n\n", completionFileName)
}

func completionHandler(args []string, data interface{}) (int, error) {
	file, err := os.Create(completionFileName)
	if err != nil {
		clog.Fatal(err)
	}
	defer file.Close()

	if err = cmdmux.GenerateCompletion("kbdashboard", file); err != nil {
		clog.Fatal(err)
	}
	clog.Printf("Create completion file '%s' OK.", completionFileName)

	return 0, nil
}
