package main

import (
	"embed"

	"github.com/SamHennessy/hlive-wails-hn/pages"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	ps := NewPageServerWails(pages.Index)

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "HLive-Wails Hacker News Reader",
		Width:            1024,
		Height:           768,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        ps.OnStartup,
		AssetsHandler:    ps.AssetsHandler(),
		LogLevel:         logger.DEBUG,
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
