// Package png allows for loading png images and applying
// image flitering effects on them
package png

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// The Image represents a structure for working with PNG images.
type ImageTask struct {
	Temp *image.RGBA64
	Out *image.RGBA64
	Bounds  image.Rectangle
	MaxX int
	MaxY int
}

type Effects struct{
	InPath string
	OutPath string
	Effects []string
	ImgTask	*ImageTask
	Threads int
}

//
// Public functions
//

// Load returns a Image that was loaded based on the filePath parameter
func Load(filePath string) (*ImageTask, error) {

	inReader, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer inReader.Close()

	inImg, err := png.Decode(inReader)

	if err != nil {
		return nil, err
	}

	inBounds := inImg.Bounds()

	temp := image.NewRGBA64(inBounds)

	outImg := image.NewRGBA64(inBounds)

	bounds := inBounds
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := inImg.At(x, y).RGBA()
			temp.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}

	return &ImageTask{Temp: temp,Out: outImg, Bounds: inBounds, MaxX: inBounds.Max.X, MaxY: inBounds.Max.Y}, nil
}

// Save saves the image to the given file
func (img *ImageTask) Save(filePath string) error {

	outWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outWriter.Close()

	err = png.Encode(outWriter, img.Temp)
	if err != nil {
		return err
	}
	return nil
}

func (img *ImageTask) Save2(filePath string) error {

	outWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outWriter.Close()

	err = png.Encode(outWriter, img.Out)
	if err != nil {
		return err
	}
	return nil
}


//clamp will clamp the comp parameter to zero if it is less than zero or to 65535 if the comp parameter
// is greater than 65535.
func clamp(comp float64) uint16 {
	return uint16(math.Min(65535, math.Max(0, comp)))
}

