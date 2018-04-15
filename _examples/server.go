package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/didip/gomet"
)

// curl command to send data:
// curl -H "Content-Type: application/json" -X POST -d '{"message":"hello world"}' http://localhost:8080/
func main() {
	lp := gomet.NewBroadcaster()

	for i := 0; i < 10; i++ {
		fmt.Printf("INFO: Running Broadcast Worker(%v)\n", i)
		go lp.BroadcastWorker()
	}

	http.HandleFunc("/stream", lp.HTTPHandler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payloadBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			println("ERROR: Failed to read payload: " + err.Error())
			return
		}

		println("INFO: Received a payload from HTTP POST request: " + string(payloadBytes))

		select {
		case lp.InChan <- payloadBytes:
		case <-time.After(1 * time.Second): // Timeout
		}
	})

	println("INFO: Running HTTP streamer on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
