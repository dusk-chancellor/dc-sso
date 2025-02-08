package server

import (
	"context"
	"fmt"
	"net"

	adapter "github.com/dusk-chancellor/dc-sso/internal/adapters/grpc"
	"github.com/dusk-chancellor/dc-sso/internal/service"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// grpc server configuration

type Server struct {
	log 	   *zap.SugaredLogger
	grpcServer *grpc.Server
	port 	   int
}

func New(log *zap.SugaredLogger, service service.Service, port int) *Server {
	// log options
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadSent, logging.PayloadReceived,
		),
		logging.WithErrorFields(
			func(err error) logging.Fields {
				return logging.Fields{
					err.Error(),
				}
			},
		),
	}
	// recovery options
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("panic", zap.Any("panic", p))

			return status.Error(codes.Internal, "internal error")
		}),
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(interceptorLogger(log), loggingOpts...),
	))

	adapter.RegisterGrpc(grpcServer, &service)

	return &Server{
		log: log,
		grpcServer: grpcServer,
		port: port,
	}
}
// 
func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		s.log.DPanic("failed running server", zap.Error(err))
	}
}
// runs server
func (s *Server) Run() error {
	childLogger := s.log.With(
		zap.String("operation", "Run()"),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		childLogger.Error("listening process error", zap.Error(err))
		return err
	}

	s.log.Info("grpc server started on", zap.String("address", l.Addr().String()))

	if err := s.grpcServer.Serve(l); err != nil {
		childLogger.Error("serving process error", zap.Error(err))
		return err
	}

	return nil
}
// stops server
func (s *Server) Stop() {
	s.log.Info("stopping server...", zap.Int("port", s.port))

	s.grpcServer.GracefulStop()
}
// convert zap.Logger to logging.Logger
func interceptorLogger(log *zap.SugaredLogger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := log.WithOptions(zap.AddCallerSkip(1)).With(f)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
