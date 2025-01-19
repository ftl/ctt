package main

import (
	"github.com/ftl/ctt/pkg/app"
	"github.com/ftl/ctt/pkg/fyneui"
)

const (
	appID   = "ft.ctt"
	version = "develop"
)

func main() {
	app := app.NewApp()

	ui := fyneui.NewApp(appID, version, app)
	app.SetReport(ui)
	app.Quitter = ui
	app.UI = ui

	ui.Run()
}
