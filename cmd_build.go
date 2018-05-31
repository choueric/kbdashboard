package main

import (
	"io"
	"os"
)

var buildHelp = &helpNode{
	cmd:      "build",
	synopsis: "Build various targets of kernel. [image|modules|dtb].",
	usage: func(w io.Writer, h *helpNode) {
		printSubcmdTitle("build image", true)
		printSubcmdInfo("Build kernel images for current profile.\n")
		printSubcmdInfo("Equal to '$kbdashboard make uImage'.\n")

		printSubcmdTitle("build modules", false)
		printSubcmdInfo("Build and install modules for current profile.\n")
		printSubcmdInfo("Eqaul to '$ make modules' then '$ make modules_install'.\n")

		printSubcmdTitle("build dtb", false)
		printSubcmdInfo("Build 'dtb' file and install into 'BuildDir'.\n")
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

	src := p.BuildDir + "/arch/" + p.Arch + "/boot/dts/" + p.DTB
	dst := p.BuildDir + "/" + p.DTB

	if err := copyFileContents(src, dst); err != nil {
		return 0, err
	}
	return 0, nil
}
