package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetDefaultZap sets default zap logger.
func SetDefaultZap() {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "msg",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// replace the zap.L and zap.S default loggers with configured one
	logger, err := config.Build()
	if err != nil {
		zap.L().Panic("couldn't build zap config", zap.Error(err))
	}

	zap.ReplaceGlobals(logger)
}

// Configure configures zap component.
func Configure(debug bool) {
	var config zap.Config
	if debug {
		config = zap.NewDevelopmentConfig()
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stdout"}
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
		// disable log sampling (it drops repeated log entries) since most of the time we don't want to skip logs.
		config.Sampling = nil
		// disable stack trace on every log level to avoid being too verbose on production systems
		config.DisableStacktrace = true
	}

	// replace the zap.L and zap.S default loggers with configured one
	logger, err := config.Build()
	if err != nil {
		zap.L().Panic("fatal error", zap.Error(err))
	}

	zap.ReplaceGlobals(logger)
}
