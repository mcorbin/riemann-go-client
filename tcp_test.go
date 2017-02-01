package riemanngo

import (
	"testing"
)

func TestUnitNewTcpClient(t *testing.T) {
	client := NewTcpClient("127.0.0.1:5555");
	if client.addr != "127.0.0.1:5555" {
		t.Error("Error creating a new tcp client")
	}
}

// i use this Riemann config for integration test
// (logging/init {:file "/var/log/riemann/riemann.log"})

// (let [host "0.0.0.0"]
//   (tcp-server)
//   (udp-server))

// (periodically-expire 60)

// (streams
//  (where (not (service #"^riemann "))
//    (index)
//    #(info %)))

func TestIntegSendEventTcp(t *testing.T) {
	c := NewTcpClient("127.0.0.1:5555");
	err := c.Connect(5)
	defer c.Close()
	if err != nil {
		t.Error("Error Tcp client Connect")
	}
	result, err := SendEvent(c, &Event{
		Service: "LOOOl",
		Metric:  100,
		Tags: []string{"nonblocking"},
	})
	if ! *result.Ok {
		t.Error("Error Tcp client SendEvent")
	}
}

func TestIntegSendEventsTcp(t *testing.T) {
	c := NewTcpClient("127.0.0.1:5555");
	err := c.Connect(5)
	defer c.Close()
	if err != nil {
		t.Error("Error Tcp client Connect")
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
	if ! *result.Ok {
		t.Error("Error Tcp client SendEvent")
	}
}


func TestIntegQueryIndex(t *testing.T) {
	c := NewTcpClient("127.0.0.1:5555");
	err := c.Connect(5)
	defer c.Close()
	if err != nil {
		t.Error("Error Tcp client Connect")
	}
	events := []Event {
		Event{
			Host: "foobaz",
			Service: "golang",
			Metric:  100,
			Tags: []string{"hello"},
		},
		Event{
			Host: "foobar",
			Service: "golang",
			Metric:  200,
			Tags: []string{"goodbye"},
		},
	}
	result, err := SendEvents(c, &events)
	if ! *result.Ok {
		t.Error("Error Tcp client SendEvent")
	}
	queryResult, err := QueryIndex(c, "(service = \"golang\")")
	if len(queryResult) != 2 {
		t.Error("Error Tcp client QueryIndex")
	}
}

func TestIntegTcpConnec(t *testing.T) {
	c := NewTcpClient("does.not.exists:8888");
	// should produce an error
	err := c.Connect(2)
	if err == nil {
		t.Error("Error, should fail")
	}
}
