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
	"os"

	"github.com/choueric/clog"
	"github.com/choueric/cmdmux"
)

func main() {
	clog.SetFlags(clog.Lshortfile | clog.LstdFlags | clog.Lcolor)

	if len(os.Args) >= 2 && os.Args[1] == "dump" {
		getConfig(true)
		return
	}

	config := getConfig(false)

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
	cmdmux.HandleFunc("/version", versionHandler)

	ret, err := cmdmux.Execute(config)
	if err != nil {
		clog.Warn("Execute error:", err)
		os.Exit(0)
	}
	os.Exit(ret)
}
