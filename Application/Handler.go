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

/*
	currently handler re-usability requires plenty of overhead due to the limited implementation of generics in Go, especially with multiple fields to access
	when go introduces generic field access this would become a lot less cumbersome, as of 1.20 this is not implemented
	using a union interface as the argument type of the handler would reduce this to

	type ConcreteHandlerAlpha struct {
		baseHandler
	}

	type ConcreteHandlerBeta struct {
		success bool
	}

	type ReceiverBeta struct {
		success bool
	}

	type AbstractReceiver interface {
		*ConcreteReceiverAlpha | *ConcreteReceiverBeta
	}

	func (h *ConcreteHandlerAlpha) execute(r AbstractReceiver) {
		r.success = true
	}
*/

/*
	Current example of re-usability without generics, using setters and getters

	type HandlerAlpha struct {
		baseHandler
	}

	type ReceiverAlpha struct {
		success bool
	}

	type ReceiverBeta struct {
		success bool
	}

	type SuccessReceiver interface {
		setSuccess(success bool)
	}

	func (r *ReceiverAlpha) setSuccess(success bool) {
		r.success = success
	}

	func (r *ReceiverBeta) setSuccess(success bool) {
		r.success = success
	}

	func (h *HandlerAlpha) execute(r SuccessReceiver) {
		r.setSuccess(true)
	}
*/
