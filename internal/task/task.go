package task

import (
	"Scotty/internal/client"
	"context"
	"log"
)

type State int

const (
	Initialize State = -1
)

type Task interface {
	// Next runs the next function of the tasks,
	// returns the next state
	Next(state State) (State, error)

	// Stop puts an end to the tasks
	Stop()

	// Context returns the context of the task
	Context() context.Context
}

type Base struct {
	// ID Of the tasks !
	ID string
	// Delay between each critical requests
	Delay int

	// context the context of the task
	// considered bad practice but for our use-case I
	// believe it will improve our over-all task purposes
	context context.Context
	cancel  func()

	// Client the client to use
	Client *client.CustomClient
}

// New creates a new TaskBase and sets our initial values
func New(id string) *Base {
	ctx, cancel := context.WithCancel(context.Background())

	taskBase := &Base{
		ID:      id,
		Client:  client.New(ctx),
		context: ctx,
		cancel:  cancel,
	}

	return taskBase
}

func (b *Base) Stop() {
	b.cancel()

}

// Context returns the task context
func (b *Base) Context() context.Context {
	return b.context
}

// RunTask runs a specific task
func RunTask(b Task) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	state := Initialize
	var err error
	for {
		select {
		default:
			state, err = b.Next(state)

			if err != nil {
				log.Println(err)
				return
			}
		case <-b.Context().Done():
			return
		}
	}
}
