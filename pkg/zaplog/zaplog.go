package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.SugaredLogger {
	cfg := zap.Config{
		Level: zap.NewDevelopmentConfig().Level,
		Development: true,
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"strerr",
		},
	}

	logger := zap.Must(cfg.Build())

	return logger.Sugar()
}

func Log(msg string) {
	zap.S().Log(zapcore.InfoLevel, msg)
}
