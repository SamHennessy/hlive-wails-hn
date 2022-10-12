package main

import (
	"net/http"

	l "github.com/SamHennessy/hlive"
)

func main() {
	// Server
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("../frontend/dist")))

	// Listen
	l.LoggerDev.Info().Str("addr", ":3000").Msg("HLive Wails Dev server listening")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		l.LoggerDev.Err(err).Msg("http listen and serve")
	}
}
