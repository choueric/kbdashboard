package main

import (
	"errors"
	"os"
)

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
