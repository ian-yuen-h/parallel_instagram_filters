// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import "image/color"

func (img *ImageTask) getKernel(s string) [9]float64{
	var kernel [9]float64
	switch s{
	case "S":
		kernel = [9]float64{0.0, -1.0, 0.0, -1.0, 5.0, -1.0, 0.0, -1.0, 0.0}
	case "E":
		kernel = [9]float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}
	case "B":
		kernel = [9]float64{1/9.0, 1 / 9, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
	case "G":
		kernel = [9]float64{}
	}
	return kernel
}


func (img *ImageTask) MConvulsionControl(grey int, y1 int, y2 int, kernel [9]float64){
	if grey != 0{
		for y := y1; y <= y2; y++ {
			for x := 0; x <= img.MaxX; x++ {
				r, g, b, a := img.Temp.At(x, y).RGBA()
				greyC := clamp(float64(r+g+b) / 3)
				img.Out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
			}
		}
	}else{
		for y := y1; y <= y2; y++ {
			for x := 0 ; x <= img.MaxX; x++ {
				neighbors := img.ExtractNeighbors(x, y, img.MaxX, img.MaxY)
				newResults := img.ApplyConvulusion(neighbors, kernel)
				img.Out.Set(x, y, color.RGBA64{clamp(newResults.Red), clamp(newResults.Green), clamp(newResults.Blue), clamp(newResults.Alpha)})

			}
		}
	}
}

func (img *ImageTask) SeqFilterControl(filter string) {

	grey := 0
	kernel := img.getKernel(filter)
	if kernel == [9]float64{}{
		grey =1
	}
	img.MConvulsionControl(grey, 0, img.MaxY, kernel)
	temp := img.Temp
	img.Temp = img.Out	//swaps pointer
	img.Out = temp
}


//neighbor struct
type Neighbors struct {
	inputs [9]*Colors
}

//color value struct
type Colors struct{
	Red float64
	Green float64
	Blue float64
	Alpha float64
}

func (img *ImageTask) ExtractNeighbors(xBase int, yBase int, maxX int, maxY int) *Neighbors {

	var counter int
	counter = 0

	var zeNeighbors Neighbors

	for yr := -1; yr <= 1; yr++ {
		for xr := -1; xr <= 1; xr++ {
			currentX := xBase + xr
			currentY:= yBase + yr
			if (currentX < 0) || (currentY <0) || (currentX > maxX) || (currentY > maxY){
				toAdd := Colors{Red: 0.0, Green: 0.0, Blue: 0.0, Alpha: 0.0}
				zeNeighbors.inputs[counter] = &toAdd
				counter += 1
			}else{
				r, g, b, a := img.Temp.At(currentX, currentY).RGBA()

				toAdd := Colors{Red: float64(r), Green: float64(g), Blue: float64(b), Alpha: float64(a)}
				zeNeighbors.inputs[counter] = &toAdd
				counter += 1
			}
		}
	}
	return &zeNeighbors
}

//given neighbors and kernel, apply convulsions
func (img *ImageTask) ApplyConvulusion(neighbors *Neighbors, kernel [9]float64) *Colors{

	var blue float64
	var green float64
	var red float64
	for r := 0; r < 9; r++ {
		red = red + (neighbors.inputs[r].Red * kernel[r])
		green = green + (neighbors.inputs[r].Green * kernel[r])
		blue = blue + (neighbors.inputs[r].Blue * kernel[r])
	}

	returned := Colors{Red: red, Green: green, Blue: blue, Alpha: neighbors.inputs[4].Alpha}
	return &returned
}



