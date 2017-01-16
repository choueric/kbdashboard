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

var (
	editProfile string
)

func initEditCmd() {
	cmdmux.HandleFunc("/edit", editConfigHandler)
	cmdmux.HandleFunc("/edit/config", editConfigHandler)

	cmdmux.HandleFunc("/edit/install", editInstallHandler)
	flagSet, err := cmdmux.FlagSet("/edit/install")
	if err != nil {
		clog.Fatal(err)
	}
	flagSet.StringVar(&editProfile, "p", "", "Specify profile by name or index.")
}

func edit_usage() {
	printTitle("- edit [config|install] [profile]")
	fmt.Printf("  Edit various configuraion or scripts using the 'Editor'.\n")
}

////////////////////////////////////////////////////////////////////////////////

func edit_config_usage() {
	printTitle("  - edit config")
	fmt.Printf("    Edit the kbdashboard's configuration file.\n")
	printDefOption("edit")
}

func editConfigHandler(args []string, data interface{}) (int, error) {
	return wrap(edit_config, args, data)
}

func edit_config(args []string, config *Config) int {
	var argv = []string{config.Editor, config.filepath}
	return execCmd(config.Editor, argv)
}

////////////////////////////////////////////////////////////////////////////////

func edit_install_usage() {
	printTitle("  - edit install [profile]")
	fmt.Printf("    Edit [profile]'s installation script.\n")
}

func editInstallHandler(args []string, data interface{}) (int, error) {
	return wrap(edit_install, args, data)
}

func edit_install(args []string, config *Config) int {
	p, _ := getProfile(editProfile, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	var argv = []string{config.Editor, config.getInstallFilename(p)}
	return execCmd(config.Editor, argv)
}
