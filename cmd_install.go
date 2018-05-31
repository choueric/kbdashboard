package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
)

const (
	builtInComments = ` 
# Built-in variables:
# - KBD_COLOR: global, colorful ouput switch.
# - KBD_DEBUG: global, debug output switch.
# - KBD_EDITOR: global, edior name.
# - KBD_CURRENT: global, current profile index.
# - KBD_NAME: profile, name.
# - KBD_SRC_DIR: profile, source directory.
# - KBD_ARCH: profile, archetect.
# - KBD_CC: profile, cross compiler.
# - KBD_TARGET: profile, build target.
# - KBD_BUILD_DIR: profile, build directory.
# - KBD_DEFCONFIG: profile, default config name.
# - KBD_DTB: profile, DTB target name.
# - KBD_MOD_DIR: profile, modules install directory.
# - KBD_THREAD_NUM: profile, thread number.
`
	scriptContent = `
install -d $DEST_DIR

case "$1" in 
	"image") cp -v $IMAGE $DEST_DIR;;
	"modules") cp -a $MODULES $DEST_DIR;;
	"dtb") cp -v $DTB $DEST_DIR;;
	*) echo "cmds: image|modules|dtb.";;
esac
`
)

var installHelp = &helpNode{
	cmd:      "install",
	synopsis: "Execute your own install script.",
	usage: func(w io.Writer, h *helpNode) {
		printCmdTitle("install [args]", false)
		printSubcmdInfo("Execute your own install script of the current profile.\n")
		printSubcmdInfo("The [args] will transferred to the script.\n")
	},
}

func createScript(fileName string, p *Profile) error {
	fmt.Printf("create install script: '%s'\n", cWrap(cGREEN, fileName))
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return err
	}
	defer file.Close()

	str := fmt.Sprintf("#!/bin/sh\n\n# install script for profile '%s'\n", p.Name)
	str += builtInComments + "\n"
	str += fmt.Sprintf("IMAGE=\"$KBD_BUILD_DIR/arch/%s/boot/%s\"\n", p.Arch, p.Target)
	str += fmt.Sprintf("MODULES=\"$KBD_MOD_DIR/lib/modules/{TODO:uname -r}\"\n")
	str += fmt.Sprintf("DTB=\"$KBD_BUILD_DIR/%s\"\n", p.DTB)
	str += "DEST_DIR=\"{TODO:dirname}\"\n"
	str += scriptContent

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
	return 0, pipeCmd(cmd, os.Stdout, true)
}
