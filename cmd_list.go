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
	"flag"
	"fmt"
)

var listVerbose bool

func listUsage() {
	printTitle("- list [-v]", false)
	fmt.Printf("  List all profiles.\n")
	fmt.Printf("  -v: Print with more information\n\n")
}

func doList(args []string, config *Config) int {
	fmt.Printf("cmd %s'list'%s:\n", CGREEN, CEND)

	flagSet := flag.NewFlagSet("list", flag.ExitOnError)
	flagSet.BoolVar(&listVerbose, "v", false, "print with more information")
	flagSet.Parse(args)

	for i, p := range config.Profiles {
		if config.Current == i {
			printProfile(&p, listVerbose, true, i)
		} else {
			printProfile(&p, listVerbose, false, i)
		}
	}

	return 0
}

func listHandler(args []string, data interface{}) (int, error) {
	return wrap(doList, args, data)
}
