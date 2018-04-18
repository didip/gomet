package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/didip/gomet"
)

// curl command to send data:
// curl -H "Content-Type: application/json" -X POST -d '{"message":"hello world"}' http://localhost:8080/
func main() {
	broadcastTimeout := 1 * time.Second

	lp := gomet.NewBroadcaster()

	go lp.Broadcast()

	http.HandleFunc("/stream", lp.HTTPHandler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payloadBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR: Failed to read payload: %v", err.Error())
			return
		}

		log.Printf("DEBUG: Received a payload from HTTP POST request: %v", string(payloadBytes))

		select {
		case lp.InChan <- payloadBytes:
		case <-time.After(broadcastTimeout): // Timeout
		}
	})

	log.Printf("INFO: Running HTTP streamer on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
