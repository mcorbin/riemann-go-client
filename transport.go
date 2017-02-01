package riemanngo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	pb "github.com/golang/protobuf/proto"
	"github.com/riemann/riemann-go-client/proto"
)

// Client is an interface to a generic client
type Client interface {
	SendRecv(message *proto.Msg) (*proto.Msg, error)
	Connect() error
	Close() error
	SendMaybeRecv(message *proto.Msg) (*proto.Msg, error)
}

// TcpClient is a type that implements the Client interface
type TcpClient struct {
	addr         string
	conn         net.Conn
	requestQueue chan request
}

// UdpClient is a type that implements the Client interface
type UdpClient  struct {
	addr         string
	conn         net.Conn
	requestQueue chan request
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

// MAX_UDP_SIZE is the maximum allowed size of a UDP packet before automatically failing the send
const MAX_UDP_SIZE = 16384

// NewTcpClient - Factory
func NewTcpClient(addr string) *TcpClient {
	t := &TcpClient{
		addr: addr,
		requestQueue: make(chan request),
	}
	go t.runRequestQueue()
	return t
}

// NewUdpClient - Factory
func NewUdpClient(addr string) *UdpClient {
	t := &UdpClient{
		addr: addr,
		requestQueue: make(chan request),
	}
	go t.runRequestQueue()
	return t
}

func (c *TcpClient) Connect() error {
	tcp, err := net.DialTimeout("tcp", c.addr, time.Second*5)
	if err != nil {
		return err
	}
	c.conn = tcp
	return nil
}


func (c *UdpClient) Connect() error {
	udp, err := net.DialTimeout("udp", c.addr, time.Second*5)
	if err != nil {
		return err
	}
	c.conn = udp
	return nil
}

// TcpClient implementation of SendRecv, queues a request to send a message to the server
func (t *TcpClient) SendRecv(message *proto.Msg) (*proto.Msg, error) {
	response_ch := make(chan response)
	t.requestQueue <- request{message, response_ch}
	r := <-response_ch
	return r.message, r.err
}

// TcpClient implementation of SendMaybeRecv, queues a request to send a message to the server
func (t *TcpClient) SendMaybeRecv(message *proto.Msg) (*proto.Msg, error) {
	return t.SendRecv(message)
}

// Close will close the TcpClient
func (t *TcpClient) Close() error {
	close(t.requestQueue)
	err := t.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// runRequestQueue services the TcpClient request queue
func (t *TcpClient) runRequestQueue() {
	for req := range t.requestQueue {
		message := req.message
		response_ch := req.response_ch

		msg, err := t.execRequest(message)

		response_ch <- response{msg, err}
	}
}

// execRequest will send a TCP message to Riemann
func (t *TcpClient) execRequest(message *proto.Msg) (*proto.Msg, error) {
	msg := &proto.Msg{}
	data, err := pb.Marshal(message)
	if err != nil {
		return msg, err
	}
	b := new(bytes.Buffer)
	if err = binary.Write(b, binary.BigEndian, uint32(len(data))); err != nil {
		return msg, err
	}
	if _, err = t.conn.Write(b.Bytes()); err != nil {
		return msg, err
	}
	if _, err = t.conn.Write(data); err != nil {
		return msg, err
	}
	var header uint32
	if err = binary.Read(t.conn, binary.BigEndian, &header); err != nil {
		return msg, err
	}
	response := make([]byte, header)
	if err = readMessages(t.conn, response); err != nil {
		return msg, err
	}
	if err = pb.Unmarshal(response, msg); err != nil {
		return msg, err
	}
	if msg.GetOk() != true {
		return msg, errors.New(msg.GetError())
	}
	return msg, nil
}

// UdpClient implementation of SendRecv, will automatically fail if called
func (t *UdpClient) SendRecv(message *proto.Msg) (*proto.Msg, error) {
	return nil, fmt.Errorf("udp doesn't support receiving acknowledgements")
}

// UdpClient implementation of SendMaybeRecv, queues a request to send a message to the server
func (t *UdpClient) SendMaybeRecv(message *proto.Msg) (*proto.Msg, error) {
	response_ch := make(chan response)
	t.requestQueue <- request{message, response_ch}
	r := <-response_ch
	return r.message, r.err
}

// Close will close the UdpClient
func (t *UdpClient) Close() error {
	close(t.requestQueue)
	err := t.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// runRequestQueue services the UdpClient request queue
func (t *UdpClient) runRequestQueue() {
	for req := range t.requestQueue {
		message := req.message
		response_ch := req.response_ch

		msg, err := t.execRequest(message)

		response_ch <- response{msg, err}
	}
}

// execRequest will send a UDP message to Riemann
func (t *UdpClient) execRequest(message *proto.Msg) (*proto.Msg, error) {
	data, err := pb.Marshal(message)
	if err != nil {
		return nil, err
	}
	if len(data) > MAX_UDP_SIZE {
		return nil, fmt.Errorf("unable to send message, too large for udp")
	}
	if _, err = t.conn.Write(data); err != nil {
		return nil, err
	}
	return nil, nil
}

// readMessages will read Riemann messages from the TCP connection
func readMessages(r io.Reader, p []byte) error {
	for len(p) > 0 {
		n, err := r.Read(p)
		p = p[n:]
		if err != nil {
			return err
		}
	}
	return nil
}
