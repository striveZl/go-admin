package logging

import (
	"context"
	"fmt"
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

func InitWithConfig(_ context.Context, cfg *LoggerConfig) error {
	zconfig, level, err := buildZapConfig(cfg)
	if err != nil {
		return err
	}

	logger, err := buildLogger(cfg, zconfig, level)
	if err != nil {
		return err
	}

	logger = logger.WithOptions(loggerOptions(cfg)...)

	zap.ReplaceGlobals(logger)

	return nil
}

func buildZapConfig(cfg *LoggerConfig) (zap.Config, zapcore.Level, error) {
	levelName := cfg.Level

	var zconfig zap.Config
	if cfg.Debug {
		levelName = "debug"
		zconfig = zap.NewDevelopmentConfig()
	} else {
		zconfig = zap.NewProductionConfig()
	}

	level, err := zapcore.ParseLevel(levelName)
	if err != nil {
		return zap.Config{}, zapcore.InfoLevel, err
	}

	zconfig.Level.SetLevel(level)
	return zconfig, level, nil
}

func buildLogger(cfg *LoggerConfig, zconfig zap.Config, level zapcore.Level) (*zap.Logger, error) {
	if !cfg.File.Enable {
		return zconfig.Build()
	}

	writer, err := buildFileWriteSyncer(cfg)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zconfig.EncoderConfig),
		writer,
		zap.NewAtomicLevelAt(level),
	)
	return zap.New(core), nil
}

func buildFileWriteSyncer(cfg *LoggerConfig) (zapcore.WriteSyncer, error) {
	filename := cfg.File.Path
	if filename == "" {
		return nil, fmt.Errorf("logger file path is required")
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0o777); err != nil {
		return nil, err
	}

	fileWriter := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.File.MaxSize,
		MaxBackups: cfg.File.MaxBackups,
		Compress:   false,
		LocalTime:  true,
	}

	return zapcore.AddSync(fileWriter), nil
}

func loggerOptions(cfg *LoggerConfig) []zap.Option {
	skip := cfg.CallerSkip
	if skip <= 0 {
		skip = 2
	}

	return []zap.Option{
		zap.WithCaller(true),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(skip),
	}
}
