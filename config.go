package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/choueric/jconfig"
)

type Config struct {
	Editor   string    `json:"editor"`
	Current  int       `json:"current"`
	Color    bool      `json:"color"`
	Profiles []Profile `json:"profile"`
	Debug    bool      `json:"debug"`
	filepath string
	jc       interface{} // must be interface{} or panic
}

const DefaultConfig = `{
	"color": true,
	"debug": false,
	"editor": "vim",
	"current": 0,
	"profile": [
	{
		"name":"profile_name",
		"src_dir":"/path/of/the/kernel/source",
		"arch":"$ARCH: arm ",
		"cross_compile":"$CROSS_COMPILE: arm-eabi-",
		"build_dir":"$O: ./_build",
		"defconfig":"at91rm9200_defconfig",
		"target":"uImage",
		"dtb":"at91rm9200ek.dtb",
		"mod_install_dir":"$INSTALL_MOD_PATH: ./_build/modules",
		"thread_num":4
	}
	]
} `

var gConfig *Config

func (c *Config) String() string {
	line := fmt.Sprintf("Config File\t: %s\nEditor\t\t: %s\nColor\t\t: %v\n"+
		"Current Profile\t: %d\n",
		c.filepath, c.Editor, c.Color, c.Current)
	for _, v := range c.Profiles {
		line += v.Name + "\n" + v.String() + "\n"
	}
	return line
}

func (c *Config) fix() {
	// validate config
	if c.Current >= len(c.Profiles) {
		fmt.Fprintf(os.Stderr, "Current in config.json is invalid: %s\n", c.Current)
		os.Exit(1)
	}

	// fix invaid configurations
	for _, p := range c.Profiles {
		p.BuildDir = fixRelativeDir(p.BuildDir, p.SrcDir)
		p.ModInstallDir = fixRelativeDir(p.ModInstallDir, p.SrcDir)
		if p.DefConfig == "" {
			p.DefConfig = "defconfig"
		}
	}
}

func (c *Config) dump() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "%v", c)
	w.Flush()
}

func (c *Config) save() {
	c.getJc().Save()
}

func (c *Config) getInstallFilename(p *Profile) string {
	return c.getJc().Dir() + "/" + p.Name + "_install.sh"
}

func (c *Config) getJc() *jconfig.JConfig {
	return c.jc.(*jconfig.JConfig)
}

func getConfig(dump bool) *Config {
	filepath := os.Getenv("HOME") + "/.kbdashboard/config.json"
	jc := jconfig.New(filepath, Config{})

	if _, err := jc.Load(DefaultConfig); err != nil {
		fmt.Fprintf(os.Stderr, "load config error: %v", err)
		os.Exit(1)
	}

	config := jc.Data().(*Config)
	config.jc = jc
	config.filepath = jc.FilePath()
	config.fix()

	if dump {
		config.dump()
	}
	return config
}

/*
 * get profile from @config by @arg.
 * @arg may be numberic index or name of profile, it cannot be empty.
 */
func doGetProfile(arg string, config *Config) (*Profile, int, error) {
	var p *Profile
	var index int
	if isNumber(arg) {
		n, _ := strconv.Atoi(arg)
		total := len(config.Profiles)
		if n >= total || n < 0 {
			errMsg := fmt.Sprintf("invalid profile index: [%d/%d].", n, total)
			return nil, -1, errors.New(errMsg)
		}
		p = &config.Profiles[n]
		index = n
	} else {
		for i, v := range config.Profiles {
			if v.Name == arg {
				p = &v
				index = i
				break
			}
		}
	}

	return p, index, nil
}

/*
 * get current profile in @config
 */
func getCurrentProfile(config *Config) (*Profile, int, error) {
	profile := strconv.Itoa(config.Current)
	return doGetProfile(profile, config)
}
