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
	"os/exec"
	"regexp"
	"strconv"

	"github.com/choueric/clog"
)

type CmdHandler func(args []string, config *Config)

func checkError(err error) {
	if err != nil {
		clog.Fatal(err)
	}
}

func isNumber(str string) bool {
	if m, _ := regexp.MatchString("^[0-9]+$", str); !m {
		return false
	} else {
		return true
	}
}

func getProfile(arg string, config *Config) (*Profile, int) {
	var p *Profile
	var index int
	if isNumber(arg) {
		n, _ := strconv.Atoi(arg)
		if n >= len(config.Profiles) || n < 0 {
			clog.Fatalf("invalid index of profile: [%d]\n", n)
		}
		p = config.Profiles[n]
		index = n
	} else {
		for i, v := range config.Profiles {
			if v.Name == arg {
				p = v
				index = i
				break
			}
		}
	}

	return p, index
}

func getProfileByCurrent(args []string, config *Config) (*Profile, int) {
	var arg string
	if len(args) == 0 {
		arg = strconv.Itoa(config.Current)
	} else {
		arg = args[0]
	}

	return getProfile(arg, config)
}

////////////////////////////////////////////////////////////////////////////////

var handlerMap = map[string]CmdHandler{
	"list":    cmd_list,
	"choose":  cmd_choose,
	"edit":    cmd_edit,
	"make":    cmd_make,
	"config":  cmd_config,
	"build":   cmd_build,
	"install": cmd_install,
	"module":  cmd_module,
}

func cmd_help(args []string, config *Config) {
	order := []string{
		"list", "choose", "edit", "make", "config",
		"build", "install", "module",
	}
	fmt.Printf("cmd %s'help'%s:\n", CGREEN, CEND)
	fmt.Printf("Usage: \n")
	for _, v := range order {
		fmt.Printf("  - %s\t: ", v)
		handlerMap[v](nil, nil)
	}
	fmt.Println("  - help\t: Display this message.")
}

func listProfile(p *Profile, verbose bool, current bool, i int) {
	if verbose {
		if current {
			fmt.Printf("\n%s*%s ", CRED, CEND)
			fmt.Printf("%s[%d]\t: '%s'%s\n", CGREEN, i, p.Name, CEND)
		} else {
			fmt.Printf("\n  %s[%d]\t: '%s'%s\n", CGREEN, i, p.Name, CEND)
		}
		fmt.Printf("  SrcDir\t\t: %s\n", p.SrcDir)
		fmt.Printf("  Arch\t\t\t: %s\n", p.Arch)
		fmt.Printf("  CC\t\t\t: %s\n", p.CrossComile)
		fmt.Printf("  Target\t\t: %s\n", p.Target)
		fmt.Printf("  BuildDir\t\t: %s\n", p.OutputDir)
		fmt.Printf("  ModInsDir\t\t: %s\n", p.ModInstallDir)
		fmt.Printf("  ThreadNum\t\t: %d\n", p.ThreadNum)
	} else {
		if current {
			fmt.Printf("\n%s*%s ", CRED, CEND)
			fmt.Printf("%s[%d]\t: '%s'%s\n", CGREEN, i, p.Name, CEND)
		} else {
			fmt.Printf("\n  %s[%d]\t: '%s'%s\n", CGREEN, i, p.Name, CEND)
		}
		fmt.Printf("  SrcDir: %s\n", p.SrcDir)
		fmt.Printf("  Arch\t: %s\n", p.Arch)
		fmt.Printf("  CC\t: %s\n", p.CrossComile)
	}
}

func cmd_list(args []string, config *Config) {
	if config == nil {
		fmt.Printf("[-v]. List all profiles. '-v' means verbose.\n")
		return
	}

	verbose := false
	if len(args) == 1 && args[0] == "-v" {
		verbose = true
	}

	fmt.Printf("cmd %s'list'%s:\n", CGREEN, CEND)
	for i, p := range config.Profiles {
		if config.Current == i {
			listProfile(p, verbose, true, i)
		} else {
			listProfile(p, verbose, false, i)
		}
	}
}

