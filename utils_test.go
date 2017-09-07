package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"testing"
)

func deepCompare(file1, file2 string) bool {
	const chunkSize = 64000

	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func Test_fixRelativeDir(t *testing.T) {
	pre := "/home"
	input := ""
	expect := ""
	output := fixRelativeDir(input, pre)
	if output != expect {
		t.Error("output is:", output)
	}

	input = "./subdir"
	expect = "/home/subdir"
	output = fixRelativeDir(input, pre)
	if output != expect {
		t.Error("output is:", output)
	}

	input = "/boot/subdir"
	expect = "/boot/subdir"
	output = fixRelativeDir(input, pre)
	if output != expect {
		t.Error("output is:", output)
	}
}

func Test_fixHomePath(t *testing.T) {
	input := "~/workspace"
	home := os.Getenv("HOME")
	expect := path.Join(home, "workspace")

	output := fixHomePath(input)
	if output != expect {
		t.Error("wrong output:", output)
	}

	input = "$HOME/workspace"
	output = fixHomePath(input)
	if output != expect {
		t.Error("wrong output:", output)
	}

	input = ""
	expect = ""
	output = fixHomePath(input)
	if output != expect {
		t.Error("wrong output:", output)
	}
}

func Test_isNumber(t *testing.T) {
	input := "12"
	expect := true
	output := isNumber(input)
	if output != expect {
		t.Error("wrong output:", output)
	}

	input = "abc"
	expect = false
	output = isNumber(input)
	if output != expect {
		t.Error("wrong output:", output)
	}

	input = "12abc"
	expect = false
	output = isNumber(input)
	if output != expect {
		t.Error("wrong output:", output)
	}

	input = ""
	expect = false
	output = isNumber(input)
	if output != expect {
		t.Error("wrong output:", output)
	}
}

func Test_checkFileExist(t *testing.T) {
	input := "./utils_test.go"
	expect := true
	output, err := checkFileExsit(input)
	if err != nil {
		t.Error(err)
	} else if output != expect {
		t.Error("output wrong:", output)
	}

	input = "./utils_test.go.test"
	expect = false
	output, err = checkFileExsit(input)
	if err != nil {
		t.Error(err)
	} else if output != expect {
		t.Error("output wrong:", output)
	}

	input = ""
	expect = false
	output, err = checkFileExsit(input)
	if err != nil {
		t.Error(err)
	} else if output != expect {
		t.Error("output wrong:", output)
	}
}

func Test_checkDirExist(t *testing.T) {
	input := "./test_dir"
	if err := checkDirExist(input); err != nil {
		t.Error(err)
	}
	os.Remove(input)
}

func Test_copyFileContents(t *testing.T) {
	src := "Makefile"
	dst := "Makefile.cpy"
	err := copyFileContents(src, dst)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(dst)

	if !deepCompare(src, dst) {
		t.Error("file dose not match")
	}
}

func Test_getCpuNum(t *testing.T) {
	// use 'grep -c ^processors.: /proc/cpuinfo' to get number of CPU
	var result bytes.Buffer
	cmd := exec.Command("grep", "-c", "^processor.:", "/proc/cpuinfo")
	err := pipeCmd(cmd, &result, false)
	if err != nil {
		t.Error(err)
	}
	num := result.String()
	expect, err := strconv.Atoi(num[0 : len(num)-1])
	if err != nil {
		t.Error(err)
	}

	output := getCpuNum()
	if expect != output {
		t.Errorf("expect %d, output %d\n", expect, output)
	}
}
