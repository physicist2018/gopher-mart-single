package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipMiddleware — мидлварь для поддержки сжатия запросов и ответов
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Распаковка запроса, если он сжат
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Unable to decompress request body", http.StatusBadRequest)
				return
			}
			defer gz.Close()
			r.Body = gz
		}

		// Обертка ответа для сжатия
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gzWriter := gzip.NewWriter(w)
			defer gzWriter.Close()

			// Обернем ResponseWriter
			gzResponseWriter := &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gzWriter,
			}
			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gzResponseWriter, r)
			return
		}

		// Если сжатие не поддерживается, просто передаем дальше
		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
