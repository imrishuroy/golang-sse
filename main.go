package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type sendHiReq struct {
	Message string `json:"message"`
}

// A struct to hold the connections for broadcasting messages
type connection struct {
	w http.ResponseWriter
}

var connections = make(map[*connection]bool)
var mu sync.Mutex

func sendHi(w http.ResponseWriter, r *http.Request) {
	// Handle POST requests to receive the message
	if r.Method == http.MethodPost {
		var req sendHiReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Broadcast the received message to all connected clients
		mu.Lock()
		for conn := range connections {
			_, _ = fmt.Fprintf(conn.w, "data: %s\n\n", req.Message)
			conn.w.(http.Flusher).Flush()
		}
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
		return
	}

	// Handle SSE for GET requests
	if r.Method == http.MethodGet {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		conn := &connection{w: w}
		mu.Lock()
		connections[conn] = true
		mu.Unlock()

		// Keep the connection alive
		for {
			time.Sleep(30 * time.Second)
			//_, _ = fmt.Fprintf(w, "data: %s\n\n", "Keep-alive message")
			flusher.Flush()
		}
	}
}

func main() {
	http.HandleFunc("/messages", sendHi)
	fmt.Println("Server started at 192.168.10.117:8080")
	err := http.ListenAndServe("192.168.10.117:8080", nil)
	if err != nil {
		fmt.Println("Error starting the server", err)
		return
	}
}

//package main
//
//import (
//	"fmt"
//	"net/http"
//	"time"
//)
//
//func main() {
//	fmt.Printf("Server started at http://localhost:8080\n")
//	http.HandleFunc("/events", eventsHandler)
//
//	http.HandleFunc("/messages", sendHi)
//
//	err := http.ListenAndServe(":8080", nil)
//	if err != nil {
//		return
//	}
//
//}
//
//func eventsHandler(w http.ResponseWriter, r *http.Request) {
//	// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
//
//	w.Header().Set("Content-Type", "text/event-stream")
//	w.Header().Set("Cache-Control", "no-cache")
//	w.Header().Set("Connection", "keep-alive")
//
//	// Simulate sending events (you can replace this with real data)
//	for i := 0; i < 200; i++ {
//		_, _ = fmt.Fprintf(w, "data: %s\n\n", fmt.Sprintf("Event %d", i))
//		time.Sleep(2 * time.Second)
//		w.(http.Flusher).Flush()
//	}
//
//	// Simulate closing the connection
//	closeNotify := w.(http.CloseNotifier).CloseNotify()
//	<-closeNotify
//}
//
////type sendHiReq struct {
////	Message string `json:"message"`
////}
//
//func sendHi(w http.ResponseWriter, r *http.Request) {
//
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
//
//	w.Header().Set("Content-Type", "text/event-stream")
//	w.Header().Set("Cache-Control", "no-cache")
//	w.Header().Set("Connection", "keep-alive")
//
//	_, _ = fmt.Fprintf(w, "data: %s\n\n", "Hi")
//	w.(http.Flusher).Flush()
//
//	// Simulate closing the connection
//	closeNotify := w.(http.CloseNotifier).CloseNotify()
//	<-closeNotify
//
//}
//
//// Parse the incoming JSON payload
////var req sendHiReq
////if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
////	http.Error(w, "Invalid request payload", http.StatusBadRequest)
////	return
////}
//
//// Process the message (for example, convert to uppercase)
////message := fmt.Sprintf("Processed Message: %s", req.Message)
