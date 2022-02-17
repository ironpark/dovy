package main

import (
	"embed"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	err := wails.Run(&options.App{
		Title:             "dovey",
		Width:             1024,
		Height:            768,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         true,
		MinWidth:          0,
		MinHeight:         0,
		MaxWidth:          0,
		MaxHeight:         0,
		StartHidden:       false,
		HideWindowOnClose: false,
		AlwaysOnTop:       false,
		RGBA:              &options.RGBA{255, 255, 255, 255},
		Assets:            assets,
		Menu:              nil,
		Logger:            nil,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnShutdown:        app.shutdown,
		OnBeforeClose:     nil,
		Bind: []interface{}{
			app,
		},
		WindowStartState: 0,
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		Mac: &mac.Options{
			TitleBar:             nil,
			Appearance:           "",
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			About:                nil,
		},
		Linux: nil,
	})

	if err != nil {
		log.Fatal(err)
	}
}
