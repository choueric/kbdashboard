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
	"github.com/choueric/cmdmux"
)

var makeProfile string

func initMakeCmd() {
	cmdmux.HandleFunc("/make", makeHandler)
	if flagSet, err := cmdmux.FlagSet("/make"); err == nil {
		flagSet.StringVar(&makeProfile, "p", "", "Specify profile by name or index.")
	}
}

func make_usage() {
	printTitle("- make <target> [profile]")
	fmt.Printf("  Execute 'make' <target> on [profile].\n")
}

func makeHandler(args []string, data interface{}) (int, error) {
	return wrap(handler_make, args, data)
}

func handler_make(args []string, config *Config) int {
	if len(args) <= 1 {
		clog.Error("need more arguments")
		return -1
	}
	target := args[0]
	args = args[1:]

	p, _ := getProfile(makeProfile, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	printCmd("build", p.Name)
	return makeKernel(p, target)
}
