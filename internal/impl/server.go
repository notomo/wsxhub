package impl

import (
	"log"
	"net/http"

	"github.com/notomo/wsxhub/internal/domain"
	"github.com/rs/xid"
	"golang.org/x/net/websocket"
)

// ServerFactoryImpl :
type ServerFactoryImpl struct {
	Port         string
	Worker       domain.Worker
	TargetWorker domain.Worker
}

// Server :
func (factory *ServerFactoryImpl) Server(
	routes ...domain.Route,
) (domain.Server, error) {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Path, websocket.Handler(func(ws *websocket.Conn) {
			conn := &ConnectionImpl{
				ws:           ws,
				worker:       factory.Worker,
				targetWorker: factory.TargetWorker,
				id:           xid.New().String(),
				done:         make(chan bool),
			}
			if err := route.Handler(conn); err != nil {
				log.Print(err)
			}
		}))
	}

	server := &http.Server{
		Addr:    ":" + factory.Port,
		Handler: mux,
	}

	return &ServerImpl{
		httpServer: server,
		conns:      make(map[string]domain.Connection),
		worker:     factory.Worker,
	}, nil
}

// ServerImpl :
type ServerImpl struct {
	httpServer *http.Server
	conns      map[string]domain.Connection
	worker     domain.Worker
}

// Connections :
func (server *ServerImpl) Connections() map[string]domain.Connection {
	return server.conns
}

// Add :
func (server *ServerImpl) Add(conn domain.Connection) error {
	server.conns[conn.ID()] = conn
	return nil
}

// Delete :
func (server *ServerImpl) Delete(conn domain.Connection) error {
	delete(server.conns, conn.ID())
	return nil
}

// Start :
func (server *ServerImpl) Start(
	send func(map[string]domain.Connection, string) error,
) error {
	go func() {
		server.httpServer.ListenAndServe()
	}()

	return server.worker.Run(server, send)
}
