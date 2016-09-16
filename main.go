package main

import (
	"log"
	"os"
)

const (
	CRED   = "\x1b[31;1m"
	CGREEN = "\x1b[32;1m"
	CEND   = "\x1b[0;m"
)

func parseConfig() *Config {
	config, err := ParseConfig("")
	if err != nil {
		log.Fatal(err)
	}
	if config == nil {
		log.Fatal("config is nil.")
	}

	//log.Println(config)
	return config
}

func handleCmd(cmd string, args []string, config *Config) {
	if cmd == "help" {
		cmd_help(nil, nil)
		return
	}

	h, ok := handerMap[cmd]
	if !ok {
		log.Fatalf("[%s] is not supported\n", cmd)
		cmd_help(nil, nil)
	}
	h(args, config)
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var (
		args []string
		cmd  string
	)

	switch len(os.Args) {
	case 2:
		cmd = os.Args[1]
	case 3:
		cmd = os.Args[1]
		args = os.Args[2:len(os.Args)]
	default:
		cmd_help(nil, nil)
		return
	}

	config := parseConfig()
	handleCmd(cmd, args, config)
}
