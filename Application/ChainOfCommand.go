package email

type baseHandler struct {
	next handler
}

// setNext sets the next handler in chain, returns the next for chaining setNext
func (h *baseHandler) setNext(next handler) handler {
	h.next = next
	return next
}

type handler interface {
	execute(r *receiver)
	setNext(handler) handler
}
