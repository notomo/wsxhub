package impl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/notomo/wsxhub/internal/domain"
	"github.com/rs/xid"
)

// ServerFactoryImpl :
type ServerFactoryImpl struct {
	Port                string
	Worker              domain.Worker
	TargetWorker        domain.Worker
	FilterClauseFactory domain.FilterClauseFactory
	OutputWriter        io.Writer
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Server :
func (factory *ServerFactoryImpl) Server(
	routes ...domain.Route,
) (domain.Server, error) {
	log.SetOutput(factory.OutputWriter)

	mux := http.NewServeMux()
	for _, route := range routes {
		mux.HandleFunc(route.Path, func(w http.ResponseWriter, req *http.Request) {
			filterClause, err := factory.FilterClauseFactory.FilterClause(req.FormValue("filter"))
			if err != nil {
				msg := fmt.Sprintf("failed to create filterClause: %s", err)
				http.Error(w, msg, http.StatusBadRequest)
				log.Printf(msg)
				return
			}

			debounce := 0
			debounceValue := req.FormValue("debounce")
			if debounceValue != "" {
				debounce, err = strconv.Atoi(debounceValue)
				if err != nil {
					msg := fmt.Sprintf("failed to parse debounce: %s", err)
					http.Error(w, msg, http.StatusBadRequest)
					log.Printf(msg)
					return
				}
			}

			ws, err := upgrader.Upgrade(w, req, nil)
			if err != nil {
				log.Printf("failed to upgrade: %s", err)
				return
			}

			conn := &ConnectionImpl{
				websocketClient: &WebsocketClientImpl{
					ws: ws,
				},
				worker:       factory.Worker,
				targetWorker: factory.TargetWorker,
				id:           xid.New().String(),
				filterClause: filterClause,
				debounce:     debounce,
			}
			defer conn.Close()

			if err := route.Handler(conn); err != nil {
				log.Print(err)
			}
		})
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
func (server *ServerImpl) Start() error {
	go func() {
		server.httpServer.ListenAndServe()
	}()

	return server.worker.Run()
}
