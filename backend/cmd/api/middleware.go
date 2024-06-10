package main

import (
	"fmt"
	"net/http"
)

type wrappedResponseWriter struct {
	wrapped       http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func newWrappedResponseWriter(w http.ResponseWriter) *wrappedResponseWriter {
	return &wrappedResponseWriter{
		wrapped:    w,
		statusCode: http.StatusOK,
	}
}

func (mw *wrappedResponseWriter) Header() http.Header {
	return mw.wrapped.Header()
}

func (mw *wrappedResponseWriter) WriteHeader(statusCode int) {
	mw.wrapped.WriteHeader(statusCode)

	if !mw.headerWritten {
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
}

func (mw *wrappedResponseWriter) Write(b []byte) (int, error) {
	mw.headerWritten = true
	return mw.wrapped.Write(b)
}

func (mw *wrappedResponseWriter) Unwrap() http.ResponseWriter {
	return mw.wrapped
}

func (app *App) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := newWrappedResponseWriter(w)
		next.ServeHTTP(mw, r)

		app.logger.Info("request", "uri", r.URL.Path, "statusCode", mw.statusCode)
	})
}

func (app *App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
