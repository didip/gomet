package main

import (
	"crypto/tls"
	"net/http"

	"github.com/didip/gomet"
)

func main() {
	lp, err := gomet.NewClient("http://localhost:8080/stream")
	if err != nil {
		println("ERROR: Unable to connect to localhost:8080: " + err.Error())
	}

	lp.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

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

	go func() {
		println("INFO: Spawning a worker to consume from response channel")

		for dataBytes := range lp.ResponseChan {
			println("INFO: Received data from long poll response channel: " + string(dataBytes))
		}
	}()

	println("INFO: Running a daemon that long polls to localhost:8080")

	lp.ConnectForever()
}
