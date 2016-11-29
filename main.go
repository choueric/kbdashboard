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

const (
	CRED   = "\x1b[31;1m"
	CGREEN = "\x1b[32;1m"
	CEND   = "\x1b[0;m"
)

func parseConfig() *Config {
	config, err := ParseConfig("")
	if err != nil {
		clog.Fatal("[ "+checkConfigFile("")+" ]: ", err)
	}
	if config == nil {
		clog.Fatal("config is nil.")
	}

	//clog.Println(config)
	return config
}

func main() {
	var cmd string
	maps := topHandlerMap

	clog.SetFlags(clog.Lshortfile | clog.LstdFlags)

	// strip program name
	args := os.Args[1:]

	argc := len(args)
	if argc >= 1 {
		cmd = args[0]
		args = args[1:]
	} else {
		maps.PrintUsage()
		os.Exit(1)
	}

	config := parseConfig()

	ret := HandleCmd(cmd, maps, args, config)
	os.Exit(ret)
}
