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

import "fmt"

const VERSION = "0.3"

func versionUsage() {
	cmdTitle("version", false)
	cmdInfo("Print the version.\n\n")
}

func versionHandler(args []string, data interface{}) (int, error) {
	fmt.Printf("Version: %s%s%s\n", CGREEN, VERSION, CEND)
	return 0, nil
}
