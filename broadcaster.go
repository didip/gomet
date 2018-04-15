package gomet

import (
	"encoding/base64"
	"io"
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
		Writers:             new(sync.Map),
		HTTPResponseHeaders: make(map[string]string),
	}

	return broadcaster
}

type Broadcaster struct {
	InChan              chan []byte
	Writers             *sync.Map
	HTTPResponseHeaders map[string]string

	mtx sync.RWMutex
}

func (broadcaster *Broadcaster) SetHTTPResponseHeaders(headers map[string]string) {
	broadcaster.mtx.Lock()
	for key, value := range headers {
		broadcaster.HTTPResponseHeaders[key] = value
	}
	broadcaster.mtx.Unlock()
}

func (broadcaster *Broadcaster) NewBroadcastWriter(w http.ResponseWriter) *BroadcastWriter {
	broadcaster.mtx.Lock()
	defer broadcaster.mtx.Unlock()

	id := time.Now().UnixNano()

	bWriter := &BroadcastWriter{
		ID:     id,
		Writer: w,
	}

	broadcaster.Writers.Store(bWriter.ID, bWriter)

	return bWriter
}

func (broadcaster *Broadcaster) DeleteBroadcastWriter(id int64) {
	broadcaster.Writers.Delete(id)
}

func (broadcaster *Broadcaster) BroadcastWorker() {
	for payloadBytes := range broadcaster.InChan {
		broadcaster.Writers.Range(func(_, iBWriter interface{}) bool {
			bWriter := iBWriter.(*BroadcastWriter)

			encodedPayloadString := base64.StdEncoding.EncodeToString(payloadBytes)
			io.WriteString(bWriter.Writer, encodedPayloadString+"\n")

			flusher, ok := bWriter.Writer.(http.Flusher)
			if ok {
				flusher.Flush()
			}

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

		_, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, `{"Error": "Streaming unsupported"}`, http.StatusInternalServerError)
			return
		}

		for key, value := range broadcaster.HTTPResponseHeaders {
			w.Header().Set(key, value)
		}

		bWriter := broadcaster.NewBroadcastWriter(w)

		// Detect closed connection from client
		// When it happened, remove the BroadcastWriter to prevent memory leak
		notify := w.(http.CloseNotifier).CloseNotify()
		<-notify
		if bWriter != nil {
			broadcaster.DeleteBroadcastWriter(bWriter.ID)
		}
	}
}
