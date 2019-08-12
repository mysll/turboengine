package workqueue

import (
	"turboengine/common/log"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin"
)

const (
	Name           = "WorkQueue"
	MAX_GO_ROUTINE = 16
	MAX_QUEUE      = 1024
)

type Task interface {
	Run()
	Complete()
}

type WorkQueue struct {
	srv      api.Service
	jobs     []chan Task
	complete chan Task
	attachid uint64
	closing  bool
}

func (w *WorkQueue) Prepare(srv api.Service) {
	w.srv = srv
	w.jobs = make([]chan Task, MAX_GO_ROUTINE)
	w.complete = make(chan Task, 512)
	w.attachid = srv.Attach(w.invokeComplete)
	for i := 0; i < MAX_GO_ROUTINE; i++ {
		w.jobs[i] = make(chan Task, MAX_QUEUE)
	}
}

func (w *WorkQueue) Run() {
	for i := 0; i < MAX_GO_ROUTINE; i++ {
		go w.work(w.jobs[i])
	}
}

func (w *WorkQueue) Shut(api.Service) {
	w.closing = true
	for i := 0; i < MAX_GO_ROUTINE; i++ {
		close(w.jobs[i])
	}
	close(w.complete)
	w.srv.Deatch(w.attachid)
}

func (w *WorkQueue) Handle(cmd string, args ...interface{}) interface{} {
	return nil
}

func (w *WorkQueue) Schedule(filt int, task Task) bool {
	if w.closing {
		return false
	}
	idx := int(uint(filt) % uint(len(w.jobs)))
	select {
	case w.jobs[idx] <- task:
		return true
	default:
		log.Error("too much jobs")
		return false
	}
}

func (w *WorkQueue) work(queue chan Task) {
	for t := range queue {
		t.Run()
		if w.closing {
			return
		}
		w.complete <- t
	}
}

// service update
func (w *WorkQueue) invokeComplete(t *utils.Time) {
	for {
		select {
		case t := <-w.complete:
			t.Complete()
		default:
			return
		}
	}
}

func init() {
	plugin.Register(Name, &WorkQueue{})
}
