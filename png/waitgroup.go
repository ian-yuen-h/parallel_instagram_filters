package png

import (
	"sync/atomic"
)

type WaitGroup interface {
	Add(amount uint)
	Done()
	Wait()
}

type waitGroup struct{
	counter int32
}

// NewWaitGroup returns a instance of a waitgroup
// This instance must be a pointer and should not
// be copied after creation.
func NewWaitGroup() WaitGroup {
	newgroup := &waitGroup{}

	return newgroup
}

func (w *waitGroup) Add(amount uint) {

	atomic.AddInt32(&(w.counter), int32(amount))

}
func (w *waitGroup) Done() {
	atomic.AddInt32(&(w.counter), -1)
}
func (w *waitGroup) Wait(){
	for w.counter != 0{
	}
}

//must have this function, to use the struct funcs
func doneg(w WaitGroup){
	w.Done()
}