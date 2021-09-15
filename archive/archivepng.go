// Package png allows for loading png images and applying
// image flitering effects on them
package archive

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

// The Image represents a structure for working with PNG images.
type Image struct {
	in  image.Image
	temp *image.RGBA64
	out *image.RGBA64
}

//
// Public functions
//

// Load returns a Image that was loaded based on the filePath parameter
func Load(filePath string) (*Image, error) {

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

	outImg := image.NewRGBA64(inBounds)

	return &Image{inImg, outImg}, nil
}

// Save saves the image to the given file
func (img *Image) Save(filePath string) error {

	outWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outWriter.Close()

	err = png.Encode(outWriter, img.out)
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


type ImageTask struct {
	in  *image.RGBA64
	out *image.RGBA64
	Bounds image.Rectangle
}




func (img *ImageTask) ConvulsionControl2(kernel [9]float64){
	//matrice := {{0,-1,0},{-1,5,-1},{0,-1,0}}
	matrice := [][]int{{0,-1,0},{-1,5,-1},{0,-1,0}}
	bounds := img.Out.Bounds()
	sumR := 0
	sumB := 0
	sumG := 0
	var r uint32
	var g uint32
	var b uint32
	var a uint32
	//maxX := bounds.Max.X
	//maxY := bounds.Max.Y
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					var imageX int
					var imageY int
					imageX = x + i
					imageY = y + j

					imageX = x + i
					imageY = y + j

					r, g, b, a = (*img.Temp).At(imageX, imageY).RGBA()
					sumG = (sumG + (int(g) * matrice[i+1][j+1]))
					sumR = (sumR + (int(r) * matrice[i+1][j+1]))
					sumB = (sumB + (int(b) * matrice[i+1][j+1]))
					img.Out.Set(x, y, color.RGBA64{clamp(float64(sumR)), clamp(float64(sumG)), clamp(float64(sumB)), clamp(float64(a))})
				}
			}
		}
	}
}



//extract for this neighbor set, only the values for this color
//func (img *ImageTask) ExtractColors(neighbors *Neighbors) (*[9]float64, *[9]float64, *[9]float64){
//	//extract Red
//	var redArr [9]float64
//	for r := 0; r < 9; r++ {
//		redArr[r] = neighbors.inputs[r].Red
//	}
//	//extract Green
//	var greenArr [9]float64
//	for r := 0; r < 9; r++ {
//		greenArr[r] = neighbors.inputs[r].Green
//	}
//	//Extract Blue
//	var blueArr [9]float64
//	for r := 0; r < 9; r++ {
//		blueArr[r] = neighbors.inputs[r].Blue
//	}
//	return &redArr, &greenArr, &blueArr
//}

////add empty
//func (img *ImageTask) MakeEmpty() *Colors{
//	empty := Colors{Red: 0.0, Green: 0.0, Blue: 0.0, Alpha: 0.0}
//	return &empty
//}


//red = red + redArr[r] * kernel[r]
//green = green + greenArr[r] * kernel[r]
//blue = blue + blueArr[r] * kernel[r]
//red = red + (neighbors.inputs[r].Red * kernel[r])
//green = green + (neighbors.inputs[r].Green * kernel[r])
//blue = blue + (neighbors.inputs[r].Blue * kernel[r])
//red = red + clamp2(redArr[r] * kernel[r])
//green = green + clamp2(greenArr[r] * kernel[r])
//blue = blue + clamp2(blueArr[r] * kernel[r])
////fmt.Println("before",r,  red)
//red = clamp2(red + (neighbors.inputs[r].Red * kernel[r]))
////fmt.Println("after",r, red)
////fmt.Println("neighbor", r,  neighbors.inputs[r].Red)
//green = clamp2(green + (neighbors.inputs[r].Green * kernel[r]))
//blue = clamp2(blue + (neighbors.inputs[r].Blue * kernel[r]))
//red = clamp2(red)
//green = clamp2(green)
//blue = clamp2(blue)


//func (img *ImageTask) ClampTotal(){
//	bounds := img.Out.Bounds()
//	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
//		for x := bounds.Min.X; x < bounds.Max.X; x++ {
//			r, g, b, a := img.Temp.At(x, y).RGBA()
//			rx := float64(r)
//			gx := float64(g)
//			bx := float64(b)
//			//img.Out.Set(x, y, color.RGBA64{uint16(newResults.Red), uint16(newResults.Green), uint16(newResults.Blue), uint16(newResults.Alpha)})
//			img.Out.Set(x, y, color.RGBA64{clamp(rx), clamp(gx), clamp(bx), uint16(a)})
//		}
//	}
//}

