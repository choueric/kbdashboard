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

	script := gConfig.getInstallFilename(p)
	ret, err := checkFileExsit(script)
	if err != nil {
		return 0, err
	}

	if ret == false {
		// create and edit script
		if err := createScript(script, p); err != nil {
			return 0, err
		}
	}
	return 0, execCmd(gConfig.Editor, []string{gConfig.Editor, script})
}
