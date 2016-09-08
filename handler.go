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
		fmt.Printf("\n%s[%d]\t%s: %v\n", CGREEN, i, CEND, p)
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

	var p *Profile
	if isNumber(args[0]) {
		n, _ := strconv.Atoi(args[0])
		p = config.Profiles[n]
	} else {
		for _, i := range config.Profiles {
			if i.Name == args[0] {
				p = i
				break
			}
		}
	}

	if p == nil {
		log.Fatalf("can not find profile [%s]\n", args[0])
	}

	makeKernel(p, "uImage")
}
