package main

import (
	"errors"
	"io"
	"os"
)

var makeHelp = &helpNode{
	cmd:      "make",
	synopsis: "Execute '$ make <target>' on current kernel.",
	usage: func(w io.Writer) {
		cmdTitle(w, false, "make <target>")
		cmdUsage(w, "This is the way to pass through kbdashboard and invoke\n"+
			"kernel's make directly.\n"+
			"So the <target> is just as the same as kernel's own\n"+
			"Makefile's target.\n")
	},
}

func makeHandler(args []string, data interface{}) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("'make' need paramters.")
	}
	target := args[0]
	args = args[1:]

	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}

	printCmd("make "+target, p.Name)
	return 0, makeKernel(p, target, os.Stdout, true)
}
