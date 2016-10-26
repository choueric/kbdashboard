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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/choueric/clog"
)

const (
	ConfigDir = ".config/kbdashboard"
)

const DefaultConfig = `
{
	"editor": "vim",
	"current": 1,
	"profile": [
	{
		"name":"demo",
		"src_dir":"/home/user/kernel"
		"arch":"arm",
		"target":"uImage",
		"defconfig":"at91rm9200_defconfig",
		"cross_compile":"arm-eabi-",
		"output_dir":"./_build",
		"mod_install_dir":"./_build/mod",
		"thread_num":4,
	}
	]
}
`

type Profile struct {
	Name          string `json:"name"`
	SrcDir        string `json:"src_dir"`
	Arch          string `json:"arch"`
	Target        string `json:"target"`
	Defconfig     string `json:"defconfig"`
	CrossComile   string `json:"cross_compile"`
	OutputDir     string `json:"output_dir"`
	ModInstallDir string `json:"mod_install_dir"`
	ThreadNum     int    `json:"thread_num"`
}

type Config struct {
	Editor     string     `json:"editor"`
	Current    int        `json:"current"`
	Profiles   []*Profile `json:"profile"`
	configFile string
}

func (p *Profile) String() string {
	return fmt.Sprintf(
		"name = %s\n  arch = %s, CC = %s, target = %s, defconfig = %s\n"+
			"  src_dir = %s\n  build_dir = %s, mod_dir = %s\n  thread num = %d\n",
		p.Name, p.Arch, p.CrossComile, p.Target, p.Defconfig, p.SrcDir,
		p.OutputDir, p.ModInstallDir, p.ThreadNum)
}

func checkConfigDir(p string) {
	homeDir := os.Getenv("HOME")
	err := os.MkdirAll(homeDir+"/"+p, os.ModeDir|0777)
	if err != nil {
		clog.Println("mkdir:", err)
	}
}

func checkConfigFile(p string) string {
	if p == "" {
		p = os.Getenv("HOME") + "/" + ConfigDir + "/config.json"
	}
	_, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		clog.Println("create an empty config file.")
		file, err := os.Create(p)
		_, err = file.Write([]byte(DefaultConfig))
		if err != nil {
			clog.Fatal(err)
		}
		file.Close()
	} else if err != nil {
		clog.Fatal(err)
	}

	return p
}

func getInstallFilename(p *Profile) string {
	return os.Getenv("HOME") + "/" + ConfigDir + "/" + p.Name + "_install.sh"
}

func checkFileExsit(p string) bool {
	_, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		return false
	} else if err != nil {
		clog.Fatal(err)
	}

	return true
}

// if OutputDir and ModInstallDir is relative, change it to absolute
// by adding SrcDir prefix.
func fixRelativeDir(p string, pre string) string {
	if !path.IsAbs(p) {
		p = path.Join(pre, p)
	}
	return p
}

func ParseConfig(p string) (*Config, error) {
	checkConfigDir(ConfigDir)
	p = checkConfigFile(p)

	file, err := os.Open(p)
	if err != nil {
		clog.Println(err)
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	c.configFile = p
	if err = json.Unmarshal(data, c); err != nil {
		return nil, err
	}

	if c.Current >= len(c.Profiles) {
		clog.Fatal("Current in config.json is invalid: ", c.Current)
	}

	for _, p := range c.Profiles {
		p.OutputDir = fixRelativeDir(p.OutputDir, p.SrcDir)
		p.ModInstallDir = fixRelativeDir(p.ModInstallDir, p.SrcDir)
		if p.Defconfig == "" {
			p.Defconfig = "defconfig"
		}
	}

	return c, nil
}

func writeConfigFile(config *Config) {
	data, err := json.MarshalIndent(config, "  ", "  ")
	if err != nil {
		clog.Fatal(err)
	}

	file, err := os.Create(config.configFile)
	if err != nil {
		clog.Fatal(err)
	}
	defer file.Close()

	file.Write(data)
}
