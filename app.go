package main

import (
	"path"
	"fmt"
	"context"
	"image"
	"image/color"
	"sync"
	"time"
	"os"
	"errors"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type returnStruct struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (a *App) SelectFile() (string, error) {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Images (*.png;*.jpeg)",
				Pattern:     "*.png;*.jpg;*.jpeg",
			},
		},
	})
	if err != nil {
		return "", err
	}

	return selection, nil
}

func (a *App) SelectFolder() (string, error) {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Select Folder",
		Filters: []runtime.FileFilter{},
	})
	if err != nil {
		return "", err
	}
	return selection, nil
}

func (a *App) GetNumberOfFiles(sourceFolderPath string) (int, error) {
	files, err := getFiles(sourceFolderPath)
	if err != nil {
		return 0, err
	}
	return len(files), nil
}

func (a *App) FetchPreview(watermarkPath, sourceFolderPath, targetFolderPath, watermarkPosition string, watermarkOpacity, watermarkScale float64) (returnStruct, error) {
	watermarkFileDescriptor, watermark, err := openImage(watermarkPath)
	if err != nil {
		return returnStruct{}, err
	}
	defer watermarkFileDescriptor.Close()

	mask := image.NewUniform(color.Alpha{uint8(watermarkOpacity * 255)})
	files, err := getFiles(sourceFolderPath)
	if err != nil {
		return returnStruct{}, err
	}
	if len(files) == 0 {
		return returnStruct{}, errors.New("no files found in source folder")
	}

	// Remove all but first file
	file := files[0]
	if err := watermarkFile(file, watermark, mask, watermarkPosition, watermarkScale, sourceFolderPath, targetFolderPath); err != nil {
		log.Printf("Failed to watermark file %s: %v", file.Name(), err)
	}

	fmt.Println("Finished watermarking file: ", file.Name())

	fpath := path.Join(targetFolderPath, file.Name())
	return returnStruct{Status: "success", Message: fpath}, nil
}




func (a *App) ProcessImages(watermarkPath, sourceFolderPath, targetFolderPath, watermarkPosition string, watermarkOpacity, watermarkScale float64) (string, error) {
    watermarkFileDescriptor, watermark, err := openImage(watermarkPath)
    if err != nil {
        return "", err
    }
    defer watermarkFileDescriptor.Close()

    mask := image.NewUniform(color.Alpha{uint8(watermarkOpacity * 255)})
    files, err := getFiles(sourceFolderPath)
    if err != nil {
        return "", err
    }

    fmt.Printf("Starting: Processing %d files\n\n", len(files))

    var wg sync.WaitGroup
    start := time.Now()

    MAX_PARALLEL_PROCESSES := 8
    sem := make(chan struct{}, MAX_PARALLEL_PROCESSES)

    for _, file := range files {
        wg.Add(1)

        go func(loopFile os.FileInfo) {
            defer wg.Done()

            sem <- struct{}{}
            defer func() { <-sem }()

            if err := watermarkFile(loopFile, watermark, mask, watermarkPosition, watermarkScale, sourceFolderPath, targetFolderPath); err != nil {
                log.Printf("Failed to watermark file %s: %v", loopFile.Name(), err)
            }

        }(file)
    }
    wg.Wait()

    elapsed := time.Since(start)
    fmt.Printf("Done: Processed %d files in %s\n", len(files), elapsed)

    return fmt.Sprintf("\nAll done! Edited %d files in %s", len(files), elapsed), nil
}


