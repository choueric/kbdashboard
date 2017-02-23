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
	"os"
	"text/tabwriter"
)

const VERSION = "0.1.2"
const GIT_DEF = "@"

var (
	BUILD_TIME = "nil"
	GIT_COMMIT = GIT_DEF
)

func versionHandler(args []string, data interface{}) (int, error) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Version\t:", cWrap(cGREEN, VERSION))
	if BUILD_TIME != "nil" {
		fmt.Fprintln(w, "Build time\t:", cWrap(cGREEN, BUILD_TIME))
	}
	if GIT_COMMIT != GIT_DEF {
		fmt.Fprintln(w, "Git Commit\t:", cWrap(cGREEN, GIT_COMMIT))
	}
	w.Flush()
	return 0, nil
}
