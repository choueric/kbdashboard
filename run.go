package main

import (
	"bufio"
	"fmt"
	"log"
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
		installModPath := []string{"INSTALL_MODE_PATH", p.ModInstallDir}
		cmdArgs = append(cmdArgs, strings.Join(installModPath, "="))
	}

	return cmdArgs
}

func makeKernel(p *Profile, target string) {
	cmdArgs := makeKernelOpt(p, target)

	cmd := exec.Command("make", cmdArgs...)
	cmd.Dir = p.SrcDir

	fmt.Printf("    %s%v%s\n", CGREEN, cmdArgs, CEND)

	runCmd(cmd)
}

func configKernel(p *Profile, target string) {
	cmdArgs := makeKernelOpt(p, target)
	args := []string{"make"}
	args = append(args, cmdArgs...)

	os.Chdir(p.SrcDir)
	execCmd("make", args)
}

func runCmd(cmd *exec.Cmd) error {
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Error creating StdoutPipe for Cmd:", err)
		return err
	}

	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		log.Println("create stderrPipe:", err)
		return err
	}

	scanner := bufio.NewScanner(stdoutReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s>>%s %s\n", CGREEN, CEND, scanner.Text())
		}
	}()

	errScanner := bufio.NewScanner(stderrReader)
	go func() {
		for errScanner.Scan() {
			fmt.Printf("%s**%s %s\n", CRED, CEND, errScanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Println("Error starting Cmd:", err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		log.Println("Error waiting for Cmd:", err)
		return err
	}

	return nil
}

func execCmd(name string, argv []string) {
	binary, err := exec.LookPath(name)
	if err != nil {
		log.Fatal(err)
	}

	env := os.Environ()

	err = syscall.Exec(binary, argv, env)
	if err != nil {
		log.Fatal(err)
	}
}
