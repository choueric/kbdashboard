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

var configHandlerPool = HandlerPool{
	&Handler{"menu", config_menu, menu_usage},
	&Handler{"def", config_def, def_usage},
	&Handler{"save", config_save, save_usage},
}

////////////////////////////////////////////////////////////////////////////////

func config_usage() {
	printTitle("- config [menu|def|save] [profile]")
	fmt.Printf("  Configure [profile] and save it.\n")
	for _, v := range configHandlerPool {
		fmt.Printf("\n")
		v.usage()
	}
}

func handler_config(args []string, config *Config) int {
	var cmd string

	argc := len(args)
	if argc == 0 {
		cmd = "menu"
	} else {
		cmd = args[0]
		args = args[1:]
	}

	return HandleCmd(cmd, configHandlerPool, args, config)
}

////////////////////////////////////////////////////////////////////////////////

func menu_usage() {
	printTitle("  - config menu [profile]")
	fmt.Printf("    Use menuconfig of kernel.\n")
	printDefOption("config")
}

func config_menu(args []string, config *Config) int {
	p, _ := getProfile(args, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	printCmd("config menu", p.Name)
	return configKernel(p, "menuconfig")
}

////////////////////////////////////////////////////////////////////////////////

func def_usage() {
	printTitle("  - config def [profile]")
	fmt.Printf("    Use default config specified in config file.\n")
}

func config_def(args []string, config *Config) int {
	p, _ := getProfile(args, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	printCmd("config def", p.Name)
	return makeKernel(p, p.Defconfig)
}

////////////////////////////////////////////////////////////////////////////////

func save_usage() {
	printTitle("  - config save [profile]")
	fmt.Printf("    Save current config to default config.\n")
	fmt.Printf("    First execute 'make savedefconfig', then substitute the " +
		"config file specified by 'DefConfig'\n")
}

func config_save(args []string, config *Config) int {
	p, _ := getProfile(args, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	printCmd("config save", p.Name)
	if makeKernel(p, "savedefconfig") != 0 {
		clog.Fatalf("config save failed.\n")
	}

	src := p.OutputDir + "/defconfig"
	dst := p.SrcDir + "/arch/" + p.Arch + "/configs/" + p.Defconfig

	if copyFileContents(src, dst) != nil {
		return 1
	} else {
		return 0
	}
}
