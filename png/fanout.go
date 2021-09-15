package png

type ParaChan struct{
	Targets []*Effects
	Threads int
}

func (control *ParaChan) QueueRelease() <-chan *Effects{
	queueChan := make(chan *Effects, control.Threads)
	go func() {
		for _, effect := range control.Targets {
			queueChan <- effect
		}
	}()
	return queueChan
}

func (control *ParaChan) GenWorkers(queueChan <-chan *Effects) chan *Effects{

	resultStream := make(chan *Effects, control.Threads)
	for r := 0; r <= control.Threads; r++ {
		go control.Worker(queueChan, resultStream)
	}
	return resultStream
}

func (control *ParaChan) ResultAggregator(resultStream chan *Effects){
	var counter int
	for{
		effect, more := <-resultStream
		if more {
			_ = effect.ImgTask.Save(effect.OutPath)
			counter = counter + 1
			if counter == len(control.Targets){		//finished processing all
				break
			}
		}
	}
}

func (control *ParaChan) Worker(queueChan <-chan *Effects, resultStream chan<- *Effects) {
	for{
		effect, more := <-queueChan
		if more {
			for _, filter := range effect.Effects {
				control.SubWorkerControl(effect.ImgTask, filter)
			}
			resultStream<- effect
		}
	}
}

func (control *ParaChan) SubWorkerControl(img *ImageTask, filter string){

	done := make(chan int, control.Threads)
	var counter int

	subPackQueueChan := control.SubPacketQueue(img)

	for i := 0; i < control.Threads; i++ {
		go control.Subworker(img, done, subPackQueueChan, filter)
	}
	for counter != control.Threads{
		temp := <-done
		counter = counter + temp
	}
	temp := img.Temp
	img.Temp = img.Out	//swaps pointer
	img.Out = temp
}

func (control *ParaChan) Subworker(img *ImageTask, done chan int, subPackQueueChan <-chan [2][2]int, filter string) {

	targetCor := <- subPackQueueChan
	_, y1 := targetCor[0][0],targetCor[0][1]
	_, y2 := targetCor[1][0], targetCor[1][1]

	grey := 0
	kernel := img.getKernel(filter)
	if kernel == [9]float64{}{
		grey =1
	}
	img.MConvulsionControl(grey, y1, y2, kernel)
	done <- 1
}

func (control *ParaChan) SubPacketQueue(img *ImageTask) <-chan [2][2]int{
	subPackQueueChan := make(chan [2][2]int, control.Threads)
	decomposed := control.Decompose(img)
	go func() {
		for _, packet := range decomposed.converted {
			subPackQueueChan <- packet
		}
	}()
	return subPackQueueChan
}

type Conversions struct {
	converted [][2][2]int
}

func (control *ParaChan) Decompose(img *ImageTask) *Conversions{

	total := img.MaxX * img.MaxY
	threads := control.Threads

	var convertObj Conversions
	convertObj.converted = make([][2][2]int, 0)

	jump := total/threads

	var endRange int
	var currRange int
	for i := 0; i < threads; i++ {
		var start [2]int
		var end [2]int
		if i == 0{
			start = [2]int{0,0}
			endRange = currRange+jump-1
			end = control.Translate(endRange, img)
		} else if (i == threads - 1) {
			end = [2]int{img.MaxX, img.MaxY}
			start = control.Translate(currRange, img)
		}else{
			endRange = currRange+jump-1
			start = control.Translate(currRange, img)
			end = control.Translate(endRange, img)
		}
		x1, y1 := start[0], start[1]
		x2, y2 := end[0], end[1]
		temp := [2][2]int{{x1, y1}, {x2, y2}}
		convertObj.converted = append(convertObj.converted, temp )
		currRange = currRange+jump
	}
	return &convertObj
}

func (control *ParaChan) Translate(i int, img *ImageTask) [2]int{

	row := i/img.MaxX
	column := i/img.MaxY

	res := [2]int{column, row}
	return res
}