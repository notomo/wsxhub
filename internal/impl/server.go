package impl

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/notomo/wsxhub/internal/domain"
	"github.com/rs/xid"
	"golang.org/x/net/websocket"
)

// ServerFactoryImpl :
type ServerFactoryImpl struct {
	Port                string
	Worker              domain.Worker
	TargetWorker        domain.Worker
	FilterClauseFactory domain.FilterClauseFactory
	OutputWriter        io.Writer
}

// Server :
func (factory *ServerFactoryImpl) Server(
	routes ...domain.Route,
) (domain.Server, error) {
	log.SetOutput(factory.OutputWriter)

	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Path, websocket.Handler(func(ws *websocket.Conn) {
			req := ws.Request()
			filterClause, err := factory.FilterClauseFactory.FilterClause(req.FormValue("filter"))
			if err != nil {
				log.Printf("failed to create filterClause: %s", err)
				return
			}

			debounce := 0
			debounceValue := req.FormValue("debounce")
			if debounceValue != "" {
				debounce, err = strconv.Atoi(debounceValue)
				if err != nil {
					log.Printf("failed to parse debounce: %s", err)
					return
				}
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
func (server *ServerImpl) Start() error {
	go func() {
		server.httpServer.ListenAndServe()
	}()

	return server.worker.Run()
}
