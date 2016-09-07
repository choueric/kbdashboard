package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/choueric/kernelBuildDashboard/kbd"
)

type CmdHandler func(args []string, config *kbd.Config)

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

func cmd_list(args []string, config *kbd.Config) {
	for i, item := range config.Items {
		fmt.Printf("\n[%d]\t: %v\n", i, item)
	}
}

func cmd_make(args []string, config *kbd.Config) {
	if len(args) == 0 {
		log.Fatal("make need item name or number")
	}

	var item *kbd.Item
	if isNumber(args[0]) {
		n, _ := strconv.Atoi(args[0])
		item = config.Items[n]
	} else {
		for _, i := range config.Items {
			if i.Name == args[0] {
				item = i
				break
			}
		}
	}

	if item == nil {
		log.Fatalf("can not find item [%s]\n", args[0])
	}

	makeKernel(item, "uImage")
}
