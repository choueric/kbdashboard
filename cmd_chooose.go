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

import "github.com/choueric/clog"

func chooseUsage() {
	cmdTitle("choose <profile>", false)
	cmdInfo("Choose <profile> as the current profile.\n\n")
}

// args[0] is the profile to be choosen
func chooseHandler(args []string, data interface{}) (int, error) {
	if len(args) == 0 || args[0] == "" {
		clog.Fatal("Must specify profile's name or index.")
	}
	p, index := doGetProfile(args[0], gConfig)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	printCmd("choose", p.Name)
	gConfig.Current = index
	gConfig.save()

	return 0, nil
}
