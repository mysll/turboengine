package workqueue

import (
	"turboengine/common/log"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin"
)

var (
	Name           = "WorkQueue"
	MAX_GO_ROUTINE = 16
	MAX_QUEUE      = 128
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
	hashmask uint64
}

func (w *WorkQueue) Prepare(srv api.Service, args ...any) {
	w.srv = srv
	if !utils.IsPowerOfTwo(MAX_GO_ROUTINE) {
		panic("MAX_GO_ROUTINE must be power of two ")
	}
	w.hashmask = uint64(MAX_GO_ROUTINE - 1)
	w.jobs = make([]chan Task, MAX_GO_ROUTINE)
	w.complete = make(chan Task, MAX_QUEUE*2)
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
	w.srv.Detach(w.attachid)
}

func (w *WorkQueue) Handle(cmd string, args ...any) any {
	return nil
}

func (w *WorkQueue) Schedule(hashkey uint64, task Task) bool {
	if w.closing {
		return false
	}

	select {
	case w.jobs[hashkey&w.hashmask] <- task:
		return true
	default:
		log.Error("too many jobs")
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
