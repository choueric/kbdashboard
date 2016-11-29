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
	"os"
	"os/exec"

	"github.com/choueric/clog"
)

func checkError(err error) {
	if err != nil {
		clog.Fatal(err)
	}
}

func printCmd(cmd string, m string) {
	fmt.Printf("execute command %s'%s'%s for %s[%s]%s\n", CGREEN, cmd, CEND,
		CGREEN, m, CEND)
}

func printDefOption(cmd string) {
	fmt.Printf("    %s*%s This is the default option for %s'%s'%s command.\n",
		CRED, CEND, CGREEN, cmd, CEND)
}

func copyFileContents(src, dst string) (err error) {
	fmt.Printf("copy %s'%s'%s -> %s'%s'%s\n", CGREEN, src, CEND,
		CGREEN, dst, CEND)
	in, err := os.Open(src)
	checkError(err)
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	checkError(err)
	defer out.Close()

	_, err = io.Copy(out, in)
	checkError(err)

	err = out.Sync()
	return
}

////////////////////////////////////////////////////////////////////////////////

func list_usage() {
	printTitle("- list [verbose]")
	fmt.Printf("  List all profiles.\n")
	fmt.Printf("  [verbose]: Print with more information\n")
}

func handler_list(args []string, config *Config) int {
	verbose := false
	if len(args) == 1 && args[0] == "verbose" {
		verbose = true
	}

	fmt.Printf("cmd %s'list'%s:\n", CGREEN, CEND)
	for i, p := range config.Profiles {
		if config.Current == i {
			printProfile(p, verbose, true, i)
		} else {
			printProfile(p, verbose, false, i)
		}
	}

	return 0
}

////////////////////////////////////////////////////////////////////////////////

func choose_usage() {
	printTitle("- choose <profile>")
	fmt.Printf("  Choose <profile> as the current one.\n")
}

func handler_choose(args []string, config *Config) int {
	if len(args) == 0 {
		clog.Fatal("Choose need profile's name or index")
	}
	p, index := doGetProfile(args[0], config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	printCmd("choose", p.Name)
	config.Current = index

	writeConfigFile(config)

	return 0
}

////////////////////////////////////////////////////////////////////////////////

func edit_usage() {
	printTitle("- edit [profile]")
	fmt.Printf("  Edit [profile]'s config using the specified editor.\n")
}

func handler_edit(args []string, config *Config) int {
	var argv = []string{config.Editor, config.configFile}
	return execCmd(config.Editor, argv)
}

////////////////////////////////////////////////////////////////////////////////

func make_usage() {
	printTitle("- make <target> [profile]")
	fmt.Printf("  Execute 'make' <target> on [profile].\n")
}

func handler_make(args []string, config *Config) int {
	target := args[0]
	args = args[1:]

	p, _ := getProfile(args, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	printCmd("build", p.Name)
	return makeKernel(p, target)
}

////////////////////////////////////////////////////////////////////////////////

func install_usage() {
	printTitle("- install [edit] [profile]")
	fmt.Printf("  Execute the install script of [profile].\n")
	fmt.Printf("  [edit]: Open the install script with editor.\n")
}

func handler_install(args []string, config *Config) int {
	var doEdit bool
	var create bool

	argc := len(args)
	if argc != 0 && args[0] == "edit" {
		doEdit = true
		args = args[1:]
	}

	p, _ := getProfile(args, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	script := getInstallFilename(p)
	if checkFileExsit(script) == false {
		// create script
		file, err := os.OpenFile(script, os.O_RDWR|os.O_CREATE, 0775)
		checkError(err)
		_, err = file.Write([]byte("#!/bin/sh"))
		checkError(err)
		file.Close()
		create = true
	}

	if doEdit {
		fmt.Printf("cmd %s'install edit'%s profile %s[%s]%s\n",
			CGREEN, CEND, CGREEN, p.Name, CEND)
		var argv = []string{config.Editor, script}
		return execCmd(config.Editor, argv)
	}

	printCmd("install", p.Name)
	if create {
		// edit script
		var argv = []string{config.Editor, script}
		return execCmd(config.Editor, argv)
	} else {
		cmd := exec.Command(script)
		cmd.Dir = p.SrcDir
		fmt.Printf("    %s%s%s\n", CGREEN, script, CEND)
		return runCmd(cmd)
	}
}

////////////////////////////////////////////////////////////////////////////////

var topHandlerPool = HandlerPool{
	&Handler{"list", handler_list, list_usage},
	&Handler{"choose", handler_choose, choose_usage},
	&Handler{"edit", handler_edit, edit_usage},
	&Handler{"config", handler_config, config_usage},
	&Handler{"build", handler_build, build_usage},
	&Handler{"install", handler_install, install_usage},
	&Handler{"make", handler_make, make_usage},
}
