package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"log"
	"os"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	logFile, err := os.OpenFile("randOpen.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		log.Println("日志文件打开失败：", err)
		return
	}
	log.SetOutput(logFile)
	app := &App{}
	err = wails.Run(&options.App{
		Title:            "randOpen",
		Width:            300,
		Height:           260,
		AssetServer:      &assetserver.Options{Assets: assets},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255}, // 白色背景
		OnStartup:        app.startup,
		Bind:             []interface{}{app},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
