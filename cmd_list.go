package main

import (
	"flag"
	"fmt"
	"io"
)

var listHelp = &helpNode{
	cmd:      "list",
	synopsis: "List profiles' information.",
	usage: func(w io.Writer) {
		cmdTitle(w, false, "list [-a|-v]")
		cmdUsage(w, "Without any options, print current profile's brief details.\n"+
			"-a: Print with all profile's brief details.\n"+
			"-v: Print current profile's full information.\n")
	},
}

func listHandler(args []string, data interface{}) (int, error) {
	var verbose, all bool
	flagSet := flag.NewFlagSet("list", flag.ExitOnError)
	flagSet.BoolVar(&verbose, "v", false, "print with more information.")
	flagSet.BoolVar(&all, "a", false, "print all profiles.")
	flagSet.Parse(args)

	printAdditional := func(p *Profile) {
		version, err := kernelFullVersion(p)
		if err == nil {
			fmt.Println("  Version   :", version)
		}
	}

	if !all {
		p := gConfig.Profiles[gConfig.Current]
		printProfile(p, verbose, true, gConfig.Current)
		if verbose {
			printAdditional(p)
		}
		return 0, nil
	}

	for i, p := range gConfig.Profiles {
		printProfile(p, verbose, gConfig.Current == i, i)
		if verbose {
			printAdditional(p)
		}
	}

	return 0, nil
}
