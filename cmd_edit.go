package main

import (
	"fmt"
	"io"
)

func editUsage(w io.Writer, m *helpMap) {
	defaultHelp(w, m)
	fmt.Printf("\n")
	editProfileUsage()
	editInstallUsage()
	fmt.Printf("\n")
}

////////////////////////////////////////////////////////////////////////////////

func editProfileUsage() {
	subcmdTitle("edit profile", true)
	subcmdInfo("Edit the kbdashboard's configuration file.\n")
}

func editProfileHandler(args []string, data interface{}) (int, error) {
	var argv = []string{gConfig.Editor, gConfig.filepath}
	return 0, execCmd(gConfig.Editor, argv)
}

////////////////////////////////////////////////////////////////////////////////

func editInstallUsage() {
	subcmdTitle("edit install", false)
	subcmdInfo("Edit current profile's installation script.\n")
}

func editInstallHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}
	var argv = []string{gConfig.Editor, gConfig.getInstallFilename(p)}
	return 0, execCmd(gConfig.Editor, argv)
}
