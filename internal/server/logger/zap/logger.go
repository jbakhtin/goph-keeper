package zap

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ToDo: рефакторинг
// ToDo: вынести объект Writer в отдельные реализации, что бы можно было подключать разные хранилища логов
// ToDo: разобраться почему некорректно работает форматирования вывода
// ToDo: разобраться почему файл обновляется только после остановки приложения

const (
	DevelopmentEnvironment = "development"
	ProductionEnvironment  = "production"
)

type Config interface {
	GetAppEnv() string
	GetLoggerFileDirectory() string
	GetLoggerFileMaxSize() int
	GetLoggerFileMaxBackups() int
	GetLoggerFileMaxAge() int
	GetLoggerFileCompress() bool
}

type Logger struct {
	zap.Logger
}

func NewLogger(cfg Config) (lgr *Logger, err error) {
	var tops []teeOption

	switch cfg.GetAppEnv() {
	case DevelopmentEnvironment:
		tops = append(tops, teeOption{
			W: os.Stdout,
			Lef: func(lvl zapcore.Level) bool {
				return true
			},
		})
	case ProductionEnvironment:
		infoLevel, err := setUpLogLevel(cfg, "info", func(lvl zapcore.Level) bool {
			return lvl <= zapcore.InfoLevel
		})
		if err != nil {
			return nil, err
		}

		errorLevel, err := setUpLogLevel(cfg, "error", func(lvl zapcore.Level) bool {
			return lvl > zapcore.InfoLevel
		})
		if err != nil {
			return nil, err
		}

		tops = append(tops, *infoLevel, *errorLevel)
	}

	cores := newTee(tops)
	lgr = &Logger{
		Logger: *zap.New(zapcore.NewTee(cores...)),
	}
	defer func() {
		syncError := lgr.Sync()

		// ToDo: нужно разобратьяс в чем заключается различия sync под osx и linux
		if cfg.GetAppEnv() != DevelopmentEnvironment {
			err = syncError
		}
	}()

	return lgr, err
}

func setUpLogLevel(cfg Config, levelName string, levelCond LevelEnablerFunc) (*teeOption, error) {
	dir := fmt.Sprintf("%v%v/", cfg.GetLoggerFileDirectory(), levelName)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return nil, err
	}

	nowTime := time.Now()
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%v/%v.log", dir, nowTime.Format("2006-01-02")),
		MaxSize:    cfg.GetLoggerFileMaxSize(),
		MaxBackups: cfg.GetLoggerFileMaxBackups(),
		MaxAge:     cfg.GetLoggerFileMaxAge(),
		Compress:   cfg.GetLoggerFileCompress(),
	})

	return &teeOption{
		W:   writer,
		Lef: levelCond,
	}, nil
}

func (l Logger) Debug(msg string, fields ...any) {
	l.Logger.Debug(msg, zap.Any("args", fields))
}

func (l Logger) Info(msg string, fields ...any) {
	l.Logger.Info(msg, zap.Any("args", fields))
}

func (l Logger) Warn(msg string, fields ...any) {
	l.Logger.Warn(msg, zap.Any("args", fields))
}

func (l Logger) Error(msg string, fields ...any) {
	l.Logger.Error(msg, zap.Any("args", fields))
}

func (l Logger) Fatal(msg string, fields ...any) {
	l.Logger.Fatal(msg, zap.Any("args", fields))
}
