package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

const VERSION = "0.2.4"
const GIT_DEF = "@"

var (
	BUILD_TIME = "nil"
	GIT_COMMIT = GIT_DEF
)
var versionHelp = &helpNode{
	cmd:      "version",
	synopsis: "Print version information.",
	usage: func(w io.Writer, h *helpNode) {
		printCmdTitle("version", false)
		printCmdInfo("Output the version information.\n")
	},
}

func versionHandler(args []string, data interface{}) (int, error) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Version\t:", cWrap(cGREEN, VERSION))
	if BUILD_TIME != "nil" {
		fmt.Fprintln(w, "Build time\t:", cWrap(cGREEN, BUILD_TIME))
	}
	if GIT_COMMIT != GIT_DEF {
		fmt.Fprintln(w, "Git Commit\t:", cWrap(cGREEN, GIT_COMMIT))
	}
	w.Flush()
	return 0, nil
}
