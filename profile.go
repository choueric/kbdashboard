package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type Profile struct {
	Name          string `json:"name"`
	SrcDir        string `json:"src_dir"`
	Arch          string `json:"arch"`
	Target        string `json:"target"`
	DefConfig     string `json:"defconfig"`
	DTB           string `json:"dtb"`
	CrossComile   string `json:"cross_compile"`
	BuildDir      string `json:"build_dir"`
	ModInstallDir string `json:"mod_install_dir"`
	ThreadNum     int    `json:"thread_num"`
}

// do not include 'Name' filed.
func (p *Profile) String() string {
	line := ""
	line += fmt.Sprintln("  SrcDir\t:", p.SrcDir)
	line += fmt.Sprintln("  Arch\t:", p.Arch)
	line += fmt.Sprintln("  CC\t:", p.CrossComile)
	line += fmt.Sprintln("  Target\t:", p.Target)
	line += fmt.Sprintln("  DefConfig\t:", p.DefConfig)
	line += fmt.Sprintln("  DTB\t:", p.DTB)
	line += fmt.Sprintln("  BuildDir\t:", p.BuildDir)
	line += fmt.Sprintln("  ModInsDir\t:", p.ModInstallDir)
	line += fmt.Sprintln("  ThreadNum\t:", p.ThreadNum)
	return line
}

func printProfile(p *Profile, verbose bool, current bool, i int) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	header := func(p *Profile, current bool) {
		if current {
			fmt.Printf("\n" + defMark() + " ")
		} else {
			fmt.Printf("\n  ")
		}
		fmt.Println(cWrap(cGREEN, fmt.Sprintf("[%d] '%s'", i, p.Name)))
	}
	if verbose {
		header(p, current)
		fmt.Fprintf(w, "%v", p)
	} else {
		header(p, current)
		fmt.Fprintf(w, "  SrcDir\t: %s\n", p.SrcDir)
		fmt.Fprintf(w, "  Arch\t: %s\n", p.Arch)
		fmt.Fprintf(w, "  CC\t: %s\n", p.CrossComile)
	}
	w.Flush()
}
