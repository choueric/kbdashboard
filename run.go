package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/choueric/kernelBuildDashboard/kbd"
)

func makeKernel(item *kbd.Item, target string) {
	cmdName := "make"
	j := []string{"-j", strconv.Itoa(item.ThreadNum)}
	output := []string{"O", item.OutputDir}
	cc := []string{"CROSS_COMPILE", item.CrossComile}
	arch := []string{"ARCH", item.Arch}
	installModPath := []string{"INSTALL_MODE_PATH", item.ModInstallDir}

	cmdArgs := []string{
		strings.Join(j, ""),
		strings.Join(output, "="),
		strings.Join(cc, "="),
		strings.Join(arch, "="),
		strings.Join(installModPath, "="),
		target,
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = item.SrcDir

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