//func (control *ParaChan) Decompose(img *ImageTask) *[]*[2]*[2]int{
//
//	var decomposed []*[2]int
//	for i := 0; i < control.Threads; i++ {
//		decomposed[i] = control.Translate(i, img)
//		if i == control.Threads - 1 {
//			//put all into
//			decomposed[i] = nil
//		}
//	}
//	return &decomposed
//}


//input targets []*effects
//func sequential(targets []*effects){
//	for _, image := range targets {
//		for _, effect := range image.Effects {
//			image.ImgTask.FilterControl(effect)
//		}
//	}
//}




//generator
//takes channel
//spawns thread, that push to said channel, usually with select
//return channel



//channel queue
//size X buffer =  queue

//piped, channel workers
//each worker, takes an ImageTas
//loop over filters
//block until hear done

//piped channel filter
//grid decomposition, with bounds
//spawn subworker threads specified

//piped channel subworker
//apply filter
//calls done

//piped channel aggregator
//aggregator channel
//done group


//test case
//requests to increment, array of array of numbers
//spawn workers
//spawn subworkers, segmenting,


//output targets []*effects
//the same one, pass around as pointers?





////returns common factors that the image could be evenly chunked into
//func (img *ImageTask) GridDecompose() *[]int{
//	bounds := img.Out.Bounds()
//	maxX := bounds.Max.X
//	maxY := bounds.Max.Y
//	res1 := factors(maxX)
//	res2 := factors(maxY)
//	res3 := cf(res1, res2)
//	return &res3
//}
//
//func factors(n int) []int {
//	res := []int{}
//	for t := n; t > 0; t-- {
//		if (n/t)*t == n {
//			res = append(res, t)
//		}
//	}
//	return res
//}
//
//func cf(l1 []int, l2 []int) []int {
//	res := []int{}
//	for len(l1) > 0 && len(l2) > 0 {
//		v1 := l1[0]
//		v2 := l2[0]
//		if v1 == v2 {
//			res = append(res, v1)
//			l2 = l2[1:]
//		}
//		if v2 > v1 {
//			l2 = l2[1:]
//		} else {
//			l1 = l1[1:]
//		}
//	}
//	return res
//}


//func main() {
//
//	jobSize := os.Args[1]
//	jobList := strings.Split(jobSize, "+")
//	for _, job := range jobList{
//		taskList := getOrders(job)
//		if len(os.Args) != 3{
//			//run sequential
//			sequential(taskList)
//			for _, effect := range taskList {
//				outPath := effect.OutPath
//				fixPath := "../data/out/" + outPath
//				_ = effect.ImgTask.Save(fixPath)
//			}
//		} else {
//			Mode := os.Args[2]
//			Threads, _ := strconv.Atoi(os.Args[3])
//			if Mode == "pipeline" {
//				parallelChannel(taskList, Threads)
//			} else {
//				parallelBSP(taskList, Threads)
//			}
//		}
//	}
//}



//fanout
//for each in coordinate range
//need new looping mechanism
//start = (x1, y1), end = (x2, y2)
//xTemp := x1
//yTemp := y1
//for (xTemp<= x2) && (yTemp <= y2){
//	neighbors := img.ExtractNeighbors(xTemp, yTemp, MaxX, MaxY)
//	newResults := img.ApplyConvulusion(neighbors, kernel)
//	img.Out.Set(xTemp, yTemp, color.RGBA64{clamp(newResults.Red), clamp(newResults.Green), clamp(newResults.Blue), clamp(newResults.Alpha)})
//	xTemp += 1
//	xTemp = xTemp % MaxX
//	yTemp +=1
//	yTemp = yTemp % MaxY
//
//}
//special lines for y1, y2 levels
//for x := 0; x < MaxX; x++ {
//	neighbors := img.ExtractNeighbors(x, y1, MaxX, MaxY)
//	newResults := img.ApplyConvulusion(neighbors, kernel)
//	img.Out.Set(x, y1, color.RGBA64{clamp(newResults.Red), clamp(newResults.Green), clamp(newResults.Blue), clamp(newResults.Alpha)})
//}
//for x := 0; x < MaxX; x++ {
//	neighbors := img.ExtractNeighbors(x, y2, MaxX, MaxY)
//	newResults := img.ApplyConvulusion(neighbors, kernel)
//	img.Out.Set(x, y2, color.RGBA64{clamp(newResults.Red), clamp(newResults.Green), clamp(newResults.Blue), clamp(newResults.Alpha)})
//}
////normal rectangle
//yMin2 := y1 +1






