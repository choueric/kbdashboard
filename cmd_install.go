package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func createScript(fileName string, p *Profile) error {
	fmt.Printf("create install script: '%s'\n", cWrap(cGREEN, fileName))
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return err
	}
	defer file.Close()

	str := fmt.Sprintf("#!/bin/sh\n\n# install script for profile '%s'\n", p.Name)
	str += fmt.Sprintf("# Built-in variables:\n")
	str += fmt.Sprintf("# - KBD_COLOR: global, colorful ouput switch.\n")
	str += fmt.Sprintf("# - KBD_DEBUG: global, debug output switch.\n")
	str += fmt.Sprintf("# - KBD_EDITOR: global, edior name.\n")
	str += fmt.Sprintf("# - KBD_CURRENT: global, current profile index.\n")
	str += fmt.Sprintf("# - KBD_NAME: profile, name.\n")
	str += fmt.Sprintf("# - KBD_SRC_DIR: profile, source directory.\n")
	str += fmt.Sprintf("# - KBD_ARCH: profile, archetect.\n")
	str += fmt.Sprintf("# - KBD_CC: profile, cross compiler.\n")
	str += fmt.Sprintf("# - KBD_TARGET: profile, build target.\n")
	str += fmt.Sprintf("# - KBD_BUILD_DIR: profile, build directory.\n")
	str += fmt.Sprintf("# - KBD_DEFCONFIG: profile, default config name.\n")
	str += fmt.Sprintf("# - KBD_DTB: profile, DTB target name.\n")
	str += fmt.Sprintf("# - KBD_MOD_DIR: profile, modules install directory.\n")
	str += fmt.Sprintf("# - KBD_THREAD_NUM: profile, thread number.\n")

	if _, err = file.Write([]byte(str)); err != nil {
		return err
	}

	return nil
}

func buildEnviron(cmd *exec.Cmd, p *Profile) {
	if cmd.Env == nil {
		cmd.Env = make([]string, len(os.Environ()))
		copy(cmd.Env, os.Environ())
	}

	addEnv := func(key, val string) {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, val))
	}

	logger.Println(gConfig.Color)
	addEnv("KBD_COLOR", strconv.FormatBool(gConfig.Color))
	addEnv("KBD_DEBUG", strconv.FormatBool(gConfig.Debug))
	addEnv("KBD_EDITOR", gConfig.Editor)
	addEnv("KBD_CURRENT", strconv.Itoa(gConfig.Current))
	addEnv("KBD_NAME", p.Name)
	addEnv("KBD_SRC_DIR", p.SrcDir)
	addEnv("KBD_ARCH", p.Arch)
	addEnv("KBD_CC", p.CrossComile)
	addEnv("KBD_TARGET", p.Target)
	addEnv("KBD_BUILD_DIR", p.BuildDir)
	addEnv("KBD_DEFCONFIG", p.DefConfig)
	addEnv("KBD_DTB", p.DTB)
	addEnv("KBD_MOD_DIR", p.ModInstallDir)
	addEnv("KBD_TREAD_NUM", strconv.Itoa(p.ThreadNum))
}

func installHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}

	script := gConfig.getInstallFilename(p)
	ret, err := checkFileExsit(script)
	if err != nil {
		return 0, err
	}

	if ret == false {
		// create and edit script
		if err := createScript(script, p); err != nil {
			return 0, err
		}
		return 0, execCmd(gConfig.Editor, []string{gConfig.Editor, script})
	}

	printCmd("install", p.Name)
	logger.Printf("%s\n", cWrap(cGREEN, script))
	cmd := exec.Command(script, args...) // args for script.
	buildEnviron(cmd, p)
	cmd.Dir = p.SrcDir
	return 0, pipeCmd(cmd)
}
