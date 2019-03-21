package main

import (
	"bytes"
	"testing"
	"text/tabwriter"
)

const helpOutput = `  list       : List profiles' information.
  choose     : Choose <profile> as the current profile.
  edit       : Edit profiles or scripts with the 'editor'. [profile|install].
  config     : Handle kernel's configuration. [menu|def|save].
  build      : Build various targets of kernel. [image|modules|dtb].
  install    : Execute your own install script.
  make       : Execute '$ make <target>' on current kernel.
  dts        : List relevant DTS files.
  version    : Print version information.
  completion : Generate a shell completion file.
  help       : Print help message for one or all commands. [<cmd>].
`

func Test_helpOutput(t *testing.T) {
	var output bytes.Buffer
	w := tabwriter.NewWriter(&output, 0, 0, 1, ' ', tabwriter.TabIndent)
	for _, h := range helpJar {
		outputSynopsis(w, h)
	}
	w.Flush()

	if helpOutput != output.String() {
		t.Errorf("expect:\n%s\noutput:\n%s\n", helpOutput, output.String())
	}
}

const configHelpOutput = `Usage of comamnd '[33;1mconfig[0;m':
[31;1m*[0;m [32;1mconfig menu[0;m
    Invoke 'make menuconfig' on the current kernel.
  [32;1mconfig def[0;m
    Invoke 'make defconfig' on the current kernel.
  [32;1mconfig save[0;m
    Save current config as the default config.
    Invoke 'make savedefconfig' and then overwrite the original config file.
`

func Test_configHelpOutput(t *testing.T) {
	gConfig = getConfig(false)
	var output bytes.Buffer
	w := tabwriter.NewWriter(&output, 0, 0, 1, ' ', tabwriter.TabIndent)
	outputBanner(w, configHelp)
	outputUsage(w, configHelp)
	w.Flush()

	if configHelpOutput != output.String() {
		t.Errorf("expect:\n%s\noutput:\n%s\n", helpOutput, output.String())
	}
}
