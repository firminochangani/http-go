package http

type Header map[string]string

func (h Header) Set(name, value string) {
	h[name] = value
}

func (h Header) Get(name string) (string, bool) {
	value, exists := h[name]
	return value, exists
}

func (h Header) Del(name string) {
	delete(h, name)
}
