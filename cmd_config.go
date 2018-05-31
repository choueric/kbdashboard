package main

import (
	"io"
	"os"
)

var configHelp = &helpNode{
	cmd:      "config",
	synopsis: "Handle kernel's configuration. [menu|def|save].",
	usage: func(w io.Writer, h *helpNode) {
		cmdTitle(w, true, "config menu")
		cmdUsage(w, "Invoke 'make menuconfig' on the current kernel.\n")

		cmdTitle(w, false, "config def")
		cmdUsage(w, "Invoke 'make defconfig' on the current kernel.\n")

		cmdTitle(w, false, "config save")
		cmdUsage(w, "Save current config as the default config.\n"+
			"Invoke 'make savedefconfig' and then overwrite the original config file.\n")
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
