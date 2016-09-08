package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
)

type CmdHandler func(args []string, config *Config)

func isNumber(str string) bool {
	if m, _ := regexp.MatchString("^[0-9]+$", str); !m {
		return false
	} else {
		return true
	}
}

func getProfile(arg string, config *Config) *Profile {
	var p *Profile
	if isNumber(arg) {
		n, _ := strconv.Atoi(arg)
		if n >= len(config.Profiles) || n < 0 {
			log.Fatalf("invalid index of profile: [%d]\n", n)
		}
		p = config.Profiles[n]
	} else {
		for _, i := range config.Profiles {
			if i.Name == arg {
				p = i
				break
			}
		}
	}

	return p
}

////////////////////////////////////////////////////////////////////////////////

var handerMap = map[string]CmdHandler{
	"list": cmd_list,
	"make": cmd_make,
}

func cmd_help(args []string, config *Config) {
	fmt.Printf("Usage: \n")
	for k, h := range handerMap {
		fmt.Printf("  - %s\t: ", k)
		h(nil, nil)
	}
	fmt.Println("  - help\t: Display this message.")
}

func cmd_list(args []string, config *Config) {
	if config == nil {
		fmt.Printf("No arg. List all profiles\n")
		return
	}

	for i, p := range config.Profiles {
		fmt.Printf("\n%s[%d]\t: '%s'%s\n", CGREEN, i, p.Name, CEND)
		fmt.Printf("SrcDir\t: %s\nArch\t: %s\nCC\t: %s\n",
			p.SrcDir, p.Arch, p.CrossComile)

	}
}

func cmd_make(args []string, config *Config) {
	if config == nil {
		fmt.Printf("{name | index}. Build kernel specified by name or index\n")
		return
	}

	if len(args) == 0 {
		log.Fatal("make need profile's name or index")
	}

	p := getProfile(args[0], config)
	if p == nil {
		log.Fatalf("can not find profile [%s]\n", args[0])
	}

	makeKernel(p, "uImage")
}
