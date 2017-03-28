// kbdashboard is used to manage building processes of multiple linux kernels.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/choueric/cmdmux"
)

var logger *log.Logger

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "dump" {
		getConfig(true)
		return
	}
	gConfig = getConfig(false)

	if gConfig.Debug {
		logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logger = log.New(ioutil.Discard, "", 0)
	}

	cmdmux.HandleFunc("/", helpHandler)
	cmdmux.HandleFunc("/help", helpHandler)
	cmdmux.HandleFunc("/list", listHandler)
	cmdmux.HandleFunc("/choose", chooseHandler)

	cmdmux.HandleFunc("/edit", editProfileHandler)
	cmdmux.HandleFunc("/edit/profile", editProfileHandler)
	cmdmux.HandleFunc("/edit/install", editInstallHandler)

	cmdmux.HandleFunc("/config", configMenuHandler)
	cmdmux.HandleFunc("/config/menu", configMenuHandler)
	cmdmux.HandleFunc("/config/def", configDefHandler)
	cmdmux.HandleFunc("/config/save", configSaveHandler)

	cmdmux.HandleFunc("/build", buildImageHandler)
	cmdmux.HandleFunc("/build/image", buildImageHandler)
	cmdmux.HandleFunc("/build/modules", buildModulesHandler)
	cmdmux.HandleFunc("/build/dtb", buildDtbHandler)

	cmdmux.HandleFunc("/install", installHandler)
	cmdmux.HandleFunc("/make", makeHandler)

	cmdmux.HandleFunc("/dts", dtsListHandler)
	cmdmux.HandleFunc("/dts/list", dtsListHandler)
	cmdmux.HandleFunc("/dts/link", dtsLinkHandler)

	cmdmux.HandleFunc("/version", versionHandler)
	cmdmux.HandleFunc("/completion", completionHandler)

	ret, err := cmdmux.Execute(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Execute Error: %v\n", err)
		os.Exit(-1)
	}
	os.Exit(ret)
}
