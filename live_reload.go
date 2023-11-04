package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

var LIVE_RELOAD_CHANNEL_STORE []*chan string = make([]*chan string, 0)

func removeChannel(ch *chan string) {
	pos := -1
	storeLen := len(LIVE_RELOAD_CHANNEL_STORE)
	for i, msgChan := range LIVE_RELOAD_CHANNEL_STORE {
		if ch == msgChan {
			pos = i
		}
	}

	if pos == -1 {
		return
	}
	LIVE_RELOAD_CHANNEL_STORE[pos] = LIVE_RELOAD_CHANNEL_STORE[storeLen-1]
	LIVE_RELOAD_CHANNEL_STORE = LIVE_RELOAD_CHANNEL_STORE[:storeLen-1]
	slog.Debug("Connection remains", "count", len(LIVE_RELOAD_CHANNEL_STORE))
}

func BroadcastLiveReloadMessage(msg string) {
	for _, ch := range LIVE_RELOAD_CHANNEL_STORE {
		*ch <- msg
	}
}

func CreateLiveReloadHandler(versionToken string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ExecuteLiveReloadHandler(w, r, versionToken)
	}
}

func ExecuteLiveReloadHandler(w http.ResponseWriter, r *http.Request, versionToken string) {
	ch := make(chan string)
	LIVE_RELOAD_CHANNEL_STORE = append(LIVE_RELOAD_CHANNEL_STORE, &ch)

	slog.Debug("Client connected", "count", len(LIVE_RELOAD_CHANNEL_STORE))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	defer func() {
		close(ch)
		removeChannel(&ch)
		slog.Debug("Client closed connection")
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		slog.Error("Could not init http.Flusher")
		os.Exit(1)
	}

	_, err := fmt.Fprintf(w, "data: %s\n\n", versionToken)
	if err != nil {
		slog.Error("Could write SSE message", "error", err.Error())
		os.Exit(1)
	}
	flusher.Flush()

	for {
		select {
		case message := <-ch:
			_, err := fmt.Fprintf(w, "data: %s\n\n", message)
			if err != nil {
				slog.Error("Could write SSE message", "error", err.Error())
				os.Exit(1)
				return
			}
			flusher.Flush()
		case <-r.Context().Done():
			slog.Debug("Client closed connection")
			return
		}
	}
}
