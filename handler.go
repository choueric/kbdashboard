package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

type CmdHandler func(args []string, config *Config)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func isNumber(str string) bool {
	if m, _ := regexp.MatchString("^[0-9]+$", str); !m {
		return false
	} else {
		return true
	}
}

func getProfile(arg string, config *Config) (*Profile, int) {
	var p *Profile
	var index int
	if isNumber(arg) {
		n, _ := strconv.Atoi(arg)
		if n >= len(config.Profiles) || n < 0 {
			log.Fatalf("invalid index of profile: [%d]\n", n)
		}
		p = config.Profiles[n]
		index = n
	} else {
		for i, v := range config.Profiles {
			if v.Name == arg {
				p = v
				index = i
				break
			}
		}
	}

	return p, index
}

func getProfileByCurrent(args []string, config *Config) (*Profile, int) {
	var arg string
	if len(args) == 0 {
		arg = strconv.Itoa(config.Current)
	} else {
		arg = args[0]
	}

	return getProfile(arg, config)
}

////////////////////////////////////////////////////////////////////////////////

var handerMap = map[string]CmdHandler{
	"list":    cmd_list,
	"choose":  cmd_choose,
	"edit":    cmd_edit,
	"make":    cmd_make,
	"config":  cmd_config,
	"build":   cmd_build,
	"install": cmd_install,
}

func cmd_help(args []string, config *Config) {
	fmt.Printf("cmd %s'help'%s:\n", CGREEN, CEND)
	fmt.Printf("Usage: \n")
	for k, h := range handerMap {
		fmt.Printf("  - %s\t: ", k)
		h(nil, nil)
	}
	fmt.Println("  - help\t: Display this message.")
}

func cmd_list(args []string, config *Config) {
	if config == nil {
		fmt.Printf("List all profiles.\n")
		return
	}

	fmt.Printf("cmd %s'list'%s:\n", CGREEN, CEND)
	for i, p := range config.Profiles {
		if config.Current == i {
			fmt.Printf("\n%s*%s ", CRED, CEND)
			fmt.Printf("%s[%d]\t: '%s'%s\n", CGREEN, i, p.Name, CEND)
		} else {
			fmt.Printf("\n%s[%d]\t: '%s'%s\n", CGREEN, i, p.Name, CEND)
		}
		fmt.Printf("SrcDir\t: %s\n", p.SrcDir)
		fmt.Printf("Arch\t: %s\n", p.Arch)
		fmt.Printf("CC\t: %s\n", p.CrossComile)
		fmt.Printf("Target\t: %s\n", p.Target)
	}
}

func cmd_make(args []string, config *Config) {
	argc := len(args)
	if config == nil || argc == 0 {
		fmt.Printf("<target> [name | index]. Execute 'make' with specify target.\n")
		return
	}

	target := args[0]
	args = args[1:len(args)]

	p, _ := getProfileByCurrent(args, config)
	if p == nil {
		log.Fatalf("can not find profile [%s]\n", args[0])
	}

	fmt.Printf("cmd %s'build %s'%s for %s[%s]%s\n", CGREEN, target, CEND, CGREEN, p.Name, CEND)
	makeKernel(p, target)
}

func cmd_build(args []string, config *Config) {
	if config == nil {
		fmt.Printf("[name | index]. Build kernel specified by name or index.\n")
		fmt.Printf("\t\t  Same as '$ kbdashboard make uImage' if target in config is uImage.\n")
		return
	}

	p, _ := getProfileByCurrent(args, config)
	if p == nil {
		log.Fatalf("can not find profile [%s]\n", args[0])
	}

	fmt.Printf("cmd %s'build'%s for %s[%s]%s\n", CGREEN, CEND, CGREEN, p.Name, CEND)
	makeKernel(p, p.Target)
}

func cmd_config(args []string, config *Config) {
	if config == nil {
		fmt.Printf("[name | index]. Configure kernel using menuconfig.\n")
		fmt.Printf("\t\t  Same as '$ kbdashboard make menuconfig'.\n")
		return
	}

	p, _ := getProfileByCurrent(args, config)
	if p == nil {
		log.Fatalf("can not find profile [%s]\n", args[0])
	}

	fmt.Printf("cmd %s'config'%s for %s[%s]%s\n", CGREEN, CEND, CGREEN, p.Name, CEND)
	configKernel(p, "menuconfig")
}

func cmd_edit(args []string, config *Config) {
	if config == nil {
		fmt.Printf("Edit the config file using editor specified in config file.\n")
		return
	}

	var argv = []string{config.Editor, config.configFile}
	execCmd(config.Editor, argv)
}

func cmd_choose(args []string, config *Config) {
	if config == nil {
		fmt.Printf("Choose current profile.\n")
		return
	}

	if len(args) == 0 {
		log.Fatal("Choose need profile's name or index")
	}
	p, index := getProfile(args[0], config)
	if p == nil {
		log.Fatalf("can not find profile [%s]\n", args[0])
	}

	fmt.Printf("cmd %s'choose'%s profile %s[%s]%s\n", CGREEN, CEND, CGREEN, p.Name, CEND)
	config.Current = index

	writeConfigFile(config)
}

func cmd_install(args []string, config *Config) {
	if config == nil {
		fmt.Printf("[edit] [name | index]. Execute or edit install script.\n")
		fmt.Printf("\t\t  If use sub-cmd 'edit', open the install script with editor.\n")
		fmt.Printf("\t\t  If no sub-cmd 'edit', execute the install script.\n")
		return
	}

	var doEdit bool
	var create bool

	argc := len(args)
	if argc != 0 && args[0] == "edit" {
		doEdit = true
		args = args[1:len(args)]
	}

	p, _ := getProfileByCurrent(args, config)
	if p == nil {
		log.Fatalf("can not find profile [%s]\n", args[0])
	}

	script := getInstallFilename(p)
	if checkFileExsit(script) == false {
		// create script
		file, err := os.OpenFile(script, os.O_RDWR|os.O_CREATE, 0775)
		checkError(err)
		_, err = file.Write([]byte("#!/bin/sh"))
		checkError(err)
		file.Close()
		create = true
	}

	if doEdit {
		fmt.Printf("cmd %s'install edit'%s profile %s[%s]%s\n",
			CGREEN, CEND, CGREEN, p.Name, CEND)
		var argv = []string{config.Editor, script}
		execCmd(config.Editor, argv)
		return
	}

	fmt.Printf("cmd %s'install'%s profile %s[%s]%s\n", CGREEN, CEND, CGREEN, p.Name, CEND)
	if create {
		// edit script
		var argv = []string{config.Editor, script}
		execCmd(config.Editor, argv)
	} else {
		cmd := exec.Command(script)
		cmd.Dir = p.SrcDir
		fmt.Printf("    %s%s%s\n", CGREEN, script, CEND)
		runCmd(cmd)
	}
}
