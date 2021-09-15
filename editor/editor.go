package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"proj2/png"
	"strconv"
	"strings"

	//"strconv"
)

func main() {

	jobSize := os.Args[1]
	jobList := strings.Split(jobSize, "+")
	for _, job := range jobList{
		taskList := getOrders(job)
		if len(os.Args) != 4{
			sequential(taskList)
			for _, effect := range taskList {
				_ = effect.ImgTask.Save(effect.OutPath)
			}
		} else{
			Mode := os.Args[2]
			Threads, _  := strconv.Atoi(os.Args[3])
			if Mode == "pipeline"{
				parallelChannel(taskList, Threads)
			}else{
				parallelBSP(taskList, Threads)
			}
		}
	}
}

//extract photos, return slice of effects = wrapper for ImgTask
func getOrders(jobSize string)[]*png.Effects{
	file, err := os.Open("../data/effects.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var targets []*png.Effects
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var t png.Effects
		a := string(scanner.Text())
		err = json.Unmarshal([]byte(a), &t)
		if err != nil {
			fmt.Println("error:", err)
		}
		targets = append(targets, &t)
	}

	//fix inPaths, outPaths
	for _, effect := range targets {
		path := effect.InPath
		fixPath := "../data/in/"+jobSize+"/"+path
		effect.InPath = fixPath
		path = effect.OutPath
		fixPath = "../data/out/" + jobSize +"_"+path
		effect.OutPath = fixPath
	}

	targets = genImageTask(targets)
	return targets
}

func genImageTask(targets []*png.Effects) []*png.Effects{
	for _, effect := range targets {
		pngImg, _ := png.Load(effect.InPath)		//loads each image as RGBA64
		effect.ImgTask = pngImg
	}
	return targets
}

func sequential(targets []*png.Effects){
	for _, image := range targets {
		for _, effect := range image.Effects {
			image.ImgTask.SeqFilterControl(effect)
		}
	}
}

func parallelChannel(targets []*png.Effects, threads int){
	paraControl := png.ParaChan{targets, threads}
	queueChan := paraControl.QueueRelease()	//queue set up
	resultStream := paraControl.GenWorkers(queueChan) //workers ready
	paraControl.ResultAggregator(resultStream)	//result set
}

func parallelBSP(targets []*png.Effects, threads int){
	paraControl := png.ParaChan{targets, threads}
	aggPool := paraControl.GenBQueues()

	var done int
	//for each in targets
	for j, image := range targets {
		if j == len(targets)-1 {		//end of filter list
			done = 1
		}else{
			done = 0
		}
		paraControl.GenBTickets(image, aggPool, done)
	}
	group2 := png.NewWaitGroup()
	go paraControl.GenBWorkers(aggPool, &group2)
	group2.Add(1)
	group2.Wait()
}