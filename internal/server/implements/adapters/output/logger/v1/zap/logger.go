package zap

import (
	"fmt"
	"os"
	"time"

	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/input/config/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/output/logger/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ToDo: нужно рефакторить
// ToDo: вынести объект Writer в отдельные реализации, что бы можно было подключать разные хранилища логов
// ToDo: разобраться почему некорректно работает форматирования вывода
// ToDo: разобраться почему файл синхронится только после остановки приложения

var _ logger.Interface = &Logger{}

const (
	DevelopmentEnvironment = "development"
	ProductionEnvironment  = "production"
)

type Logger struct {
	zap.Logger
}

func NewLogger(cfg config.Interface) (lgr *Logger, err error) {
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
		err = lgr.Sync()
	}()

	return lgr, err
}

func setUpLogLevel(cfg config.Interface, levelName string, levelCond LevelEnablerFunc) (*teeOption, error) {
	dir := fmt.Sprintf("%v%v/", cfg.GetLoggerFileDirectory(), levelName)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return nil, err
	}

	time := time.Now()
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%v/%v.log", dir, time.Format("2006-01-02")),
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

func (z Logger) Debug(msg string, fields ...any) {
	z.Logger.Debug(msg, zap.Any("args", fields))
}

func (z Logger) Info(msg string, fields ...any) {
	z.Logger.Info(msg, zap.Any("args", fields))
}

func (z Logger) Warn(msg string, fields ...any) {
	z.Logger.Warn(msg, zap.Any("args", fields))
}

func (z Logger) Error(msg string, fields ...any) {
	z.Logger.Error(msg, zap.Any("args", fields))
}

func (z Logger) Fatal(msg string, fields ...any) {
	z.Logger.Fatal(msg, zap.Any("args", fields))
}
