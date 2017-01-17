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
}

func helpUsage() {
	printTitle("- help", false)
	fmt.Printf("  Print help message for one or all commands.\n")
}

func helpHandler(args []string, data interface{}) (int, error) {
	if len(args) == 0 {
		fmt.Printf("Usage of %s:\n", os.Args[0])
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

	fmt.Printf("Usage of %s:\n", cmd)
	f()

	return 1, nil
}
