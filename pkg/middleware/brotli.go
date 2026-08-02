package middleware

import (
	"net/http"

	"github.com/andybalholm/brotli"
)

func Brotli(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")

		compressor := brotli.HTTPCompressor(w, r)
		defer compressor.Close()

		cw := &brotliResponseWriter{
			ResponseWriter: w,
			writer:         compressor,
		}

		next.ServeHTTP(cw, r)
	})
}

type brotliResponseWriter struct {
	http.ResponseWriter
	writer interface {
		Write([]byte) (int, error)
	}
}

func (cw *brotliResponseWriter) Write(b []byte) (int, error) {
	return cw.writer.Write(b)
}
