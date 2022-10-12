package main

import (
	"context"
	"net/http"

	"github.com/SamHennessy/hlive"
	"github.com/rs/zerolog"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type PageServerWails struct {
	pageFN   func(ctx context.Context) *hlive.Page
	page     *hlive.Page
	logger   zerolog.Logger
	ctxWails context.Context
}

func NewPageServerWails(pageFN func(ctxWails context.Context) *hlive.Page) *PageServerWails {
	newPageFN := func(ctxWails context.Context) *hlive.Page {
		page := pageFN(ctxWails)
		page.DOM().Head().Add(
			hlive.T("script", hlive.Attrs{"src": "/wails/ipc.js"}),
			hlive.T("script", hlive.Attrs{"src": "/wails/runtime.js"}),
		)

		return page
	}

	return &PageServerWails{
		pageFN: newPageFN,
		// TODO: get logger from page
		logger: hlive.LoggerDev.Level(zerolog.DebugLevel),
	}
}

func (ps *PageServerWails) Serve(ctxWails context.Context) {
	ps.ctxWails = ctxWails

	hlive.LoggerDev.Debug().Msg("Serve")

	runtime.EventsOn(ctxWails, "connect", func(optionalData ...interface{}) {
		ps.logger.Debug().Msg("Connect")

		if ps.page == nil {
			ps.page = ps.pageFN(ctxWails)
		}

		fromWails := make(chan hlive.MessageWS)
		toWails := make(chan hlive.MessageWS)

		runtime.EventsOn(ctxWails, "out", func(optionalData ...interface{}) {
			var message []byte

			if len(optionalData) > 0 {
				str, ok := optionalData[0].(string)
				if ok {
					message = []byte(str)
				}
			}

			// TODO: file upload
			fromWails <- hlive.MessageWS{
				Message:  message,
				IsBinary: false,
			}
		})

		go func() {
			for {
				select {
				case <-ctxWails.Done():
					return
				case msg := <-toWails:
					runtime.EventsEmit(ctxWails, "in", string(msg.Message))
				}
			}
		}()

		go func() {
			err := ps.page.ServeWS(ctxWails, "wails", toWails, fromWails)
			if err != nil {
				hlive.LoggerDev.Err(err).Msg("ServeWS")
			}

			hlive.LoggerDev.Info().Msg("ServeWS: done")
		}()
	})
}

func (ps *PageServerWails) AssetsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/hlive.html" {
			if ps.page == nil {
				ps.page = ps.pageFN(r.Context())
			}

			if _, err := ps.page.RunRenderPipeline(r.Context(), w); err != nil {
				hlive.LoggerDev.Err(err).Msg("pageFN.RunRenderPipeline")

				return
			}
		}
	})
}

func (ps *PageServerWails) OnStartup(ctx context.Context) {
	go ps.Serve(ctx)
}
