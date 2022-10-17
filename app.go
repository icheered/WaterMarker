package main

import (
	"context"

	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/nfnt/resize"
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

func (a *App) SelectFile() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Images (*.png;*.jpeg)",
				Pattern:     "*.png;*.jpg;*.jpeg",
			},
		},
	})
	//fmt.Print(selection)
	if err != nil {
		return ""
	}

	return selection
}

func (a *App) SelectFolder() string {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Select Folder",
		Filters: []runtime.FileFilter{},
	})
	if err != nil {
		return "error"
	}
	return selection
}

func (a *App) GetNumberOfFiles(sourcefolderPath string) int {
	files := getFiles(sourcefolderPath)
	return len(files)
}

func (a *App) FetchPreview(watermarkPath, sourcefolderPath, targetfolderPath, watermarkPosition string, watermarkOpacity, watermarkScale float64) returnStruct {
	watermark, err := openImage(watermarkPath)
	if err != nil {
		return returnStruct{Status: "error", Message: err.Error()}
	}

	mask := image.NewUniform(color.Alpha{uint8(watermarkOpacity * 255)})
	files := getFiles(sourcefolderPath)
	// Remove all but first file
	file := files[0]
	watermarkFile(file, watermark, mask, watermarkPosition, watermarkScale, sourcefolderPath, targetfolderPath)

	fmt.Print("Done")

	fpath := path.Join(targetfolderPath, file.Name())
	return returnStruct{Status: "success", Message: fpath}
}

func (a *App) ProcessImages(watermarkPath, sourcefolderPath, targetfolderPath, watermarkPosition string, watermarkOpacity, watermarkScale float64) returnStruct {

	watermark, err := openImage(watermarkPath)
	if err != nil {
		return returnStruct{Status: "error", Message: err.Error()}
	}

	mask := image.NewUniform(color.Alpha{uint8(watermarkOpacity * 255)})
	files := getFiles(sourcefolderPath)

	fmt.Printf("Starting: Processing %d files\n\n", len(files))

	var wg sync.WaitGroup
	wg.Add(len(files))
	start := time.Now()

	for _, file := range files {
		go func(loopFile os.FileInfo) {
			watermarkFile(loopFile, watermark, mask, watermarkPosition, watermarkScale, sourcefolderPath, targetfolderPath)
			defer wg.Done()
		}(file)
	}
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Print("Done: ", elapsed)

	return returnStruct{Status: "success", Message: fmt.Sprintf("\nAll done! Editted %d files in %s", len(files), elapsed)}
}

func watermarkFile(file os.FileInfo, watermark image.Image, mask image.Image, watermarkPosition string, watermarkScale float64, sourcefolderPath string, targetfolderPath string) {
	if !(strings.HasSuffix(file.Name(), ".jpg")) && !(strings.HasSuffix(file.Name(), ".jpeg")) {
		fmt.Printf("Skipping photo '%s' because it is not a .jpg or .jpeg\n", file.Name())
		return
	}
	fmt.Printf("Processing photo '%s'\n", file.Name())
	srcImage, err := openImage(path.Join(sourcefolderPath, file.Name()))
	if err != nil {
		return
	}
	imgSize := srcImage.Bounds()

	imgheight := min(imgSize.Dx(), imgSize.Dy())
	scaledWatermark := resize.Resize(0, uint(watermarkScale*float64(imgheight)), watermark, resize.NearestNeighbor)

	wmSize := scaledWatermark.Bounds()
	canvas := image.NewRGBA(imgSize)
	var watermarkOffset image.Point
	if watermarkPosition == "bottom-left" {
		watermarkOffset = image.Point{0, imgSize.Max.Y - wmSize.Max.Y}
	} else if watermarkPosition == "bottom-right" {
		watermarkOffset = image.Point{imgSize.Max.X - wmSize.Max.X, imgSize.Max.Y - wmSize.Max.Y}
	} else if watermarkPosition == "top-left" {
		watermarkOffset = image.Point{0, 0}
	} else if watermarkPosition == "top-right" {
		watermarkOffset = image.Point{imgSize.Max.X - wmSize.Max.X, 0}
	} else if watermarkPosition == "center" {
		watermarkOffset = image.Point{(imgSize.Max.X - wmSize.Max.X) / 2, (imgSize.Max.Y - wmSize.Max.Y) / 2}
	}

	draw.Draw(canvas, imgSize, srcImage, image.Point{0, 0}, draw.Src)
	draw.DrawMask(canvas, imgSize.Add(watermarkOffset), scaledWatermark, image.Point{0, 0}, mask, image.Point{0, 0}, draw.Over)
	saveImage(canvas, targetfolderPath, file.Name())
	fmt.Printf("Finished processing photo '%s'\n", file.Name())
}

func getFiles(dirname string) []os.FileInfo {
	entries, err := os.ReadDir(dirname)
	if err != nil {

		log.Fatal(err)
	}
	infos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Fatal(err)
		}
		infos = append(infos, info)
	}
	return infos
}

func saveImage(img image.Image, pname, fname string) error {
	fpath := path.Join(pname, fname)
	outputFile, err := os.Create(fpath)

	if err != nil {
		log.Fatalf("failed to create: %s", err)
	}
	var opt jpeg.Options
	opt.Quality = 95

	jpeg.Encode(outputFile, img, &opt)
	defer outputFile.Close()
	return nil
}

func openImage(fname string) (image.Image, error) {
	inputfile, err := os.Open(fname)
	if err != nil {
		fmt.Print("Failed to open: " + fname)
		return image.NewUniform(color.Black), errors.New("Failed to open: " + fname)
	}

	var srcimage image.Image
	if fname[len(fname)-4:] == ".png" {
		srcimage, err = png.Decode(inputfile)
		if err != nil {
			return image.NewUniform(color.Black), errors.New("Failed to decode: " + fname)
		}
		defer inputfile.Close()
	} else if fname[len(fname)-4:] == ".jpg" || fname[len(fname)-5:] == ".jpeg" {
		srcimage, err = jpeg.Decode(inputfile)
		if err != nil {
			return image.NewUniform(color.Black), errors.New("Failed to decode: " + fname)
		}
		defer inputfile.Close()
	} else {
		return image.NewUniform(color.Black), errors.New("Failed to open: " + fname + ". Not PNG/JPG/JPEG")
	}
	return srcimage, nil
}

func min(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}
