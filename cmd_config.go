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

func configUsage() {
	printTitle("- config [menu|def|save]", false)
	fmt.Printf("  Configure kernel of current profile or save it.\n")
	configMenuUsage()
}

////////////////////////////////////////////////////////////////////////////////

func configMenuUsage() {
	printTitle("  - config menu", true)
	fmt.Printf("    Invoke 'make menuconfig' on the current kernel.\n")
}

func doConfigMenu(args []string, config *Config) int {
	p, _ := getCurrentProfile(config)
	printCmd("config menu", p.Name)
	return configKernel(p, "menuconfig")
}

func configMenuHandler(args []string, data interface{}) (int, error) {
	return wrap(doConfigMenu, args, data)
}

////////////////////////////////////////////////////////////////////////////////

func configDefUsage() {
	printTitle("  - config def", false)
	fmt.Printf("    Invoke 'make defconfig' on the current kernel.\n")
}

func doConfigDef(args []string, config *Config) int {
	p, _ := getCurrentProfile(config)
	printCmd("config def", p.Name)
	return makeKernel(p, p.Defconfig)
}

func configDefHandler(args []string, data interface{}) (int, error) {
	return wrap(doConfigDef, args, data)
}

////////////////////////////////////////////////////////////////////////////////

func configSaveUsage() {
	printTitle("  - config save", false)
	fmt.Printf("    Save current config as the default config.\n")
	fmt.Printf("    First execute 'make savedefconfig', then replace the " +
		"config file specified by 'DefConfig'.\n")
}

func doConfigSave(args []string, config *Config) int {
	p, _ := getCurrentProfile(config)

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

func configSaveHandler(args []string, data interface{}) (int, error) {
	return wrap(doConfigSave, args, data)
}
