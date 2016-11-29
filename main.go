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

const (
	CRED   = "\x1b[31;1m"
	CGREEN = "\x1b[32;1m"
	CEND   = "\x1b[0;m"
)

func printTitle(format string, v ...interface{}) {
	fmt.Printf("%s%s%s\n", CGREEN, fmt.Sprintf(format, v...), CEND)
}

func parseConfig(dump bool) *Config {
	config, err := ParseConfig("")
	if err != nil {
		clog.Fatal("[ "+checkConfigFile("")+" ]: ", err)
	}
	if config == nil {
		clog.Fatal("config is nil.")
	}

	if dump {
		fmt.Println(config)
	}
	return config
}

func main() {
	var cmd string
	pool := topHandlerPool

	clog.SetFlags(clog.Lshortfile | clog.LstdFlags)

	// strip program name
	args := os.Args[1:]
	argc := len(args)

	if argc >= 1 && args[0] == "dump" {
		parseConfig(true)
		return
	}

	if argc >= 1 {
		cmd = args[0]
		args = args[1:]
	} else {
		pool.PrintUsage()
		printTitle("- help")
		fmt.Println("  Display this message.")
		os.Exit(1)
	}

	config := parseConfig(false)

	ret := HandleCmd(cmd, pool, args, config)
	os.Exit(ret)
}
