package riemanngo

import (
	"testing"
)

func TestUnitNewUdpClient(t *testing.T) {
	client := NewUdpClient("127.0.0.1:5555");
	if client.addr != "127.0.0.1:5555" {
		t.Error("Error creating a new tcp client")
	}
}


func TestIntegSendEventUdp(t *testing.T) {
	c := NewUdpClient("127.0.0.1:5555");
	err := c.Connect(5)
	defer c.Close()
	if err != nil {
		t.Error("Error Udp client Connect")
	}
	result, err := SendEvent(c, &Event{
		Service: "LOOOl",
		Metric:  100,
		Tags: []string{"nonblocking"},
	})
	if result != nil || err != nil {
		t.Error("Error Udp client SendEvent")
	}
}

func TestIntegSendEventsUdp(t *testing.T) {
	c := NewUdpClient("127.0.0.1:5555");
	err := c.Connect(5)
	defer c.Close()
	if err != nil {
		t.Error("Error Udp client Connect")
	}
	events := []Event {
		Event{
			Service: "hello",
			Metric:  100,
			Tags: []string{"hello"},
		},
		Event{
			Service: "goodbye",
			Metric:  200,
			Tags: []string{"goodbye"},
		},
	}
	result, err := SendEvents(c, &events)
	if result != nil || err != nil {
		t.Error("Error Udp client SendEvent")
	}
}

func TestIntegUdpConnec(t *testing.T) {
	c := NewUdpClient("does.not.exists:8888");
	// should produce an error
	err := c.Connect(2)
	if err == nil {
		t.Error("Error, should fail")
	}
}
