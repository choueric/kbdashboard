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
	"os"

	"github.com/choueric/clog"
)

func main() {
	var cmd string
	pool := topHandlerPool

	clog.SetFlags(clog.Lshortfile | clog.LstdFlags)

	// strip program name
	args := os.Args[1:]
	argc := len(args)

	if argc >= 1 && args[0] == "dump" {
		getConfig(true)
		return
	}

	if argc >= 1 {
		cmd = args[0]
		args = args[1:]
	} else {
		pool.PrintUsage()
		PrintHelpMessage()
		os.Exit(1)
	}

	config := getConfig(false)

	ret := HandleCmd(cmd, pool, args, config)
	os.Exit(ret)
}
