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

func makeKernel(p *Profile, target string) {
	cmdName := "make"
	j := []string{"-j", strconv.Itoa(p.ThreadNum)}
	output := []string{"O", p.OutputDir}
	cc := []string{"CROSS_COMPILE", p.CrossComile}
	arch := []string{"ARCH", p.Arch}
	installModPath := []string{"INSTALL_MODE_PATH", p.ModInstallDir}

	cmdArgs := []string{
		strings.Join(j, ""),
		strings.Join(output, "="),
		strings.Join(cc, "="),
		strings.Join(arch, "="),
		strings.Join(installModPath, "="),
		target,
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = p.SrcDir

	runCmd(cmd)
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

	fmt.Printf("%s %v\n", binary, argv)
	err = syscall.Exec(binary, argv, env)
	if err != nil {
		log.Fatal(err)
	}
}
