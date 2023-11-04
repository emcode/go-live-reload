package main

import (
	"html/template"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {

	slog.Info("Starting example server...")
	liveReloadVersionToken := strconv.FormatInt(int64(rand.Int()), 10)

	http.HandleFunc("/", index)
	http.HandleFunc("/live-reload", CreateLiveReloadHandler(liveReloadVersionToken))

	ticker := time.NewTicker(time.Second * 10)
	go runTickerTick(ticker, liveReloadVersionToken)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Could not start to listen on port: 8080")
		os.Exit(1)
	}
}

func runTickerTick(ticker *time.Ticker, message string) {
	for _ = range ticker.C {
		BroadcastLiveReloadMessage(message)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	HandleError(w, err)
}

func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(err.Error())
		os.Exit(1)
	}
}
