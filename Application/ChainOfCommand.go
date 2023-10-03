package email

type baseHandler struct {
	next Handler
}

// setNext sets the next Handler in chain, returns the next for chaining setNext
func (h *baseHandler) setNext(next Handler) Handler {
	h.next = next
	return next
}

type Handler interface {
	execute(r *receiver)
	setNext(Handler) Handler
}
