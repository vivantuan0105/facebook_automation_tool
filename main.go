package main

import (
	"embed"
	"log"

	"socialmanager/backend/app"
	"socialmanager/backend/database"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	database.InitDB()

	application := app.NewApp()

	err := wails.Run(&options.App{
		Title:  "Social Desktop Manager",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 24, G: 24, B: 27, A: 1},
		OnStartup:        application.Startup,
		Bind: []interface{}{
			application,
		},
	})

	if err != nil {
		log.Fatal("Error starting app:", err)
	}
}
