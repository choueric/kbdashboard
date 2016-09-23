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
	j := []string{"-j", strconv.Itoa(p.ThreadNum)}
	output := []string{"O", p.OutputDir}
	cc := []string{"CROSS_COMPILE", p.CrossComile}
	arch := []string{"ARCH", p.Arch}
	installModPath := []string{"INSTALL_MODE_PATH", p.ModInstallDir}

	fmt.Println(p)

	cmdArgs := []string{
		strings.Join(j, ""),
		strings.Join(output, "="),
		strings.Join(cc, "="),
		strings.Join(arch, "="),
		strings.Join(installModPath, "="),
		target,
	}

	return cmdArgs
}

func makeKernel(p *Profile, target string) {
	cmdArgs := makeKernelOpt(p, target)

	cmd := exec.Command("make", cmdArgs...)
	cmd.Dir = p.SrcDir

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
