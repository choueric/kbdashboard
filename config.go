package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	return fmt.Sprintf("name = %s\n  arch = %s, CC = %s, target = %s\n"+
		"  src_dir = %s\n  build_dir = %s, mod_dir = %s\n  thread num = %d\n",
		p.Name, p.Arch, p.CrossComile, p.Target, p.SrcDir, p.OutputDir, p.ModInstallDir,
		p.ThreadNum)
}

func checkConfigDir(path string) {
	homeDir := os.Getenv("HOME")
	err := os.MkdirAll(homeDir+"/"+path, os.ModeDir|0777)
	if err != nil {
		log.Println("mkdir:", err)
	}
}

func checkConfigFile(path string) string {
	if path == "" {
		path = os.Getenv("HOME") + "/" + ConfigDir + "/config.json"
	}
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		log.Println("create an empty config file.")
		file, err := os.Create(path)
		_, err = file.Write([]byte(DefaultConfig))
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	} else if err != nil {
		log.Fatal(err)
	}

	return path
}

func getInstallFilename(p *Profile) string {
	return os.Getenv("HOME") + "/" + ConfigDir + "/" + p.Name + "_install.sh"
}

func checkFileExsit(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	} else if err != nil {
		log.Fatal(err)
	}

	return true
}

func ParseConfig(path string) (*Config, error) {
	checkConfigDir(ConfigDir)
	path = checkConfigFile(path)

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	config.configFile = path
	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if config.Current >= len(config.Profiles) {
		log.Fatal("Current in config.json is invalid: ", config.Current)
	}

	return config, nil
}

func writeConfigFile(config *Config) {
	data, err := json.MarshalIndent(config, "  ", "  ")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(config.configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Write(data)
}
