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
	"io"
)

func editUsage(w io.Writer, m *helpMap) {
	defaultHelp(w, m)
	fmt.Printf("\n")
	editProfileUsage()
	editInstallUsage()
	fmt.Printf("\n")
}

////////////////////////////////////////////////////////////////////////////////

func editProfileUsage() {
	subcmdTitle("edit profile", true)
	subcmdInfo("Edit the kbdashboard's configuration file.\n")
}

func editProfileHandler(args []string, data interface{}) (int, error) {
	var argv = []string{gConfig.Editor, gConfig.filepath}
	return 0, execCmd(gConfig.Editor, argv)
}

////////////////////////////////////////////////////////////////////////////////

func editInstallUsage() {
	subcmdTitle("edit install", false)
	subcmdInfo("Edit current profile's installation script.\n")
}

func editInstallHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}
	var argv = []string{gConfig.Editor, gConfig.getInstallFilename(p)}
	return 0, execCmd(gConfig.Editor, argv)
}
