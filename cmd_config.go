package main

import (
	"fmt"
	"io"
	"os"
)

func configUsage(w io.Writer, m *helpMap) {
	defaultHelp(w, m)
	fmt.Printf("\n")
	configMenuUsage()
	configDefUsage()
	configSaveUsage()
	fmt.Printf("\n")
}

////////////////////////////////////////////////////////////////////////////////

func configMenuUsage() {
	subcmdTitle("config menu", true)
	subcmdInfo("Invoke 'make menuconfig' on the current kernel.\n")
}

func configMenuHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}
	printCmd("config menu", p.Name)
	return 0, configKernel(p, "menuconfig")
}

////////////////////////////////////////////////////////////////////////////////

func configDefUsage() {
	subcmdTitle("config def", false)
	subcmdInfo("Invoke 'make defconfig' on the current kernel.\n")
}

func configDefHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}
	printCmd("config def", p.Name)
	return 0, makeKernel(p, p.DefConfig, os.Stdout, true)
}

////////////////////////////////////////////////////////////////////////////////

func configSaveUsage() {
	subcmdTitle("config save", false)
	subcmdInfo("Save current config as the default config.\n")
	subcmdInfo("Invoke 'make savedefconfig' and then overwrite the original config file.\n")
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