////BSP save
//package png
//
//import (
//"fmt"
//"image/color"
//"sync"
//)
//
////local queue
//type PoolQueue struct{
//	Queue []*Tickets
//}
//
////pointers to local queue pools
//type AggPoolQueue struct{
//	MQueue []*PoolQueue
//}
//
////tickets in each queue
//type Tickets struct{
//	Effect *Effects
//	Range [2][2]int
//	GAction string	//"p" = pointer, "s" = save
//	CondVar *sync.Cond
//	CondMutex *sync.Mutex
//	Counter *int
//	Done int
//	filter string
//}
//
//func (control *ParaChan) GenBQueues() *AggPoolQueue{
//
//	//initialize AggPoolQueue
//	var AggPoolQ AggPoolQueue
//	//AggPoolQ.MQueue = make([]PoolQueue, control.Threads)
//
//	//spawn threads
//	for r := 0; r <= control.Threads; r++ {
//		nQueue := PoolQueue{}
//		AggPoolQ.MQueue = append(AggPoolQ.MQueue, &nQueue)
//	}
//	return &AggPoolQ
//}
//
//func (control *ParaChan) GenBWorkers( AggPoolQueue *AggPoolQueue, group2 *WaitGroup){
//
//	group := waitGroup{}
//	for _, queue := range AggPoolQueue.MQueue {
//		//fmt.Println("here", queue)
//		go control.BWorker(queue, &group)
//		group.Add(1)
//	}
//	//fmt.Println(2)
//	group.Wait()
//	doneg(*group2)
//	//fmt.Println(1)
//}
//
////need mutext to set? maybe not because specific xy value
//func (control *ParaChan) BWorker(nQueue *PoolQueue, group *waitGroup){	//PoolQueue arg, Cond Var
//
//	fmt.Println(nQueue)
//	var position int
//	for j, image := range targets {
//
//	}
//	for {
//		//dequeue logic
//		//fmt.Println(nQueue.Queue)
//		//fmt.Println(len(nQueue.Queue))
//		if len(nQueue.Queue) < 1{	//loop until some value to dequeue
//			//fmt.Println("oops")
//			continue
//		}
//		if len(nQueue.Queue) < position{	//no new items yet
//			continue
//		}
//		newTicket := nQueue.Queue[position]		//pop
//		if newTicket.Done == 200{		//got done ticket
//			break
//		}
//		position += 1		//adjust
//
//
//		//convolusion logic
//		bounds := newTicket.Effect.ImgTask.Bounds
//		maxX := bounds.Max.X
//		maxY := bounds.Max.Y
//
//		targetCor := newTicket.Range
//		_, y1 := targetCor[0][0],targetCor[0][1]
//		_, y2 := targetCor[1][0], targetCor[1][1]
//
//		grey := 0
//
//		var kernel [9]float64
//		switch newTicket.filter{
//		case "S":
//			kernel = [9]float64{0.0, -1.0, 0.0, -1.0, 5.0, -1.0, 0.0, -1.0, 0.0}
//		case "E":
//			kernel = [9]float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}
//		case "B":
//			kernel = [9]float64{1/9.0, 1 / 9, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
//		case "G":
//			grey = 1
//		}
//
//		//apply convulsions
//		if grey != 0{
//			for y := y1; y < y2; y++ {
//				for x := 0; x < maxX; x++ {
//					r, g, b, a := newTicket.Effect.ImgTask.Temp.At(x, y).RGBA()
//					greyC := clamp(float64(r+g+b) / 3)
//					newTicket.Effect.ImgTask.Out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
//				}
//			}
//		}else{
//			//for each in coordinate range
//			for y := y1; y < y2; y++ {
//				for x := 0; x < maxX; x++ {
//					neighbors := newTicket.Effect.ImgTask.ExtractNeighbors(x, y, maxX, maxY)
//					newResults := newTicket.Effect.ImgTask.ApplyConvulusion(neighbors, kernel)
//					newTicket.Effect.ImgTask.Out.Set(x, y, color.RGBA64{clamp(newResults.Red), clamp(newResults.Green), clamp(newResults.Blue), clamp(newResults.Alpha)})
//				}
//			}
//		}
//		newTicket.CondMutex.Lock()		//do we need separate mutex
//		newC := *newTicket.Counter +1	//individual ticket counter
//		//fmt.Println(newC)
//		*newTicket.Counter +=1
//		//fmt.Println(*newTicket.Counter)
//		//newTicket.CondMutex.Unlock()
//		if newC == control.Threads{		//last thread enters
//			//do global step
//			if newTicket.GAction == "p"{
//				fmt.Println(1)
//				temp := newTicket.Effect.ImgTask.Temp
//				newTicket.Effect.ImgTask.Temp = newTicket.Effect.ImgTask.Out	//swaps pointer
//				newTicket.Effect.ImgTask.Out = temp
//			}else{
//				fmt.Println(2)
//				outPath := newTicket.Effect.OutPath
//				fixPath := "../data/out/" + outPath
//				_ = newTicket.Effect.ImgTask.Save(fixPath)
//			}
//			fmt.Println("Broadcast")
//			newTicket.CondVar.Broadcast()
//		}else{
//			//fmt.Println("waiting")
//			newTicket.CondVar.Wait()
//		}
//
//	}
//	group.Done()
//
//}
//
//
////for each effect, break into task Tickets, put into queue
//func (control *ParaChan) GenBTickets(effect *Effects, AggPoolQueue *AggPoolQueue, done int){
//
//	mutex := sync.Mutex{}
//	condVar := sync.NewCond(&mutex)
//	//mutex2 := sync.Mutex{}
//	//mutex2 = mutex
//	//fmt.Println(mutex == mutex2)
//	//fmt.Println(&mutex2)
//	conversions := control.Decompose(effect.ImgTask)
//	var action string
//
//	for j, filter := range effect.Effects {	//each filter
//		if j == len(effect.Effects)-1 {		//end of filter list
//			action = "s"
//		}else{
//			action = "p"
//		}
//		for i, single := range conversions.converted {	//for each separated range
//			counter := 0
//			//create Ticket
//			newTicket := Tickets{effect, single, action, condVar, &mutex, &counter, done, filter}
//			//fmt.Println(newTicket)
//			//pass to related queue
//			AggPoolQueue.MQueue[i].Queue = append(AggPoolQueue.MQueue[i].Queue, &newTicket)
//			//fmt.Println(AggPoolQueue.MQueue[i].Queue)
//
//		}
//	}
//
//	//lastly, append done ticket, if last task passed in
//	if done == 1 {
//		for i, _ := range conversions.converted { //for each separated range
//			//create Ticket
//			done = 200
//			newTicket := Tickets{Done: done}
//			AggPoolQueue.MQueue[i].Queue = append(AggPoolQueue.MQueue[i].Queue, &newTicket)
//			//fmt.Println("done")
//		}
//	}
//}



