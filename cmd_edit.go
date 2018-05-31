package main

import "io"

var editHelp = &helpNode{
	cmd:      "edit",
	synopsis: "Edit profiles or scripts with the 'editor'. [profile|install].",
	usage: func(w io.Writer, h *helpNode) {
		printSubcmdTitle("edit profile", true)
		printSubcmdInfo("Edit the kbdashboard's configuration file.\n")

		printSubcmdTitle("edit install", false)
		printSubcmdInfo("Edit current profile's installation script.\n")
	},
}

func editProfileHandler(args []string, data interface{}) (int, error) {
	var argv = []string{gConfig.Editor, gConfig.filepath}
	return 0, execCmd(gConfig.Editor, argv)
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
