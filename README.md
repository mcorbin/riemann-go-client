# Riemann client (Golang)

## Introduction

Go client library for [Riemann](https://github.com/riemann/riemann).

Features:
* Idiomatic concurrency
* Sending events, queries.
* Support tcp and udp client.

This client is a fork of Goryman, a Riemann go client written by Christopher Gilbert. Thanks to him ! The initial Goryman repository (https://github.com/bigdatadev/goryman) has been deleted. We used @rikatz fork (https://github.com/rikatz/goryman/) to create this repository.

We renamed the package name of the client `riemanngo` instead of `goryman`

## Installation

To install the package for use in your own programs:

```
go get github.com/riemann/riemann-go-client
```

If you're a developer, Riemann uses [Google Protocol Buffers](https://github.com/golang/protobuf), so make sure that's installed and available on your PATH.

```
go get github.com/golang/protobuf/{proto,protoc-gen-go}
```

## Getting Started

First we'll need to import the library:

```go
import (
    "github.com/riemann/riemann-go-client/"
)
```

Next we'll need to establish a new client using Connect. The parameter is the connection timeout duration. You can use a TCP client:

```go
c := riemanngo.NewTcpClient("127.0.0.1:5555")
err := c.Connect(5)
if err != nil {
    panic(err)
}
```

Or a UDP client:
```go
c := riemanngo.NewUdpClient("127.0.0.1:5555")
err := c.Connect(5)
if err != nil {
    panic(err)
}
```

Don't forget to close the client connection when you're done:

```go
defer c.Close()
```

Sending events is easy ([list of valid event properties](http://riemann.io/concepts.html)):

```go
result, err := riemanngo.SendEvent(c, &riemanngo.Event{
		Service: "hello",
		Metric:  100,
		Tags: []string{"riemann ftw"},
	})
```
The Hostname and Time in events will automatically be replaced with the hostname of the server and the current time if none is specified.

You can also send batch of events:

```go
events = []riemanngo.Event {
    riemanngo.Event{
        Service: "hello",
        Metric:  100,
        Tags: []string{"hello"},
    },
riemanngo.Event{
        Service: "goodbye",
        Metric:  200,
        Tags: []string{"goodbye"},
    },
}
```

You can also query the Riemann index (using the TCP client):

```go
events, err := QueryEvents(c, "service = \"hello\"")
if err != nil {
    panic(err)
}
```


## Tests

You can lauch Unit tests using

```go
go test -run Unit
```

and Integration tests using:

```go
go test -run Integ
```

## Copyright

See [LICENSE](LICENSE) document
