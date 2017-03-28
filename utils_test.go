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
