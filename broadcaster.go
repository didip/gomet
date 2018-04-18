package gomet

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type BroadcastWriter struct {
	ID     int64
	Writer http.ResponseWriter

	mtx sync.RWMutex
}

func NewBroadcaster() *Broadcaster {
	broadcaster := &Broadcaster{
		InChan:              make(chan []byte),
		WorkerInChans:       new(sync.Map),
		HTTPResponseHeaders: make(map[string]string),
	}

	return broadcaster
}

type Broadcaster struct {
	InChan              chan []byte
	WorkerInChans       *sync.Map
	HTTPResponseHeaders map[string]string

	onFlusherCastError       func()
	onCloseNotifierCastError func()

	mtx sync.RWMutex
}

func (broadcaster *Broadcaster) NewWorkerInChan() (int64, chan []byte) {
	id := time.Now().UnixNano()
	inChan := make(chan []byte)
	broadcaster.WorkerInChans.Store(id, inChan)

	return id, inChan
}

func (broadcaster *Broadcaster) DeleteWorkerInChan(id int64) {
	broadcaster.WorkerInChans.Delete(id)
}

func (broadcaster *Broadcaster) OnFlusherCastError(f func()) {
	broadcaster.mtx.Lock()
	broadcaster.onFlusherCastError = f
	broadcaster.mtx.Unlock()
}

func (broadcaster *Broadcaster) OnCloseNotifierCastError(f func()) {
	broadcaster.mtx.Lock()
	broadcaster.onCloseNotifierCastError = f
	broadcaster.mtx.Unlock()
}

// Broadcast payload from the main InChan to all of workers' InChan
func (broadcaster *Broadcaster) Broadcast() {
	for payload := range broadcaster.InChan {
		broadcaster.WorkerInChans.Range(func(_, iWorkerInChan interface{}) bool {
			workerInChan := iWorkerInChan.(chan []byte)
			workerInChan <- payload
			return true
		})
	}
}

func (broadcaster *Broadcaster) HTTPHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Connection", "keep-alive")

		// Don't cache response:
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")                                         // Proxies.

		flusher, ok := w.(http.Flusher)
		if !ok {
			if broadcaster.onFlusherCastError != nil {
				broadcaster.onFlusherCastError()
			}
			http.Error(w, `{"Error": "Streaming unsupported"}`, http.StatusInternalServerError)
			return
		}

		closeNotifier, ok := w.(http.CloseNotifier)
		if !ok {
			if broadcaster.onCloseNotifierCastError != nil {
				broadcaster.onCloseNotifierCastError()
			}
			http.Error(w, `{"Error": "Streaming unsupported"}`, http.StatusInternalServerError)
			return
		}

		for key, value := range broadcaster.HTTPResponseHeaders {
			w.Header().Set(key, value)
		}

		workerID, workerInChan := broadcaster.NewWorkerInChan()

		for {
			select {
			case <-closeNotifier.CloseNotify():
				broadcaster.DeleteWorkerInChan(workerID)
				close(workerInChan)
				return

			case payloadBytes := <-workerInChan:
				encodedPayloadString := base64.StdEncoding.EncodeToString(payloadBytes)
				fmt.Fprintf(w, "%s\n", encodedPayloadString)
				flusher.Flush()
			}
		}
	}
}
