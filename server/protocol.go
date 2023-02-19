package server

import (
	"fmt"
)

type Protocol map[string]Method

func NewProtocol() Protocol {
	res := make(Protocol)

	return res
}

func (p Protocol) Find(name string) Method {
	res, ok := p[name]

	if !ok {
		return nil
	}

	return res
}

type NotificationHandler[PA any] func(s *Session, params PA)

func RegisterNotification[PA any](p *Protocol, method string, handler NotificationHandler[PA]) {
	p.register(notification{
		n: method,
	})
}

type RequestHandler[PA any, RE any] func(s *Session, params PA, result *RE) error

func RegisterRequest[PA any, RE any](p *Protocol, method string, handler RequestHandler[PA, RE]) {
	p.register(request{
		n: method,
	})
}

func (p Protocol) register(m Method) {
	name := m.name()

	if p.Find(name) != nil {
		panic(fmt.Errorf("protocol already has a method named '%s'", name))
	}

	p[name] = m
}

type Method interface {
	name() string
}

type notification struct {
	n string
}

func (n notification) name() string {
	return n.n
}

type request struct {
	n string
}

func (n request) name() string {
	return n.n
}
