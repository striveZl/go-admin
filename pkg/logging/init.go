package logging

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

type LoggerConfig struct {
	Debug      bool
	Level      string // debug/info/warn/error/dpanic/panic/fatal
	CallerSkip int
	File       struct {
		Enable     bool
		Path       string
		MaxSize    int
		MaxBackups int
	}
	Hooks []*HookConfig
}

type HookConfig struct {
	Enable    bool
	Level     string
	Type      string // gorm
	MaxBuffer int
	MaxThread int
	Options   map[string]string
	Extra     map[string]string
}

func InitWithConfig(ctx context.Context, cfg *LoggerConfig) error {
	var zconfig zap.Config
	if cfg.Debug {
		cfg.Level = "debug"
		zconfig = zap.NewDevelopmentConfig()

	} else {
		zconfig = zap.NewProductionConfig()
	}

	level, err := zapcore.ParseLevel(cfg.Level)

	if err != nil {
		return err
	}

	zconfig.Level.SetLevel(level)

	var (
		logger   *zap.Logger
		cleanFns []func()
	)

	if cfg.File.Enable {
		filename := cfg.File.Path
		_ = os.MkdirAll(filepath.Dir(filename), 0777)
		fileWriter := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    cfg.File.MaxSize,
			MaxBackups: cfg.File.MaxBackups,
			Compress:   false,
			LocalTime:  true,
		}

		cleanFns = append(cleanFns, func() {
			_ = fileWriter.Close()
		})

		zc := zapcore.NewCore(
			zapcore.NewJSONEncoder(zconfig.EncoderConfig),
			zapcore.AddSync(fileWriter),
			zconfig.Level,
		)
		logger = zap.New(zc)
	} else {
		ilogger, err := zconfig.Build()
		if err != nil {
			return err
		}
		logger = ilogger
	}

	skip := cfg.CallerSkip
	if skip <= 0 {
		skip = 2
	}

	logger = logger.WithOptions(
		zap.WithCaller(true),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(skip),
	)

	zap.ReplaceGlobals(logger)

	return nil
}
