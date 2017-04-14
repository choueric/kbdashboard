package main

import (
	"os"
	"path"
	"testing"
)

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
}

func Test_checkDirExist(t *testing.T) {
	p := "./test_dir"
	if err := checkDirExist(p); err != nil {
		t.Error(err)
	}
	os.Remove(p)
}
