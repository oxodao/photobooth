package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/oxodao/photobooth/cmd"
	"github.com/oxodao/photobooth/config"
	"github.com/oxodao/photobooth/logs"
	"github.com/oxodao/photobooth/migrations"
	"github.com/oxodao/photobooth/services"
	"github.com/oxodao/photobooth/utils"
)

//go:embed gui/dist
var webapp embed.FS

//go:embed gui_admin/dist
var adminapp embed.FS

//go:embed sql
var dbScripts embed.FS

func main() {
	if err := config.Load(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err := utils.MakeOrCreateFolder("")
	if err != nil {
		fmt.Println("Failed to create root folder")
		os.Exit(1)
	}

	logs.Init()

	//#region Weird hack to embed webapp only on build + Instanciating the logger
	execPath := os.Args[0]
	var webappFs *fs.FS = nil
	var adminappFs *fs.FS = nil

	if !strings.HasPrefix(execPath, "/tmp/") {
		subfs, err := fs.Sub(webapp, "gui/dist")
		if err != nil {
			logs.Warn("Failed to get webapp path. Not loading the webapp", err)
		} else {
			webappFs = &subfs
		}

		subfs2, err := fs.Sub(adminapp, "gui_admin/dist")
		if err != nil {
			logs.Warn("Failed to get adminapp path. Not loading the adminapp", err)
		} else {
			adminappFs = &subfs2
		}
	}
	//#endregion

	if err := migrations.CheckDbExists(dbScripts); err != nil {
		logs.Error(err)
		os.Exit(1)
	}

	if err := services.Load(webappFs, adminappFs); err != nil {
		logs.Error(err)
		os.Exit(1)
	}

	cmd.Execute()
}
