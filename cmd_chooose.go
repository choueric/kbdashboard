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

func chooseUsage() {
	printTitle("- choose <profile>", false)
	fmt.Printf("  Choose one <profile> as the current one.\n\n")
}

// args[0] is the profile to be choosen
func doChoose(args []string, config *Config) int {
	if len(args) == 0 || args[0] == "" {
		clog.Fatal("Must specify profile's name or index.")
	}
	p, index := doGetProfile(args[0], config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	printCmd("choose", p.Name)
	config.Current = index
	config.save()

	return 0
}

func chooseHandler(args []string, data interface{}) (int, error) {
	return wrap(doChoose, args, data)
}
