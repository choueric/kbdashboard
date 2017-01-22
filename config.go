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

	"github.com/choueric/jconfig"
)

const DefaultConfig = `
{
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

func (p *Profile) String() string {
	return fmt.Sprintf(
		"name = %s%s%s\n"+
			"  arch = %s, target = %s, defconfig = %s\n  DTB = %s\n"+
			"  CC = %s\n"+
			"  src_dir = %s\n  build_dir = %s\n  mod_dir = %s\n  thread num = %d\n",
		CGREEN, p.Name, CEND, p.Arch, p.Target, p.Defconfig, p.DTB, p.CrossComile,
		p.SrcDir, p.OutputDir, p.ModInstallDir, p.ThreadNum)
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
	if verbose {
		if current {
			fmt.Printf("\n%s*%s ", CRED, CEND)
			fmt.Printf("%s[%d]\t: '%s'%s\n", CGREEN, i, p.Name, CEND)
		} else {
			fmt.Printf("\n  %s[%d]\t: '%s'%s\n", CGREEN, i, p.Name, CEND)
		}
		fmt.Printf("  SrcDir\t\t: %s\n", p.SrcDir)
		fmt.Printf("  Arch\t\t\t: %s\n", p.Arch)
		fmt.Printf("  CC\t\t\t: %s\n", p.CrossComile)
		fmt.Printf("  Target\t\t: %s\n", p.Target)
		fmt.Printf("  Defconfig\t\t: %s\n", p.Defconfig)
		fmt.Printf("  DTB\t\t\t: %s\n", p.DTB)
		fmt.Printf("  BuildDir\t\t: %s\n", p.OutputDir)
		fmt.Printf("  ModInsDir\t\t: %s\n", p.ModInstallDir)
		fmt.Printf("  ThreadNum\t\t: %d\n", p.ThreadNum)
	} else {
		if current {
			fmt.Printf("\n%s*%s ", CRED, CEND)
			fmt.Printf("%s[%d]\t: '%s'%s\n", CGREEN, i, p.Name, CEND)
		} else {
			fmt.Printf("\n  %s[%d]\t: '%s'%s\n", CGREEN, i, p.Name, CEND)
		}
		fmt.Printf("  SrcDir: %s\n", p.SrcDir)
		fmt.Printf("  Arch\t: %s\n", p.Arch)
		fmt.Printf("  CC\t: %s\n", p.CrossComile)
	}
}

type Config struct {
	Editor   string    `json:"editor"`
	Current  int       `json:"current"`
	Profiles []Profile `json:"profile"`
	filepath string
	jc       interface{} // must be interface{} or panic
}

func (c *Config) String() string {
	line := fmt.Sprintf("Config File\t:%s\nEditor\t\t:%s\nCurrent Profile\t:%d\n",
		c.filepath, c.Editor, c.Current)
	for _, v := range c.Profiles {
		line += v.String()
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
		fmt.Println(config)
	}
	return config
}
