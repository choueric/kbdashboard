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
	cmdTitle("config [menu|def|save]", false)
	cmdInfo("Handle kernel's configuration.\n")
	configMenuUsage()
	configDefUsage()
	configSaveUsage()
	fmt.Printf("\n")
}

////////////////////////////////////////////////////////////////////////////////

func configMenuUsage() {
	subcmdTitle("config menu", true)
	subcmdInfo("Invoke 'make menuconfig' on the current kernel.\n")
}

func configMenuHandler(args []string, data interface{}) (int, error) {
	p, _ := getCurrentProfile(gConfig)
	printCmd("config menu", p.Name)
	return configKernel(p, "menuconfig"), nil
}

////////////////////////////////////////////////////////////////////////////////

func configDefUsage() {
	subcmdTitle("config def", false)
	subcmdInfo("Invoke 'make defconfig' on the current kernel.\n")
}

func configDefHandler(args []string, data interface{}) (int, error) {
	p, _ := getCurrentProfile(gConfig)
	printCmd("config def", p.Name)
	return makeKernel(p, p.Defconfig), nil
}

////////////////////////////////////////////////////////////////////////////////

func configSaveUsage() {
	subcmdTitle("config save", false)
	subcmdInfo("Save current config as the default config.\n")
	subcmdInfo("Invoke 'make savedefconfig' and then overwrite the original config file.\n")
}

func configSaveHandler(args []string, data interface{}) (int, error) {
	p, _ := getCurrentProfile(gConfig)

	printCmd("config save", p.Name)
	if makeKernel(p, "savedefconfig") != 0 {
		clog.Fatalf("config save failed.\n")
	}

	src := p.OutputDir + "/defconfig"
	dst := p.SrcDir + "/arch/" + p.Arch + "/configs/" + p.Defconfig

	if copyFileContents(src, dst) != nil {
		return 1, nil
	} else {
		return 0, nil
	}
}