///archive channel
//for y := y1; y < y2; y++ {
//	var a int
//	var aMax int
//	if y == y1{
//		a = x1
//		aMax = MaxX
//	}else if y == (y2-1){
//		a = 0
//		aMax = x2
//	} else{
//		a = 0
//		aMax = MaxX
//	}
//
//	for x := a ; x < aMax; x++ {
//		neighbors := img.ExtractNeighbors(x, y, MaxX, MaxY)
//		newResults := img.ApplyConvulusion(neighbors, kernel)
//		img.Out.Set(x, y, color.RGBA64{clamp(newResults.Red), clamp(newResults.Green), clamp(newResults.Blue), clamp(newResults.Alpha)})
//
//	}
//}



//Waitgroup

//func WaitGroup(threads int){
//	done := make(chan int, threads)
//	var counter int
//
//	//spawn threads
//	for i := 0; i < threads; i++ {
//		go worker(done)
//	}
//
//	for counter != threads{
//		temp := <-done
//		counter = counter + temp
//	}
//}
//
//func worker(done chan int) {
//	done <- 1
//}
//
//func (img *ImageTask) ConvulsionControl(kernel [9]float64){
//	bounds := img.Out.Bounds()
//	maxX := bounds.Max.X
//	maxY := bounds.Max.Y
//	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
//		for x := bounds.Min.X; x < bounds.Max.X; x++ {
//			neighbors := img.ExtractNeighbors(x, y, maxX, maxY)
//			newResults := img.ApplyConvulusion(neighbors, kernel)
//			img.Out.Set(x, y, color.RGBA64{clamp(newResults.Red), clamp(newResults.Green), clamp(newResults.Blue), clamp(newResults.Alpha)})
//		}
//	}
//}