package main

import (
	"flag"
	"io"
)

func listUsage(w io.Writer, m *helpMap) {
	cmdTitle("list [-v|-c]", false)
	cmdInfo("List all profiles.\n")
	cmdInfo("-v: Print with more information\n")
	cmdInfo("-c: Print full information of current profile.\n\n")
}

func listHandler(args []string, data interface{}) (int, error) {
	var verbose, current bool
	flagSet := flag.NewFlagSet("list", flag.ExitOnError)
	flagSet.BoolVar(&verbose, "v", false, "print with more information")
	flagSet.BoolVar(&current, "c", false, "print all information of current profile")
	flagSet.Parse(args)

	if current {
		p := gConfig.Profiles[gConfig.Current]
		printProfile(p, true, true, gConfig.Current)
		return 0, nil
	}

	for i, p := range gConfig.Profiles {
		printProfile(p, verbose, gConfig.Current == i, i)
	}

	return 0, nil
}
