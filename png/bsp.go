package png

import (
	"sync"
)

//local queue pool
type PoolQueue struct{
	Queue []*Tickets
}

//pointers to local queue pools
type AggPoolQueue struct{
	MQueue []*PoolQueue
}

//tickets in each queue
type Tickets struct{
	Effect *Effects
	Range [2][2]int
	GAction string
	CondVar *sync.Cond
	CondMutex *sync.Mutex
	Counter *int32
	Done int
	filter string
}

func (control *ParaChan) GenBQueues() *AggPoolQueue{
	var AggPoolQ AggPoolQueue
	for r := 0; r < control.Threads; r++ {
		nQueue := PoolQueue{}
		AggPoolQ.MQueue = append(AggPoolQ.MQueue, &nQueue)
	}
	return &AggPoolQ
}

func (control *ParaChan) GenBWorkers( AggPoolQueue *AggPoolQueue, group2 *WaitGroup){

	group := waitGroup{}
	for i := 0; i < control.Threads ; i++ {
		go control.BWorker(AggPoolQueue.MQueue[i], &group)
		group.Add(1)
	}
	group.Wait()
	doneg(*group2)
}

func (control *ParaChan) BWorker(nQueue *PoolQueue, group *waitGroup){

	if nQueue == nil {
		return
	}
	for _ , newTicket := range nQueue.Queue {
		if newTicket.Done == 200{
			break
		}

		targetCor := newTicket.Range
		_, y1 := targetCor[0][0],targetCor[0][1]
		_, y2 := targetCor[1][0], targetCor[1][1]

		grey := 0
		kernel := newTicket.Effect.ImgTask.getKernel(newTicket.filter)
		if kernel == [9]float64{}{
			grey =1
		}

		newTicket.Effect.ImgTask.MConvulsionControl(grey, y1, y2, kernel)

		newTicket.CondMutex.Lock()
		*newTicket.Counter +=1
		if int(*newTicket.Counter) == control.Threads{		//last thread enters, do global step
			if newTicket.GAction == "p"{
				temp := *newTicket.Effect.ImgTask.Temp
				newTicket.Effect.ImgTask.Temp = newTicket.Effect.ImgTask.Out
				newTicket.Effect.ImgTask.Out = &temp
			}else{
				_ = newTicket.Effect.ImgTask.Save2(newTicket.Effect.OutPath)
			}
			newTicket.CondVar.Broadcast()
		}else{
			newTicket.CondVar.Wait()
		}
		newTicket.CondMutex.Unlock()
	}
	group.Done()
}



//for each effect, break into task Tickets, put into queue
func (control *ParaChan) GenBTickets(effect *Effects, AggPoolQueue *AggPoolQueue, done int){

	conversions := control.Decompose(effect.ImgTask)
	var action string

	mutex := sync.Mutex{}
	condVar := sync.NewCond(&mutex)

	for j, filter := range effect.Effects {	//each filter
		if j == len(effect.Effects)-1 {		//end of filter list
			action = "s"
		}else{
			action = "p"
		}
		var counter int32
		for i, single := range conversions.converted {	//for each separated range
			newTicket := Tickets{effect, single, action, condVar, &mutex, &counter, done, filter}
			AggPoolQueue.MQueue[i].Queue = append(AggPoolQueue.MQueue[i].Queue, &newTicket)

		}
	}

	if done == 1 {		//done ticket
		for i, _ := range conversions.converted {
			done = 200
			newTicket := Tickets{Done: done}
			AggPoolQueue.MQueue[i].Queue = append(AggPoolQueue.MQueue[i].Queue, &newTicket)
		}
	}
}
