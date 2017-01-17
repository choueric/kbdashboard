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

func makeUsage() {
	printTitle("- make <target>", false)
	fmt.Printf("  Execute 'make' <target> on [profile].\n")
}

func handler_make(args []string, config *Config) int {
	if len(args) <= 1 {
		clog.Error("'make' needs <target> and its parameters.")
		return -1
	}
	target := args[0]
	args = args[1:]

	p, _ := getCurrentProfile(config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	printCmd("build", p.Name)
	return makeKernel(p, target)
}

func makeHandler(args []string, data interface{}) (int, error) {
	return wrap(handler_make, args, data)
}
