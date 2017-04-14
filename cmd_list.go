package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
)

func listUsage(w io.Writer, m *helpMap) {
	printCmdTitle("list [-v|-a]", false)
	printCmdInfo("List all profiles.\n")
	printCmdInfo("-v: Print with more information\n")
	printCmdInfo("-c: Print full information of current profile.\n\n")
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
