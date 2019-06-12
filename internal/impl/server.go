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
				websocketClient: &WebsocketClientImpl{
					ws: ws,
				},
				worker:       factory.Worker,
				targetWorker: factory.TargetWorker,
				id:           xid.New().String(),
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
		worker:     factory.Worker,
	}, nil
}

// ServerImpl :
type ServerImpl struct {
	httpServer *http.Server
	worker     domain.Worker
}

// Start :
func (server *ServerImpl) Start(
	send func(map[string]domain.Connection, string) error,
) error {
	go func() {
		server.httpServer.ListenAndServe()
	}()

	return server.worker.Run(send)
}
