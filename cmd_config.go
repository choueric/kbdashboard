package main

import (
	"io"
	"os"
)

var configHelp = &helpNode{
	cmd:      "config",
	synopsis: "Handle kernel's configuration. [menu|def|save].",
	usage: func(w io.Writer, h *helpNode) {
		printSubcmdTitle("config menu", true)
		printSubcmdInfo("Invoke 'make menuconfig' on the current kernel.\n")

		printSubcmdTitle("config def", false)
		printSubcmdInfo("Invoke 'make defconfig' on the current kernel.\n")

		printSubcmdTitle("config save", false)
		printSubcmdInfo("Save current config as the default config.\n")
		printSubcmdInfo("Invoke 'make savedefconfig' and then overwrite the original config file.\n")
	},
}

func configMenuHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}
	printCmd("config menu", p.Name)
	return 0, configKernel(p, "menuconfig")
}

func configDefHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}
	printCmd("config def", p.Name)
	return 0, makeKernel(p, p.DefConfig, os.Stdout, true)
}

func configSaveHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}

	printCmd("config save", p.Name)
	if err := makeKernel(p, "savedefconfig", os.Stdout, true); err != nil {
		return 0, err
	}

	src := p.BuildDir + "/defconfig"
	dst := p.SrcDir + "/arch/" + p.Arch + "/configs/" + p.DefConfig

	if err := copyFileContents(src, dst); err != nil {
		return 0, err
	}
	return 0, nil
}
