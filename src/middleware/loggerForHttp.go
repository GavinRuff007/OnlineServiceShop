package middleware

import (
	"RestGoTest/src/config"
	"RestGoTest/src/pkg/logging"
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func DefaultStructuredLoggerForHttp(cfg *config.Config) func(http.Handler) http.Handler {
	logger := logging.NewLogger(cfg)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "swagger") {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			var reqBody []byte
			if r.Body != nil {
				reqBody, _ = io.ReadAll(r.Body)
				r.Body.Close()
				r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			lrw := &loggingResponseWriter{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(lrw, r)

			duration := time.Since(start)

			keys := map[logging.ExtraKey]interface{}{}
			keys[logging.Path] = r.URL.Path
			keys[logging.ClientIp] = r.RemoteAddr
			keys[logging.Method] = r.Method
			keys[logging.Latency] = duration
			keys[logging.StatusCode] = lrw.statusCode
			keys[logging.ErrorMessage] = "" // در صورت نیاز میشه از context گرفتن
			keys[logging.BodySize] = lrw.body.Len()
			keys[logging.RequestBody] = string(reqBody)
			keys[logging.ResponseBody] = lrw.body.String()

			logger.Info(logging.RequestResponse, logging.Api, "", keys)
		})
	}
}
