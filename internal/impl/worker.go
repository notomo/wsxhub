package impl

import (
	"encoding/json"
	"io"
	"log"

	"github.com/notomo/wsxhub/internal/domain"
)

// WorkerImpl :
type WorkerImpl struct {
	Name         string
	Joined       chan domain.Connection
	Received     chan string
	Left         chan domain.Connection
	Done         chan bool
	Conns        map[string]domain.Connection
	OutputWriter io.Writer
}

// Run :
func (worker *WorkerImpl) Run() error {
	log.SetOutput(worker.OutputWriter)

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

			var targetMap map[string]interface{}
			if err := json.Unmarshal([]byte(message), &targetMap); err != nil {
				log.Printf("(%s) failed to unmarshal message: %s", worker.Name, err)
				break
			}

			for _, conn := range worker.Conns {
				if !conn.IsTarget(targetMap) {
					log.Printf("(%s) skipped", worker.Name)
					continue
				}
				if err := conn.Send(message); err != nil {
					log.Printf("(%s) failed to send: %s", worker.Name, err)
				}
				log.Printf("(%s) sent", worker.Name)
			}

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
func (worker *WorkerImpl) Receive(message string) error {
	worker.Received <- message
	return nil
}

// Delete :
func (worker *WorkerImpl) Delete(conn domain.Connection) error {
	worker.Left <- conn
	return nil
}
