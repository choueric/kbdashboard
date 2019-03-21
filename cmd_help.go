package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/choueric/cmdmux"
)

var COMP_FILENAME = "kbdashboard.bash-completion"

type helpSubNode struct {
	name  string
	usage func(w io.Writer)
}

/*
 * One cmd's help output:
 *   ```
 *   1. Banner: "Usage of comamnd '${cmd}':"
 *   2. Synopsis: "${cmd}: ${synopsis}"
 *   3. Usage: "${usage}"
 *   ```
 *
 * All cmds's help output:
 *   ```
 *   Synopsis: "${cmd}: ${synopsis}"
 *   ...
 *   Synopsis: "${cmd}: ${synopsis}"
 *   ```
 */
type helpNode struct {
	cmd      string // one line
	synopsis string // one line
	subs     []helpSubNode
	usage    func(w io.Writer)
}

var helpJar = []*helpNode{
	listHelp, chooseHelp, editHelp, configHelp, buildHelp,
	installHelp, makeHelp, dtsHelp, versionHelp, completionHelp,
	&helpNode{
		cmd:      "help",
		synopsis: "Print help message for one or all commands.",
		subs: []helpSubNode{
			{"<cmd>", func(w io.Writer) {
				cmdTitle(w, false, "help [cmd]")
				cmdUsage(w, "Print the all commands or [cmd]'s help message.\n")
			}},
		},
	},
}

func outputBanner(w io.Writer, h *helpNode) {
	fmt.Fprintf(w, "Usage of comamnd '%s':\n", cWrap(cYELLOW, h.cmd))
}

func outputSynopsis(w io.Writer, h *helpNode) {
	fmt.Fprintf(w, "  %s\t: %s", h.cmd, h.synopsis)
	if len(h.subs) != 0 {
		str := " ["
		for i, sub := range h.subs {
			if i != 0 {
				str += "|"
			}
			str += sub.name

		}
		str += "]."
		fmt.Fprintf(w, "%s", str)
	}
	fmt.Fprintf(w, "\n")
}

func outputUsage(w io.Writer, h *helpNode) {
	if h.usage != nil {
		h.usage(w)
		return
	}

	for _, sub := range h.subs {
		sub.usage(w)
	}
}

func helpHandler(args []string, data interface{}) (int, error) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	if len(args) == 0 {
		for _, h := range helpJar {
			outputSynopsis(w, h)
		}
		w.Flush()
		return 1, nil
	}

	cmd := args[0]
	for _, h := range helpJar {
		if cmd != h.cmd {
			continue
		}
		outputBanner(w, h)
		outputUsage(w, h)
		w.Flush()
		return 1, nil
	}
	return 0, errors.New(fmt.Sprintf("invalid command '%s'.", cmd))
}

var completionHelp = &helpNode{
	cmd:      "completion",
	synopsis: "Generate a shell completion file.",
	usage: func(w io.Writer) {
		cmdTitle(w, false, "completion")
		cmdUsage(w, "Generate a shell completion file '%s'.\n", COMP_FILENAME)
	},
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
