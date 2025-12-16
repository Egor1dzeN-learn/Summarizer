package tasksq

import (
	"github.com/google/uuid"
)

type Task[Request any] struct {
	ID string

	Req *Request
}

type CompletedTask[Request any, Result any] struct {
	Task[Request]

	Res *Result
	Err error
}

type TaskHandler[Request any, Result any] interface {
	Process(task Task[Request]) (*Result, error)
	Proceed(completedTask CompletedTask[Request, Result])
}

type TaskQueue[Request any, Result any] struct {
	handler TaskHandler[Request, Result]

	inQueue  chan Task[Request]
	outQueue chan CompletedTask[Request, Result]
}

func NewTaskQueue[Request any, Result any](handler TaskHandler[Request, Result], workers int, queueSize int) TaskQueue[Request, Result] {
	tq := TaskQueue[Request, Result]{
		inQueue:  make(chan Task[Request], queueSize),
		outQueue: make(chan CompletedTask[Request, Result], queueSize),
		handler:  handler,
	}

	for range workers {
		go tq.wkIn()
		go tq.wkOut()
	}

	return tq
}

func (tq *TaskQueue[Request, Result]) Enqueue(req *Request, bind func(taskID string)) string {
	id := uuid.New().String()
	bind(id)
	tq.inQueue <- Task[Request]{ID: id, Req: req}
	return id
}

func (tq *TaskQueue[Request, Result]) Close() {
	close(tq.outQueue)
	close(tq.inQueue)
}

func (tq *TaskQueue[Request, Result]) wkIn() {
	for task := range tq.inQueue {
		completedTask := CompletedTask[Request, Result]{Task: task}

		res, err := tq.handler.Process(task)
		if err != nil {
			completedTask.Err = err
		} else {
			completedTask.Res = res
		}

		tq.outQueue <- completedTask
	}
}
func (tq *TaskQueue[Request, Result]) wkOut() {
	for task := range tq.outQueue {
		tq.handler.Proceed(task)
	}
}
