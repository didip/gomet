[![GoDoc](https://godoc.org/github.com/didip/gomet?status.svg)](http://godoc.org/github.com/didip/gomet)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/didip/gomet/master/LICENSE)

## Gomet

Simple HTTP client & server long poll library for Go.


## Five Minute Tutorial

See: https://github.com/didip/gomet/tree/master/_examples

```
cd _examples
go run server.go &
go run client.go &

curl -H "Content-Type: application/json" -X POST -d '{"message":"hello world"}' http://localhost:8080/
```

## Features

1. Server-side broadcaster conforms to regular HTTP request handler. This means it can be part of middleware stacks.

2. Broadcaster is able to send raw binary since every payload is base64 encoded.

3. Client automatically reconnects when disconnected.

4. Client's HTTP client and HTTP request are customizable.

    ```go
    lp, err := gomet.NewClient("https://localhost:8080/stream")

    lp.HTTPClient.Transport = &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    req, _ := http.NewRequest("GET", "https://localhost:9090/stream", nil)
    lp.HTTPRequest = req
    ```

5. Define your own error handling for client-side.

    ```go
    lp, err := gomet.NewClient("http://localhost:8080/stream")

    lp.SetOnConnectError(func(err error) {
        println("ERROR: Failed to maintain HTTP connection. URL: " + lp.HTTPRequest.URL.String())
    })
    lp.SetOnReadBytesError(func(err error) {
        println("ERROR: Failed to read payload: " + err.Error())
    })
    lp.SetOnBase64DecodeError(func(err error) {
        println("ERROR: Failed to read payload: " + err.Error())
    })
    lp.SetOnPayloadReceived(func(payloadBytes []byte) {
        println("INFO: Received a line: " + string(payloadBytes))
    })
    ```

## My other Go libraries

* [Tollbooth](https://github.com/didip/tollbooth): Simple middleware to rate-limit HTTP requests.

* [Stopwatch](https://github.com/didip/stopwatch): A small library to measure latency of things. Useful if you want to report latency data to Graphite.

* [LaborUnion](https://github.com/didip/laborunion): A worker pool library. It is able to dynamically resize the number of workers.
