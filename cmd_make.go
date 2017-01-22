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
	"errors"

	"github.com/choueric/clog"
)

func makeUsage() {
	cmdTitle("make <target>", false)
	cmdInfo("Execute '$ make <target>' on current kernel.\n\n")
}

func makeHandler(args []string, data interface{}) (int, error) {
	if len(args) <= 1 {
		clog.Error("'make' needs <target> and its parameters.")
		return 0, errors.New("need paramters.")
	}
	target := args[0]
	args = args[1:]

	p, _ := getCurrentProfile(gConfig)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	printCmd("build", p.Name)
	return makeKernel(p, target), nil
}
