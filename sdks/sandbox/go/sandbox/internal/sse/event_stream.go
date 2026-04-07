package sse

type Event struct {
	Event string
	Data  []byte
}

type Stream struct {
	Events <-chan Event
	Done   <-chan error
}
