package main

import (
	"strings"
	"testing"
)

func Test_doMakeKernelOpt(t *testing.T) {
	target := "build"
	n := 4
	b := "./build_dir"
	c := "arm-"
	a := "arm"
	m := "./mod_install"
	extra := ""

	expect := "build -j4 O=./build_dir CROSS_COMPILE=arm- ARCH=arm INSTALL_MOD_PATH=./mod_install"
	output := doMakeKernelOpt(target, n, b, c, a, m, extra)
	o := strings.Join(output, " ")
	if o != expect {
		t.Error("output is wrong:", o)
	}

	target = "config"
	n = 2
	b = ""
	c = "arm-"
	a = ""
	m = ""
	extra = "CFLAGS_KERNEL=-march=armv7-a   CFLAGS_GCOV=-fprofile-arcs"

	expect = "config -j2 CROSS_COMPILE=arm- CFLAGS_KERNEL=-march=armv7-a CFLAGS_GCOV=-fprofile-arcs"
	output = doMakeKernelOpt(target, n, b, c, a, m, extra)
	o = strings.Join(output, " ")
	if o != expect {
		t.Error("output is wrong:", o)
	}
}
