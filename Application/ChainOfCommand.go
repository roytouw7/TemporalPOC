package email

type baseHandler struct {
	next handler
}

func (h *baseHandler) setNext(next handler) {
	h.next = next
}

type handler interface {
	execute(r *receiver)
	setNext(handler)
}
