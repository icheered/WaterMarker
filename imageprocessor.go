package main

import (
	"fmt"
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"
	"math"
	"io/ioutil"

	"github.com/nfnt/resize"
)


func watermarkFile(file os.FileInfo, watermark image.Image, mask image.Image, watermarkPosition string, watermarkScale float64, sourceFolderPath string, targetFolderPath string) error {
    if !(strings.HasSuffix(file.Name(), ".jpg")) && !(strings.HasSuffix(file.Name(), ".jpeg")) {
        return fmt.Errorf("skipping photo '%s' because it is not a .jpg or .jpeg", file.Name())
    }
    fmt.Printf("Processing photo '%s'\n", file.Name())

    fileHandle, srcImage, err := openImage(path.Join(sourceFolderPath, file.Name()))
    if err != nil {
        return fmt.Errorf("failed to open image: %w", err)
    }
    defer fileHandle.Close()

    imgSize := srcImage.Bounds()

    imgHeight := int(math.Min(float64(imgSize.Dx()), float64(imgSize.Dy())))
    scaledWatermark := resize.Resize(0, uint(watermarkScale*float64(imgHeight)), watermark, resize.NearestNeighbor)

    wmSize := scaledWatermark.Bounds()
    canvas := image.NewRGBA(imgSize)
    var watermarkOffset image.Point

    switch watermarkPosition {
    case "bottom-left":
        watermarkOffset = image.Point{0, imgSize.Max.Y - wmSize.Max.Y}
    case "bottom-right":
        watermarkOffset = image.Point{imgSize.Max.X - wmSize.Max.X, imgSize.Max.Y - wmSize.Max.Y}
    case "top-left":
        watermarkOffset = image.Point{0, 0}
    case "top-right":
        watermarkOffset = image.Point{imgSize.Max.X - wmSize.Max.X, 0}
    case "center":
        watermarkOffset = image.Point{(imgSize.Max.X - wmSize.Max.X) / 2, (imgSize.Max.Y - wmSize.Max.Y) / 2}
    default:
        return fmt.Errorf("invalid watermark position: %s", watermarkPosition)
    }

    draw.Draw(canvas, imgSize, srcImage, image.Point{0, 0}, draw.Src)
    draw.DrawMask(canvas, imgSize.Add(watermarkOffset), scaledWatermark, image.Point{0, 0}, mask, image.Point{0, 0}, draw.Over)

	err = saveImage(canvas, targetFolderPath, file.Name())
    if err != nil {
        return fmt.Errorf("failed to save image: %w", err)
    }

    fmt.Printf("Finished processing photo '%s'\n", file.Name())
    return nil
}


func openImage(fname string) (*os.File, image.Image, error) {
	inputfile, err := os.Open(fname)
	if err != nil {
		return nil, image.NewUniform(color.Black), fmt.Errorf("failed to open: %s", err)
	}

	var img image.Image
	switch path.Ext(fname) {
	case ".png":
		img, err = png.Decode(inputfile)
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(inputfile)
	default:
		inputfile.Close()
		return nil, image.NewUniform(color.Black), fmt.Errorf("file is not PNG/JPG/JPEG: %s", fname)
	}

	if err != nil {
		inputfile.Close()
		return nil, image.NewUniform(color.Black), fmt.Errorf("failed to decode: %w", err)
	}

	return inputfile, img, nil
}

func saveImage(img image.Image, pname, fname string) error {
    fpath := path.Join(pname, fname)

    buf := new(bytes.Buffer)
    var opt jpeg.Options
    opt.Quality = 90

    if err := jpeg.Encode(buf, img, &opt); err != nil {
        return fmt.Errorf("failed to encode jpeg: %s", err)
    }

    if err := ioutil.WriteFile(fpath, buf.Bytes(), 0644); err != nil {
        return fmt.Errorf("failed to write file: %s", err)
    }

    return nil
}
