package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipResponseWriter wraps http.ResponseWriter to provide gzip compression
type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write method implementation for GzipResponseWriter
func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipMiddleware applies gzip compression to the response
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the client accepts gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Set the response header to indicate gzip compression
		w.Header().Set("Content-Encoding", "gzip")

		// Wrap the ResponseWriter with GzipResponseWriter
		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzrw := GzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzrw, r)
	})
}
