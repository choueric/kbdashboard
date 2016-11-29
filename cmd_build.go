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

var buildHandlerMap = HandlerMap{
	"image":   {build_image, image_usage},
	"modules": {build_modules, modules_usage},
	"dtb":     {build_dtb, dtb_usage},
}

////////////////////////////////////////////////////////////////////////////////

func build_usage() {
	fmt.Printf("- build [image|modules|dtb] [profile]\n")
	fmt.Printf("  Build various targets.")
	fmt.Printf(" Same as '$ kbdashboard make uImage' if target in config is uImage.\n")
}

func handler_build(args []string, config *Config) int {
	var cmd string

	argc := len(args)
	if argc == 0 {
		cmd = "image"
	} else {
		cmd = args[0]
		args = args[1:]
	}

	return HandleCmd(cmd, buildHandlerMap, args, config)
}

////////////////////////////////////////////////////////////////////////////////

func image_usage() {
	fmt.Printf("build image usage\n")
}

func build_image(args []string, config *Config) int {
	p, _ := getProfile(args, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}
	printCmd("build image", p.Name)
	return makeKernel(p, p.Target)
}

////////////////////////////////////////////////////////////////////////////////

func modules_usage() {
	fmt.Printf("- module [profile]\n")
	fmt.Printf("  Build and install modules.")
	fmt.Printf(" Same as '$ kbdashboard make modules' follwing\n")
	fmt.Printf("  '$ kbdashboard make modules_install'.\n")
}

func build_modules(args []string, config *Config) int {
	p, _ := getProfile(args, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	printCmd("module", p.Name)

	ret := makeKernel(p, "modules")
	if ret != 0 {
		clog.Fatalf("make modules failed.\n")
	}

	return makeKernel(p, "modules_install")
}

////////////////////////////////////////////////////////////////////////////////

func dtb_usage() {
	fmt.Printf("build dtb usage\n")
}

func build_dtb(args []string, config *Config) int {
	return 0
}
