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
	"errors"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/choueric/cmdmux"
)

type helpMap struct {
	cmd      string
	usage    func(w io.Writer, m *helpMap)
	synopsis string
}

var cmdHelpMap = []helpMap{
	{"list", listUsage, "List all profiles."},
	{"choose", defaultHelp, "Choose <profile> as the current profile."},
	{"edit", editUsage, "Edit profiles or scripts with the 'editor'. [profile|install]."},
	{"config", configUsage, "Handle kernel's configuration. [menu|def|save]."},
	{"build", buildUsage, "Build various targets of kernel. [image|modules|dtb]."},
	{"install", defaultHelp, "Execute the install script of current profile with [option]."},
	{"make", defaultHelp, "Execute '$ make <target>' on current kernel."},
	{"dts", dtsUsage, "List relevant DTS files."},
	{"version", defaultHelp, "Print the version."},
	{"completion", completionUsage, "Generate a shell completion file."},
	{"help", defaultHelp, "Print help message for one or all commands. [cmd]."},
}

var COMP_FILENAME = "kbdashboard.bash-completion"

func defaultHelp(w io.Writer, m *helpMap) {
	fmt.Fprintf(w, "  %s\t: %s\n", m.cmd, m.synopsis)
}

func helpHandler(args []string, data interface{}) (int, error) {
	if len(args) == 0 {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
		for _, v := range cmdHelpMap {
			fmt.Fprintf(w, "  %s\t: %s\n", v.cmd, v.synopsis)
		}
		w.Flush()
		return 1, nil
	}

	cmd := args[0]
	for _, v := range cmdHelpMap {
		if cmd == v.cmd {
			fmt.Printf("Usage of comamnd '%s':\n\n", cWrap(cYELLOW, cmd))
			v.usage(os.Stdout, &v)
			return 1, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("invalid command '%s'.", cmd))
}

func completionUsage(w io.Writer, m *helpMap) {
	cmdInfo("Generate a shell completion file '%s'.\n\n", COMP_FILENAME)
}

func completionHandler(args []string, data interface{}) (int, error) {
	file, err := os.Create(COMP_FILENAME)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	if err = cmdmux.GenerateCompletion("kbdashboard", file); err != nil {
		return 0, err
	}
	fmt.Printf("Create completion file '%s' OK.\n", COMP_FILENAME)

	return 0, nil
}
