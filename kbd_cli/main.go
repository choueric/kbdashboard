package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/choueric/kernelBuildDashboard/kbd"
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

func buildKernel(item *kbd.Item) {
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
		item.Target,
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

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	pretools := []string{"make"}
	for _, v := range pretools {
		checkPretools(v)
	}

	config, err := kbd.ParseConfig("")
	if err != nil {
		log.Fatal(err)
	}
	if config == nil {
		log.Fatal("config is nil.")
	}
	if len(config.Items) == 0 {
		log.Println("no items in config file.")
	} else {
		for _, item := range config.Items {
			fmt.Println(item, "\n")
		}
	}

	buildKernel(config.Items[0])
}
