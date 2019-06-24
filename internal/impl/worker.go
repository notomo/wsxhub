package impl

import (
	"log"

	"github.com/notomo/wsxhub/internal/domain"
)

// WorkerImpl :
type WorkerImpl struct {
	Name               string
	Joined             chan domain.Connection
	Received           chan domain.Message
	Left               chan domain.Connection
	NotifiedSendResult chan error
	Done               chan bool
	Conns              map[string]domain.Connection
}

// NewWorker :
func NewWorker(name string) domain.Worker {
	return &WorkerImpl{
		Name:               name,
		Joined:             make(chan domain.Connection),
		Received:           make(chan domain.Message),
		Left:               make(chan domain.Connection),
		NotifiedSendResult: make(chan error),
		Done:               make(chan bool),
		Conns:              make(map[string]domain.Connection),
	}
}

// Run :
func (worker *WorkerImpl) Run() error {
	log.Printf("(%s) start", worker.Name)
	for {
		select {

		case conn := <-worker.Joined:
			worker.Conns[conn.ID()] = conn
			log.Printf("(%s) joined: %s, count: %d", worker.Name, conn.ID(), len(worker.Conns))

		case conn := <-worker.Left:
			delete(worker.Conns, conn.ID())
			log.Printf("(%s) left: %s, count: %d", worker.Name, conn.ID(), len(worker.Conns))

		case message := <-worker.Received:
			log.Printf("(%s) received", worker.Name)

			for _, conn := range worker.Conns {
				sent, err := conn.Send(message)
				if err != nil {
					log.Printf("(%s) failed to send: %s", worker.Name, err)
					continue
				}
				if sent {
					log.Printf("(%s) sent", worker.Name)
				}
			}

		case err := <-worker.NotifiedSendResult:
			if err != nil {
				log.Printf("(%s) failed to send: %s", worker.Name, err)
				continue
			}
			log.Printf("(%s) sent", worker.Name)

		case <-worker.Done:
			return nil
		}
	}
}

// Add :
func (worker *WorkerImpl) Add(conn domain.Connection) error {
	worker.Joined <- conn
	return nil
}

// Receive :
func (worker *WorkerImpl) Receive(message domain.Message) error {
	worker.Received <- message
	return nil
}

// Delete :
func (worker *WorkerImpl) Delete(conn domain.Connection) error {
	worker.Left <- conn
	return nil
}

// NotifySendResult :
func (worker *WorkerImpl) NotifySendResult(err error) {
	worker.NotifiedSendResult <- err
}
