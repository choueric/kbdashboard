package main

import (
	"os"
	"testing"
)

func Test_checkDirExist(t *testing.T) {
	p := "./test_dir"
	if err := checkDirExist(p); err != nil {
		t.Error(err)
	}
	os.Remove(p)
}
