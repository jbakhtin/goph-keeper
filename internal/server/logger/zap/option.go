package zap

import (
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LevelEnablerFunc func(lvl zapcore.Level) bool

type RotateOptions struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

type teeOption struct {
	W   io.Writer
	Lef LevelEnablerFunc
}

func newTee(tops []teeOption) []zapcore.Core {
	var cores []zapcore.Core
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}

	for _, top := range tops {
		top := top
		if top.W == nil {
			ms := "the writer is nil"
			panic(ms)
		}

		lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return top.Lef(zapcore.Level(lvl))
		})

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.AddSync(top.W),
			lv,
		)
		cores = append(cores, core)
	}

	return cores
}
