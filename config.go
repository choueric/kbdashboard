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
		"name":"demo",
		"src_dir":"/home/user/kernel",
		"arch":"arm",
		"cross_compile":"arm-eabi-",
		"target":"uImage",
		"output_dir":"./_build",
		"defconfig":"at91rm9200_defconfig",
		"dtb":"at91rm9200ek.dtb",
		"mod_install_dir":"./_build/mod",
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
		p.OutputDir = fixRelativeDir(p.OutputDir, p.SrcDir)
		p.ModInstallDir = fixRelativeDir(p.ModInstallDir, p.SrcDir)
		if p.Defconfig == "" {
			p.Defconfig = "defconfig"
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
