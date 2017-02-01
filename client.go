// A Riemann client for Go, featuring concurrency, sending events and state updates, queries
//
// Copyright (C) 2014 by Christopher Gilbert <christopher.john.gilbert@gmail.com>
package riemanngo

import (
	"github.com/riemann/riemann-go-client/proto"
)

// Client is an interface to a generic client
type Client interface {
	Send(message *proto.Msg) (*proto.Msg, error)
	Connect() error
	Close() error
}

// request encapsulates a request to send to the Riemann server
type request struct {
	message     *proto.Msg
	response_ch chan response
}

// response encapsulates a response from the Riemann server
type response struct {
	message *proto.Msg
	err     error
}

// Send an event
func SendEvent(c Client, e *Event) (*proto.Msg, error) {
	epb, err := EventToProtocolBuffer(e)
	if err != nil {
		return nil, err
	}
	message := &proto.Msg{}
	message.Events = append(message.Events, epb)

	msg, err := c.Send(message)
	return msg, err
}
