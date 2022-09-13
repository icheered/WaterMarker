package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("************************************************************")
	fmt.Println("************************************************************")
	fmt.Println("WaterMarker v1.0 - Written by Tjeerd Bakker (ICheered) in Go")
	fmt.Println("************************************************************")
	fmt.Println("************************************************************")
	fmt.Println("")
	fmt.Println("For help: run the program from command line with the -h flag")
	fmt.Println("Having issues? Please let me know at Tjeerd992@gmail.com")
	fmt.Println("")
	watermarkOpacity, watermarkLocation, watermarkFile, sourceDir, targetDir, force := getParameters()

	fmt.Println("Using following parameters:")
	fmt.Printf("- Opacity:          %d\n", watermarkOpacity)
	fmt.Printf("- Location:         %s\n", watermarkLocation)
	fmt.Printf("- Watermark:        %s\n", watermarkFile)
	fmt.Printf("- Source directory: %s\n", sourceDir)
	fmt.Printf("- Target directory: %s\n", targetDir)
	fmt.Println("")

	if _, err := os.Stat(watermarkFile); errors.Is(err, os.ErrNotExist) {
		// Watermark file does not exist
		fmt.Printf("ERROR: Watermark file '%s' does not exist in this directory\n", watermarkFile)
		os.Exit(1)
	}

	if !strings.HasSuffix(watermarkFile, ".png") {
		fmt.Printf("ERROR: Watermark file '%s' is not a PNG file\n", watermarkFile)
		os.Exit(1)
	}

	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		// Source folder does not exist
		fmt.Printf("ERROR: Source folder (folder containing images) '%s' does not exist in this directory\n", sourceDir)
		os.Exit(1)
	}

	if _, err := os.Stat(targetDir); err == nil {
		// Target dir already exists
		fmt.Printf("WARNING: Target folder '%s' already exists in this directory. \n", targetDir)
		if force {
			fmt.Println("         Using --force, so will overwrite existing files")
		} else {
			fmt.Println("         Use --force to overwrite existing files")
			fmt.Println("         Exiting to avoid overwriting existing files.")
			os.Exit(1)
		}
	} else {
		os.Mkdir(targetDir, 0755)
	}
	fmt.Print("\n--------------------------------------\n")

	watermark := openImage(watermarkFile, "png")
	mask := image.NewUniform(color.Alpha{uint8(watermarkOpacity * 255)})
	files := getFiles(sourceDir)

	fmt.Printf("Starting: Processing %d files\n\n", len(files))

	var wg sync.WaitGroup
	wg.Add(len(files))
	start := time.Now()
	for _, file := range files {
		go func(file os.FileInfo, watermark image.Image, mask image.Image, watermarkLocation string, sourceDir string, targetDir string) {
			defer wg.Done()
			if !(strings.HasSuffix(file.Name(), ".jpg")) && !(strings.HasSuffix(file.Name(), ".jpeg")) {
				fmt.Printf("Skipping photo '%s' because it is not a .jpg or .jpeg\n", file)
				return
			}

			srcImage := openImage(path.Join(sourceDir, file.Name()), "jpeg")

			imgSize := srcImage.Bounds()
			wmSize := watermark.Bounds()
			canvas := image.NewRGBA(imgSize)
			var watermarkOffset image.Point
			if watermarkLocation == "left" {
				watermarkOffset = image.Point{0, imgSize.Max.Y - wmSize.Max.Y}
			} else if watermarkLocation == "right" {
				watermarkOffset = image.Point{imgSize.Max.X - wmSize.Max.X, imgSize.Max.Y - wmSize.Max.Y}
			}

			draw.Draw(canvas, imgSize, srcImage, image.Point{0, 0}, draw.Src)
			draw.DrawMask(canvas, imgSize.Add(watermarkOffset), watermark, image.Point{0, 0}, mask, image.Point{0, 0}, draw.Over)

			saveImage(canvas, targetDir, file.Name())
		}(file, watermark, mask, watermarkLocation, sourceDir, targetDir)
	}
	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("\nAll done! Editted %d files in %s\n", len(files), elapsed)

}

func getFiles(dir string) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func getParameters() (int, string, string, string, string, bool) {
	paramOpacity := flag.Int("opacity", 60, "Watermark opacity between 0 and 100")
	paramLocation := flag.String("location", "right", "Location of watermark [left, right]")
	paramWatermark := flag.String("watermark", "watermark.png", "Name of PNG image to be used as watermark")
	paramSourceDir := flag.String("source", "photos", "Source directory (location to find un-watermarked photos)")
	paramTargetDir := flag.String("target", "watermarked", "Target directory (location to put watermarked photos")
	paramForce := flag.Bool("force", false, "Force overwrite of target directory if it already exists")

	flag.Parse()
	opacity := *paramOpacity
	location := *paramLocation
	watermark := *paramWatermark
	sourceDir := *paramSourceDir
	targetDir := *paramTargetDir
	force := *paramForce

	return opacity, location, watermark, sourceDir, targetDir, force
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

func openImage(fname string, ftype string) image.Image {
	inputfile, err := os.Open(fname)
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}

	var srcimage image.Image

	if ftype == "jpeg" {
		srcimage, err = jpeg.Decode(inputfile)
		if err != nil {
			log.Fatalf("failed to decode: %s", err)
		}
		defer inputfile.Close()
	} else if ftype == "png" {
		srcimage, err = png.Decode(inputfile)
		if err != nil {
			log.Fatalf("failed to decode: %s", err)
		}
		defer inputfile.Close()
	}

	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer inputfile.Close()
	return srcimage
}
