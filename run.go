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

func doMakeKernelOpt(t string, n int, b, c, a, m string) []string {
	cmdArgs := []string{}

	if t != "" {
		cmdArgs = append(cmdArgs, t)
	}

	if n > 0 {
		j := []string{"-j", strconv.Itoa(n)}
		cmdArgs = append(cmdArgs, strings.Join(j, ""))
	}

	if b != "" {
		output := []string{"O", b}
		cmdArgs = append(cmdArgs, strings.Join(output, "="))
	}

	if c != "" {
		cc := []string{"CROSS_COMPILE", c}
		cmdArgs = append(cmdArgs, strings.Join(cc, "="))
	}

	if a != "" {
		arch := []string{"ARCH", a}
		cmdArgs = append(cmdArgs, strings.Join(arch, "="))
	}

	if m != "" {
		installModPath := []string{"INSTALL_MOD_PATH", m}
		cmdArgs = append(cmdArgs, strings.Join(installModPath, "="))
	}

	return cmdArgs
}

func makeKernelOpt(p *Profile, target string) []string {
	return doMakeKernelOpt(target, p.ThreadNum, p.BuildDir, p.CrossComile,
		p.Arch, p.ModInstallDir)
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
	if target == "kernelversion" {
		cmd.Dir = p.BuildDir
	} else {
		cmd.Dir = p.SrcDir
	}

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
