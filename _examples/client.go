package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/didip/gomet"
)

func main() {
	lp, err := gomet.NewClient("http://localhost:8080/stream")
	if err != nil {
		log.Printf("ERROR: Unable to connect to localhost:8080: %v", err.Error())
	}

	lp.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	lp.SetOnConnectError(func(err error) {
		log.Printf("ERROR: Failed to maintain HTTP connection. URL: %v", lp.HTTPRequest.URL.String())
	})
	lp.SetOnReadBytesError(func(err error) {
		log.Printf("ERROR: Failed to read payload: %v", err.Error())
	})
	lp.SetOnBase64DecodeError(func(err error) {
		log.Printf("ERROR: Failed to read payload: %v", err.Error())
	})
	lp.SetOnPayloadReceived(func(payloadBytes []byte) {
		log.Printf("INFO: Received a line: %v", string(payloadBytes))
	})

	go func() {
		log.Printf("INFO: Spawning a worker to consume from response channel")

		for dataBytes := range lp.ResponseChan {
			log.Printf("INFO: Received data from long poll response channel: %v", string(dataBytes))
		}
	}()

	log.Printf("INFO: Running a daemon that long polls to localhost:8080")

	lp.ConnectForever()
}
