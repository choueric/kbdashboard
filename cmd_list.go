package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
)

var listHelp = &helpNode{
	cmd:      "list",
	synopsis: "List profiles' information.",
	usage: func(w io.Writer, h *helpNode) {
		cmdTitle(w, false, "list [-a|-v]")
		cmdUsage(w, "Without any options, print current profile's brief details.\n"+
			"-a: Print with all profile's brief details.\n"+
			"-v: Print current profile's full information.\n")
	},
}

func printAdditional(p *Profile) {
	var result bytes.Buffer

	err := makeKernel(p, "kernelversion", &result, false)
	if err == nil {
		fmt.Printf("  Version   : %s", result.String())
	}
}

func listHandler(args []string, data interface{}) (int, error) {
	var verbose, all bool
	flagSet := flag.NewFlagSet("list", flag.ExitOnError)
	flagSet.BoolVar(&verbose, "v", false, "print with more information.")
	flagSet.BoolVar(&all, "a", false, "print all profiles.")
	flagSet.Parse(args)

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
