package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("proto", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", time.Since(startTime)).
		Msg("Request from gRPC")
	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode  int
	reponseBody []byte
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.reponseBody = b
	return r.ResponseWriter.Write(b)
}

func HttpLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			responseRecorder := &ResponseRecorder{
				ResponseWriter: w, // compose original http.ResponseWriter
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(responseRecorder, r)
			duration := time.Since(startTime)

			logger := log.Info()
			if responseRecorder.statusCode != http.StatusOK {
				logger = log.Error().Bytes("response_body", responseRecorder.reponseBody)
			}

			logger.Str("proto", "http").
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status_code", responseRecorder.statusCode).
				Str("status_text", http.StatusText(responseRecorder.statusCode)).
				Dur("duration", duration).
				Msg("Request from http")
		},
	)
}
