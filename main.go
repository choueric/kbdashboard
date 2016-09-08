package main

import (
	"fmt"
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

	return config
}

func handleCmd(cmd string, args []string, config *Config) {
	h, ok := handerMap[cmd]
	if !ok {
		printUsage()
		log.Fatalf("[%s] is not supported\n", cmd)
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
		printUsage()
		return
	}

	config := parseConfig()
	handleCmd(cmd, args, config)
}

func printUsage() {
	fmt.Printf("Usage: \n")
	for k, _ := range handerMap {
		fmt.Printf("  - %s\n", k)
	}
}
