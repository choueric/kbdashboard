package main

import (
	"io"
	"os"
)

var buildHelp = &helpNode{
	cmd:      "build",
	synopsis: "Build various targets of kernel.",
	subs: []helpSubNode{
		{"image", func(w io.Writer) {
			cmdTitle(w, true, "build image")
			cmdUsage(w, "Build kernel images for current profile.\n"+
				"Equal to '$kbdashboard make uImage'.\n")
		}},
		{"modules", func(w io.Writer) {
			cmdTitle(w, false, "build modules")
			cmdUsage(w, "Build and install modules for current profile.\n"+
				"Eqaul to '$ make modules' then '$ make modules_install'.\n")
		}},
		{"dtb", func(w io.Writer) {
			cmdTitle(w, false, "build dtb")
			cmdUsage(w, "Build 'dtb' file and install into 'BuildDir'.\n")
		}},
	},
}

func buildImageHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}
	printCmd("build image", p.Name)
	return 0, makeKernel(p, p.Target, os.Stdout, true)
}

func buildModulesHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}
	printCmd("build modules", p.Name)

	err = makeKernel(p, "modules", os.Stdout, true)
	if err != nil {
		return 0, err
	}

	return 0, makeKernel(p, "modules_install", os.Stdout, true)
}

func buildDtbHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}
	printCmd("build DTB", p.Name)

	if err := makeKernel(p, p.DTB, os.Stdout, true); err != nil {
		return 0, err
	}

	return 0, nil
}
