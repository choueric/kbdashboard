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

const DefaultConfig = `
{
	"color": true
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
}
`

var gConfig *Config

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

/*
 * get profile from @config by @arg.
 * @arg may be numberic index or name of profile, it cannot be empty.
 */
func doGetProfile(arg string, config *Config) (*Profile, int, error) {
	var p *Profile
	var index int
	if isNumber(arg) {
		n, _ := strconv.Atoi(arg)
		if n >= len(config.Profiles) || n < 0 {
			return nil, -1, errors.New(fmt.Sprintf("invalid profile index: [%d/%d].",
				n, len(config.Profiles)))
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

type Config struct {
	Editor   string    `json:"editor"`
	Current  int       `json:"current"`
	Color    bool      `json:"color"`
	Profiles []Profile `json:"profile"`
	filepath string
	jc       interface{} // must be interface{} or panic
}

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
