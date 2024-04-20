package http_go

type Router interface {
	Handle(r *Request, w *Response) error
}

type Handler func(r *Request, w *Response) error
type httpMethod = string
type routePaths map[string]Handler

type ServerDefaultNaiveRouter struct {
	paths map[httpMethod]routePaths
}

func (h ServerDefaultNaiveRouter) Handle(r *Request, w *Response) error {
	handlerFunc, exists := h.paths[r.Method][r.URL.Path]
	if !exists {
		w.WriteStatus(404)
		return w.Write([]byte("Not found"))
	}

	return handlerFunc(r, w)
}

func NewServerDefaultRouter() ServerDefaultNaiveRouter {
	paths := make(map[httpMethod]routePaths)
	paths[MethodGET] = make(routePaths)
	paths[MethodPOST] = make(routePaths)
	paths[MethodPUT] = make(routePaths)
	paths[MethodDELETE] = make(routePaths)
	paths[MethodPATCH] = make(routePaths)
	paths[MethodCONNECT] = make(routePaths)
	paths[MethodHEAD] = make(routePaths)
	paths[MethodOPTIONS] = make(routePaths)
	paths[MethodTRACE] = make(routePaths)

	return ServerDefaultNaiveRouter{
		paths: paths,
	}
}

func (h ServerDefaultNaiveRouter) GET(path string, handler Handler) {
	// how to register a route having in mind performant lookups
	// naive implementation of a router
	h.paths[MethodGET][path] = handler
}
