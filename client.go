// A Riemann client for Go, featuring concurrency, sending events and state updates, queries,
// and feature parity with the reference implementation written in Ruby.
//
// Copyright (C) 2014 by Christopher Gilbert <christopher.john.gilbert@gmail.com>
package riemanngo

import (
	pb "github.com/golang/protobuf/proto"
	"github.com/riemann/riemann-go-client/proto"
)

// Send an event
func SendEvent(c Client, e *Event) error {
	epb, err := EventToProtocolBuffer(e)
	if err != nil {
		return err
	}
	message := &proto.Msg{}
	message.Events = append(message.Events, epb)

	_, err = c.SendMaybeRecv(message)
	return err
}

// Query the server for events
func QueryIndex(c TcpClient, q string) ([]Event, error) {
	query := &proto.Query{}
	query.String_ = pb.String(q)

	message := &proto.Msg{}
	message.Query = query

	response, err := c.SendRecv(message)
	if err != nil {
		return nil, err
	}

	return ProtocolBuffersToEvents(response.GetEvents()), nil
}