func cmd_choose(args []string, config *Config) {
	if config == nil {
		fmt.Printf("Choose current profile.\n")
		return
	}

	if len(args) == 0 {
		clog.Fatal("Choose need profile's name or index")
	}
	p, index := getProfile(args[0], config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	fmt.Printf("cmd %s'choose'%s profile %s[%s]%s\n", CGREEN, CEND, CGREEN, p.Name, CEND)
	config.Current = index

	writeConfigFile(config)
}

func cmd_edit(args []string, config *Config) {
	if config == nil {
		fmt.Printf("Edit the config file using editor specified in config file.\n")
		return
	}

	var argv = []string{config.Editor, config.configFile}
	execCmd(config.Editor, argv)
}

func cmd_make(args []string, config *Config) {
	argc := len(args)
	if config == nil || argc == 0 {
		fmt.Printf("<target> [name | index]. Execute 'make' with specify target.\n")
		return
	}

	target := args[0]
	args = args[1:]

	p, _ := getProfileByCurrent(args, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	fmt.Printf("cmd %s'build %s'%s for %s[%s]%s\n", CGREEN, target, CEND, CGREEN, p.Name, CEND)
	makeKernel(p, target)
}

func cmd_config(args []string, config *Config) {
	if config == nil {
		fmt.Printf("[name | index]. Configure kernel using menuconfig.\n")
		fmt.Printf("\t\t  Same as '$ kbdashboard make menuconfig'.\n")
		return
	}

	p, _ := getProfileByCurrent(args, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	fmt.Printf("cmd %s'config'%s for %s[%s]%s\n", CGREEN, CEND, CGREEN, p.Name, CEND)
	configKernel(p, "menuconfig")
}

func cmd_build(args []string, config *Config) {
	if config == nil {
		fmt.Printf("[name | index]. Build kernel specified by name or index.\n")
		fmt.Printf("\t\t  Same as '$ kbdashboard make uImage' if target in config is uImage.\n")
		return
	}

	p, _ := getProfileByCurrent(args, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	fmt.Printf("cmd %s'build'%s for %s[%s]%s\n", CGREEN, CEND, CGREEN, p.Name, CEND)
	makeKernel(p, p.Target)
}

func cmd_install(args []string, config *Config) {
	if config == nil {
		fmt.Printf("[edit] [name | index]. Execute or edit install script.\n")
		fmt.Printf("\t\t  If use sub-cmd 'edit', open the install script with editor.\n")
		fmt.Printf("\t\t  If no sub-cmd 'edit', execute the install script.\n")
		return
	}

	var doEdit bool
	var create bool

	argc := len(args)
	if argc != 0 && args[0] == "edit" {
		doEdit = true
		args = args[1:]
	}

	p, _ := getProfileByCurrent(args, config)
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
		execCmd(config.Editor, argv)
		return
	}

	fmt.Printf("cmd %s'install'%s profile %s[%s]%s\n", CGREEN, CEND, CGREEN, p.Name, CEND)
	if create {
		// edit script
		var argv = []string{config.Editor, script}
		execCmd(config.Editor, argv)
	} else {
		cmd := exec.Command(script)
		cmd.Dir = p.SrcDir
		fmt.Printf("    %s%s%s\n", CGREEN, script, CEND)
		runCmd(cmd)
	}
}

func cmd_module(args []string, config *Config) {
	if config == nil {
		fmt.Printf("[name | index]. Build and install modules.\n")
		fmt.Printf("\t\t  Same as '$ kbdashboard make modules' follwing\n")
		fmt.Printf("\t\t  '$ kbdashboard make modules_install'.\n")
		return
	}

	p, _ := getProfileByCurrent(args, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	fmt.Printf("cmd %s'module'%s for %s[%s]%s\n", CGREEN, CEND, CGREEN, p.Name, CEND)
	if makeKernel(p, "modules") != nil {
		clog.Fatalf("make modules failed.\n")
	}
	if makeKernel(p, "modules_install") != nil {
		clog.Fatalf("make modules_install faild.\n")
	}
}
