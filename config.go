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
	"fmt"
	"os"
	"strconv"

	"github.com/choueric/clog"
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

var (
	gJConfig     *jconfig.JConfig
	defConfigDir = os.Getenv("HOME") + "/.config/kbdashboard"
)

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

type Config struct {
	Editor   string     `json:"editor"`
	Current  int        `json:"current"`
	Profiles []*Profile `json:"profile"`
	filepath string
}

func (p *Profile) String() string {
	return fmt.Sprintf(
		"name = %s%s%s\n"+
			"  arch = %s, CC = %s, target = %s, defconfig = %s, DTB = %s\n"+
			"  src_dir = %s\n  build_dir = %s, mod_dir = %s\n  thread num = %d\n",
		CGREEN, p.Name, CEND, p.Arch, p.CrossComile, p.Target, p.Defconfig, p.DTB,
		p.SrcDir, p.OutputDir, p.ModInstallDir, p.ThreadNum)
}

func (c *Config) String() string {
	return fmt.Sprintf("Config File\t:%s\nEditor\t\t:%s\nCurrent Profile\t:%d\n%v\n",
		c.filepath, c.Editor, c.Current, c.Profiles)
}

func FixConfig(c *Config) {
	// validate config
	if c.Current >= len(c.Profiles) {
		clog.Fatal("Current in config.json is invalid: ", c.Current)
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

/*
 * get profile from @config by @arg.
 * @arg may be numberic index or name of profile, it cannot be empty.
 */
func doGetProfile(arg string, config *Config) (*Profile, int) {
	var p *Profile
	var index int
	if isNumber(arg) {
		n, _ := strconv.Atoi(arg)
		if n >= len(config.Profiles) || n < 0 {
			clog.Fatalf("invalid index of profile: [%d]\n", n)
		}
		p = config.Profiles[n]
		index = n
	} else {
		for i, v := range config.Profiles {
			if v.Name == arg {
				p = v
				index = i
				break
			}
		}
	}

	return p, index
}

/*
 * get profile specified by @arg from @config.
 * @arg may be numberic index or name of profile.
 *      If it is empty, then return the chosen profile.
 */
func getProfile(args []string, config *Config) (*Profile, int) {
	var arg string
	if len(args) == 0 {
		arg = strconv.Itoa(config.Current)
	} else {
		arg = args[0]
	}

	return doGetProfile(arg, config)
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

func getInstallFilename(p *Profile) string {
	return defConfigDir + "/" + p.Name + "_install.sh"
}

func getConfig(dump bool) *Config {
	gJConfig = jconfig.New(defConfigDir, "config.json", Config{})

	if _, err := gJConfig.Load(DefaultConfig); err != nil {
		clog.Fatal("load config error:", err)
	}

	config := gJConfig.Data().(*Config)
	config.filepath = gJConfig.FilePath()
	FixConfig(config)

	if dump {
		fmt.Println(config)
	}
	return config
}

func saveConfig() {
	if gJConfig == nil {
		clog.Warn("gJConfig is nil")
		return
	}
	gJConfig.Save()
}
