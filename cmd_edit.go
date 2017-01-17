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

func editUsage() {
	printTitle("- edit [profile|install]", false)
	fmt.Printf("  Edit profiles or scripts using the 'Editor'.\n")
	editProfileUsage()
	editInstallUsage()
	fmt.Printf("\n")
}

////////////////////////////////////////////////////////////////////////////////

func editProfileUsage() {
	printTitle("  - edit profile", true)
	fmt.Printf("    Edit the kbdashboard's configuration file.\n")
}

func doEditProfile(args []string, config *Config) int {
	var argv = []string{config.Editor, config.filepath}
	return execCmd(config.Editor, argv)
}

func editProfileHandler(args []string, data interface{}) (int, error) {
	return wrap(doEditProfile, args, data)
}

////////////////////////////////////////////////////////////////////////////////

func editInstallUsage() {
	printTitle("  - edit install", false)
	fmt.Printf("    Edit current profile's installation script.\n")
}

func doEditInstall(args []string, config *Config) int {
	p, _ := getCurrentProfile(config)
	var argv = []string{config.Editor, config.getInstallFilename(p)}
	return execCmd(config.Editor, argv)
}

func editInstallHandler(args []string, data interface{}) (int, error) {
	return wrap(doEditInstall, args, data)
}
