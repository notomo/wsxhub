package domain

// ServerFactory :
type ServerFactory interface {
	Server(...Route) (Server, error)
}

// Server :
type Server interface {
	Start() error
}

// Route :
type Route struct {
	Path    string
	Handler func(Connection) error
}

// NewRoute :
func NewRoute(path string, handler func(Connection) error) Route {
	return Route{
		Path:    path,
		Handler: handler,
	}
}
