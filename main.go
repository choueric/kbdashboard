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

func handleCmd(cmd string, args []string, config *Config) {
	if cmd == "help" {
		cmd_help(nil, nil)
		return
	}

	h, ok := handlerMap[cmd]
	if !ok {
		clog.Fatalf("[%s] is not supported\n", cmd)
		cmd_help(nil, nil)
	}
	h(args, config)
}

func main() {
	clog.SetFlags(clog.Lshortfile | clog.LstdFlags)

	var (
		args []string
		cmd  string
	)

	switch len(os.Args) {
	case 2:
		cmd = os.Args[1]
	case 3:
		fallthrough
	case 4:
		cmd = os.Args[1]
		args = os.Args[2:]
	default:
		cmd_help(nil, nil)
		return
	}

	config := parseConfig()
	handleCmd(cmd, args, config)
}
