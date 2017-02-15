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

import "errors"

func makeUsage(sub bool) {
	cmdTitle("make <target>", false)
	cmdInfo("Execute '$ make <target>' on current kernel.\n\n")
}

func makeHandler(args []string, data interface{}) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("'make' need paramters.")
	}
	target := args[0]
	args = args[1:]

	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}

	printCmd("build", p.Name)
	return 0, makeKernel(p, target)
}
