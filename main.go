package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"context"


	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

type FileLoader struct {
	http.Handler
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	var requestedFilename string = ""
	if runtime.GOOS == "windows" {
		requestedFilename = strings.TrimPrefix(req.URL.Path, "/")
	} else {
		requestedFilename = req.URL.Path
	}
	println("Requesting file:", requestedFilename)
	fileData, err := os.ReadFile(requestedFilename)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFilename)))
	}

	res.Write(fileData)
}

func main() {
	// Create an instance of the app structure
	app := NewApp()
	app.startup(context.Background())

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "WaterMarker",
		Width:            1024,
		Height:           700,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		AssetsHandler:    NewFileLoader(),
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
