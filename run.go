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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func makeKernelOpt(p *Profile, target string) []string {
	cmdArgs := []string{}

	if target != "" {
		cmdArgs = append(cmdArgs, target)
	}

	if p.ThreadNum > 0 {
		j := []string{"-j", strconv.Itoa(p.ThreadNum)}
		cmdArgs = append(cmdArgs, strings.Join(j, ""))
	}

	if p.OutputDir != "" {
		output := []string{"O", p.OutputDir}
		cmdArgs = append(cmdArgs, strings.Join(output, "="))
	}

	if p.CrossComile != "" {
		cc := []string{"CROSS_COMPILE", p.CrossComile}
		cmdArgs = append(cmdArgs, strings.Join(cc, "="))
	}

	if p.Arch != "" {
		arch := []string{"ARCH", p.Arch}
		cmdArgs = append(cmdArgs, strings.Join(arch, "="))
	}

	if p.ModInstallDir != "" {
		installModPath := []string{"INSTALL_MOD_PATH", p.ModInstallDir}
		cmdArgs = append(cmdArgs, strings.Join(installModPath, "="))
	}

	return cmdArgs
}

func makeKernel(p *Profile, target string) error {
	if err := checkDirExist(p.OutputDir); err != nil {
		return err
	}
	if err := checkDirExist(p.ModInstallDir); err != nil {
		return err
	}

	cmdArgs := makeKernelOpt(p, target)

	cmd := exec.Command("make", cmdArgs...)
	cmd.Dir = p.SrcDir

	fmt.Println(cWrap(cGREEN, fmt.Sprintf("    %v", cmdArgs)))

	return pipeCmd(cmd)
}

func configKernel(p *Profile, target string) error {
	if err := checkDirExist(p.OutputDir); err != nil {
		return err
	}
	if err := checkDirExist(p.ModInstallDir); err != nil {
		return err
	}

	cmdArgs := makeKernelOpt(p, target)
	args := []string{"make"}
	args = append(args, cmdArgs...)

	os.Chdir(p.SrcDir)
	return execCmd("make", args)
}

// execute command with Stdout and Stderr being piped in a new process.
func pipeCmd(cmd *exec.Cmd) error {
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdoutReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(cWrap(cGREEN, ">>"), scanner.Text())
		}
	}()

	errScanner := bufio.NewScanner(stderrReader)
	go func() {
		for errScanner.Scan() {
			fmt.Println(cWrap(cRED, "!!"), errScanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

// execute command directly.
func execCmd(name string, argv []string) error {
	binary, err := exec.LookPath(name)
	if err != nil {
		return err
	}

	env := os.Environ()
	err = syscall.Exec(binary, argv, env)
	if err != nil {
		return err
	}

	return nil
}
