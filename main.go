package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

const (
	CRED   = "\x1b[31;1m"
	CGREEN = "\x1b[32;1m"
	CEND   = "\x1b[0;m"
)

func checkPretools(n string) {
	path, err := exec.LookPath(n)
	if err != nil {
		log.Fatal(err, ":", n, " is not found in $PATH")
	}
	log.Printf("'%s' is at %s\n", n, path)
}

func runCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
}

func buildKernel() {
	cmdName := "make"
	cmdArgs := []string{
		"-j4",
		"O=./_build",
		"CROSS_COMPILE=/home/zhs/workspace/TK1/android_dev/prebuilts/gcc/linux-x86/arm/arm-eabi-4.8/bin/arm-eabi-",
		"ARCH=arm",
		"INSTALL_MOD_PATH=./_build/mod",
		"uImage",
	}
	kernelDir := "/home/zhs/workspace/TK1/kernel_android"

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = kernelDir

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

func main() {
	pretools := []string{"make"}

	for _, v := range pretools {
		checkPretools(v)
	}

	buildKernel()
}
