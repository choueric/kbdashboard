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

	"github.com/choueric/clog"
)

type HandlerFunc func(args []string, config *Config) int

type Handler struct {
	handle HandlerFunc
	usage  func()
}

type HandlerMap map[string]Handler

func (m HandlerMap) PrintUsage() {
	fmt.Printf("cmd %s'help'%s:\n", CGREEN, CEND)
	fmt.Printf("Usage: \n\n")
	for _, v := range m {
		v.usage()
		fmt.Printf("\n")
	}
	fmt.Println("- help\n  Display this message.")
}

func HandleCmd(cmd string, maps HandlerMap, args []string, config *Config) int {
	if cmd == "help" {
		maps.PrintUsage()
		return 0
	}

	h, ok := maps[cmd]
	if !ok {
		maps.PrintUsage()
		clog.Fatalf("[%s] is not supported\n", cmd)
	}
	return h.handle(args, config)
}
