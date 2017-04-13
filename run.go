package main

import (
	"bufio"
	"fmt"
	"io"
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

	if p.BuildDir != "" {
		output := []string{"O", p.BuildDir}
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

func makeKernel(p *Profile, target string, w io.Writer, useMarker bool) error {
	if err := checkDirExist(p.BuildDir); err != nil {
		return err
	}
	if err := checkDirExist(p.ModInstallDir); err != nil {
		return err
	}

	cmdArgs := makeKernelOpt(p, target)
	logger.Println(cWrap(cGREEN, fmt.Sprintf("%v", cmdArgs)))

	cmd := exec.Command("make", cmdArgs...)
	cmd.Dir = p.BuildDir

	return pipeCmd(cmd, w, useMarker)
}

func configKernel(p *Profile, target string) error {
	if err := checkDirExist(p.BuildDir); err != nil {
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
// wait until this cmd finishes.
func pipeCmd(cmd *exec.Cmd, w io.Writer, useMarker bool) error {
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdoutMarker := cWrap(cGREEN, ">>")
	stderrMarker := cWrap(cRED, "!!")

	scanner := bufio.NewScanner(stdoutReader)
	go func() {
		for scanner.Scan() {
			if useMarker {
				fmt.Fprintln(w, stdoutMarker, scanner.Text())
			} else {
				fmt.Fprintln(w, scanner.Text())
			}
		}
	}()

	errScanner := bufio.NewScanner(stderrReader)
	go func() {
		for errScanner.Scan() {
			if useMarker {
				fmt.Fprintln(w, stderrMarker, errScanner.Text())
			} else {
				fmt.Fprintln(w, errScanner.Text())
			}
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

	logger.Println("end of pipeCmd.")
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
