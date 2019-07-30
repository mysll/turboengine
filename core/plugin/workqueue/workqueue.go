package workqueue

import (
	"turboengine/core/api"
)

const (
	NAME           = "workqueue"
	MAX_GO_ROUTINE = 10
	MAX_QUEUE      = 100
)

type Task interface {
	Run()
}

type WorkQueue struct {
	srv  api.Service
	jobs []chan Task
}

func (w *WorkQueue) Prepare(srv api.Service) {
	w.srv = srv
	w.jobs = make([]chan Task, MAX_GO_ROUTINE)
	for i := 0; i < MAX_GO_ROUTINE; i++ {
		w.jobs[i] = make(chan Task, MAX_QUEUE)
		go w.Work(w.jobs[i])
	}
}

func (w *WorkQueue) Shut(api.Service) {
	for i := 0; i < MAX_GO_ROUTINE; i++ {
		close(w.jobs[i])
	}
}

func (w *WorkQueue) Handle(cmd string, args ...interface{}) interface{} {
	return nil
}

func (w *WorkQueue) Schedule(filt int, task Task) {
	idx := int(uint(filt) % uint(len(w.jobs)))
	w.jobs[idx] <- task
}

func (w *WorkQueue) Work(queue chan Task) {
	for t := range queue {
		t.Run()
	}
}
