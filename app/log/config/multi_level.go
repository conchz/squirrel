package config

import (
	"fmt"
	"github.com/lavenderx/squirrel/app/log/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sort"
	"time"
)

type MultiLevel struct {
	zap.Config `json:",inline" yaml:",inline"`
	LevelPaths map[zapcore.Level][]string `json:"levelPaths" yaml:"levelPaths"`
}

func (c *MultiLevel) Build(opts ...zap.Option) (*zap.Logger, error) {
	enc, err := c.buildEncoder()
	if err != nil {
		return nil, err
	}

	sink, levelSinks, errSink, err := c.openSinks()
	if err != nil {
		return nil, err
	}

	baseCore := zapcore.NewCore(enc, sink, c.Level)
	multiLevelCore := core.NewMultiLevel(baseCore)
	for level, sink := range levelSinks {
		multiLevelCore.SetCore(level, zapcore.NewCore(enc, sink, c.Level))
	}
	log := zap.New(
		multiLevelCore,
		c.buildOptions(errSink)...,
	)
	if len(opts) > 0 {
		log = log.WithOptions(opts...)
	}
	return log, nil
}

func (c *MultiLevel) buildEncoder() (encoder zapcore.Encoder, err error) {
	if len(c.Encoding) == 0 {
		err = errNoEncoderNameSpecified
		return
	}
	switch c.Encoding {
	case "console":
		encoder = zapcore.NewConsoleEncoder(c.EncoderConfig)
	case "json":
		encoder = zapcore.NewJSONEncoder(c.EncoderConfig)
	default:
		err = fmt.Errorf("no encoder registered for name %q", c.Encoding)
	}
	return
}

func (c *MultiLevel) openSinks() (zapcore.WriteSyncer, map[zapcore.Level]zapcore.WriteSyncer, zapcore.WriteSyncer, error) {
	sink, closeOut, err := zap.Open(c.OutputPaths...)
	if err != nil {
		closeOut()
		return nil, nil, nil, err
	}
	levelSinks, err := c.openLevelSinks()
	if err != nil {
		closeOut()
		return nil, nil, nil, err
	}
	errSink, closeErr, err := zap.Open(c.ErrorOutputPaths...)
	if err != nil {
		closeOut()
		closeErr()
		return nil, nil, nil, err
	}
	return sink, levelSinks, errSink, nil
}

func (c *MultiLevel) openLevelSinks() (map[zapcore.Level]zapcore.WriteSyncer, error) {
	levelSinks := make(map[zapcore.Level]zapcore.WriteSyncer, len(c.LevelPaths))
	var closes []func()
	for level, paths := range c.LevelPaths {
		sink, close, err := zap.Open(paths...)
		closes = append(closes, close)
		if err != nil {
			for _, close := range closes {
				close()
			}
			return nil, err
		}
		levelSinks[level] = sink
	}
	return levelSinks, nil
}

func (c *MultiLevel) buildOptions(errSink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errSink)}

	if c.Development {
		opts = append(opts, zap.Development())
	}

	if !c.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if c.Development {
		stackLevel = zap.WarnLevel
	}
	if !c.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if c.Sampling != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSampler(core, time.Second, int(c.Sampling.Initial), int(c.Sampling.Thereafter))
		}))
	}

	if len(c.InitialFields) > 0 {
		fs := make([]zapcore.Field, 0, len(c.InitialFields))
		keys := make([]string, 0, len(c.InitialFields))
		for k := range c.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, c.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	return opts
}
