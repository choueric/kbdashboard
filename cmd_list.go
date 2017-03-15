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
	"io"
)

func listUsage(w io.Writer, m *helpMap) {
	cmdTitle("list [-v|-c]", false)
	cmdInfo("List all profiles.\n")
	cmdInfo("-v: Print with more information\n")
	cmdInfo("-c: Print full information of current profile.\n\n")
}

func listHandler(args []string, data interface{}) (int, error) {
	fmt.Printf("cmd '%s':\n", cWrap(cGREEN, "list"))

	var verbose, current bool
	flagSet := flag.NewFlagSet("list", flag.ExitOnError)
	flagSet.BoolVar(&verbose, "v", false, "print with more information")
	flagSet.BoolVar(&current, "c", false, "print all information of current profile")
	flagSet.Parse(args)

	if current {
		p := gConfig.Profiles[gConfig.Current]
		printProfile(&p, true, true, gConfig.Current)
		return 0, nil
	}

	for i, p := range gConfig.Profiles {
		printProfile(&p, verbose, gConfig.Current == i, i)
	}

	return 0, nil
}
