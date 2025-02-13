package sse

import (
	"fmt"
	"net/http"
	"skripsi/constant"
	"skripsi/helper"
	"sync"
	"time"
)

type SSEHandler struct {
	j helper.JWTInterface
}

func NewSSEHandler(jwtHelper helper.JWTInterface) *SSEHandler {
	return &SSEHandler{j: jwtHelper}
}

var clients = make(map[string]chan string) // Key: UserID, Value: Channel untuk mengirim event
var clientsMutex sync.Mutex                // Mutex untuk mencegah race condition

func (h *SSEHandler) SseHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		http.Error(w, "Missing or Invalid Authorization Header", http.StatusUnauthorized)
		return
	}

	token, err := h.j.ValidateToken(r.Context(), authHeader)
	if err != nil {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}

	claims := h.j.ExtractUserToken(token)
	if claims == nil {
		http.Error(w, "Token Expired or Invalid Claims", http.StatusUnauthorized)
		return
	}
	userID := claims[constant.JWT_ID].(string)
	if userID == "" {
		http.Error(w, "User ID not found in token", http.StatusUnauthorized)
		return
	}

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Enable CORS locally
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel for the client
	messageChan := make(chan string, 10)

	// Store the channel in the clients map
	clientsMutex.Lock()
	clients[userID] = messageChan
	clientsMutex.Unlock()

	// Ensure the client channel is cleaned up when disconnected
	defer func() {
		clientsMutex.Lock()
		delete(clients, userID)
		clientsMutex.Unlock()
		close(messageChan)
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send periodic heartbeat messages to keep the connection alive
	go func() {
		for {
			select {
			case <-r.Context().Done():
				return
			case <-time.After(15 * time.Second): // Heartbeat setiap 15 detik
				_, err := fmt.Fprintf(w, "data: heartbeat\n\n")
				if err != nil {
					return
				}
				flusher.Flush()
			}
		}
	}()

	// Listen for server events and forward to the client
	for {
		select {
		case msg := <-messageChan:
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				return
			}
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func SendSSENotification(userID, tanggal, jamMulai, jamAkhir string) {
	clientsMutex.Lock()
	clientChan, exists := clients[userID]
	clientsMutex.Unlock()

	if exists {
		message := fmt.Sprintf("Jadwal baru tersedia pada %s pukul %s sampai %s", tanggal, jamMulai, jamAkhir)
		clientChan <- message
	}
}
