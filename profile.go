/*
 * Copyright (C) 2016 Eric Chou <zhssmail@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
	Defconfig     string `json:"defconfig"`
	DTB           string `json:"dtb"`
	CrossComile   string `json:"cross_compile"`
	OutputDir     string `json:"output_dir"`
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
	line += fmt.Sprintln("  Defconfig\t:", p.Defconfig)
	line += fmt.Sprintln("  DTB\t:", p.DTB)
	line += fmt.Sprintln("  BuildDir\t:", p.OutputDir)
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
