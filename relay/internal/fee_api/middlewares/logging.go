package middlewares

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type (
	// struct for holding response details
	responseData struct {
		status int
		size   int
	}

	// our http.ResponseWriter implementation
	loggingResponseWriter struct {
		http.ResponseWriter // compose original http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
	r.responseData.size += size            // capture size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
	r.responseData.status = statusCode       // capture status code
}

func LoggingMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		loggingFn := func(rw http.ResponseWriter, req *http.Request) {
			start := time.Now()

			responseData := &responseData{
				status: 0,
				size:   0,
			}
			lrw := loggingResponseWriter{
				ResponseWriter: rw, // compose original http.ResponseWriter
				responseData:   responseData,
			}
			next.ServeHTTP(&lrw, req) // inject our implementation of http.ResponseWriter

			duration := time.Since(start)

			logger.Info().Fields(map[string]interface{}{
				"uri":      req.RequestURI,
				"method":   req.Method,
				"status":   responseData.status,
				"duration": duration,
				"size":     responseData.size,
			}).Msg("request completed")
		}
		return http.HandlerFunc(loggingFn)
	}
}
