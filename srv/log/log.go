package log

import (
	"context"
	"log/slog"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

type Logger struct {
	c Config
	*slog.Logger
}

func New(c Config) *Logger {
	lvl := slog.LevelInfo
	if c.Debug {
		lvl = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: lvl,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case "protocol", "grpc.start_time", "peer.address":
				return slog.Attr{}
			case "grpc.method":
				a.Key = "method"
			case "grpc.code":
				a.Key = "status"
			case "grpc.time_ms":
				return slog.String("responseTime", a.Value.String()+"ms")
			case "grpc.request.deadline":
				a.Key = "requestDeadline"
			case "grpc.error":
				a.Key = "error"
			}
			return a
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	if c.UseJSON {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	return &Logger{
		Logger: logger,
		c:      c,
	}
}

func (l Logger) NewUnaryInterceptor() grpc.UnaryServerInterceptor {
	ev := []logging.LoggableEvent{logging.FinishCall}
	if l.c.Debug {
		ev = append(ev, logging.PayloadReceived, logging.PayloadSent)
	}

	return logging.UnaryServerInterceptor(
		l.gRPCLogger(),
		logging.WithLogOnEvents(ev...),
		logging.WithDisableLoggingFields(
			logging.ServiceFieldKey, logging.MethodTypeFieldKey,
			logging.ComponentFieldKey,
		),
	)
}

func (l Logger) NewStreamInterceptor() grpc.StreamServerInterceptor {
	ev := []logging.LoggableEvent{logging.StartCall, logging.FinishCall}
	if l.c.Debug {
		ev = append(ev, logging.PayloadReceived, logging.PayloadSent)
	}

	return logging.StreamServerInterceptor(
		l.gRPCLogger(),
		logging.WithLogOnEvents(ev...),
		logging.WithDisableLoggingFields(
			logging.ServiceFieldKey, logging.MethodTypeFieldKey,
			logging.ComponentFieldKey,
		),
	)
}

func (l Logger) gRPCLogger() logging.LoggerFunc {
	return logging.LoggerFunc(func(
		ctx context.Context,
		lvl logging.Level,
		msg string,
		fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
