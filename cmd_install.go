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
)

var installOption string

func installHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}

	script := gConfig.getInstallFilename(p)
	ret, err := checkFileExsit(script)
	if err != nil {
		return 0, err
	}

	if ret == false {
		// create and edit script
		fmt.Printf("create install script: '%s'\n", cWrap(cGREEN, script))
		file, err := os.OpenFile(script, os.O_RDWR|os.O_CREATE, 0775)
		if err != nil {
			return 0, err
		}
		defer file.Close()

		str := fmt.Sprintf("#!/bin/sh\n\n# install script for profile '%s'", p.Name)
		if _, err = file.Write([]byte(str)); err != nil {
			return 0, err
		}
		return 0, execCmd(gConfig.Editor, []string{gConfig.Editor, script})
	}

	printCmd("install", p.Name)
	fmt.Printf("    %s\n", cWrap(cGREEN, script))
	cmd := exec.Command(script, args...) // args are the arguments for script.
	cmd.Dir = p.SrcDir
	return 0, pipeCmd(cmd)
}
