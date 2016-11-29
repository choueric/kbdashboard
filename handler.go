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
	cmd    string
	handle HandlerFunc
	usage  func()
}

type HandlerPool []*Handler

func (p HandlerPool) PrintUsage() {
	for _, v := range p {
		v.usage()
		fmt.Printf("\n")
	}
}

func HandleCmd(cmd string, pool HandlerPool, args []string, config *Config) int {
	if cmd == "help" {
		pool.PrintUsage()
		printTitle("- help")
		fmt.Println("  Display this message.")
		return 0
	}

	var h *Handler = nil
	for _, v := range pool {
		if v.cmd == cmd {
			h = v
			break
		}
	}

	if h == nil {
		pool.PrintUsage()
		clog.Fatalf("[%s] is not supported\n", cmd)
	}
	return h.handle(args, config)
}
