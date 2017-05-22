package log

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sort"
	"time"
)

// reference: https://github.com/emonuh/zap-error-dispatcher

type ZapLogConfig struct {
	zap.Config           `json:",inline" yaml:",inline"`
	ErrorDispatcherPaths []string `json:"errorDispatcherPaths" yaml:"errorDispatcherPaths"`
}

var errNoEncoderNameSpecified = errors.New("no encoder name specified")

func New(config ZapLogConfig) (l *zap.SugaredLogger) {
	var err error
	l, err = config.Build()
	if err != nil {
		panic(err)
	}

	return
}

func (c *ZapLogConfig) Build(options ...zap.Option) (*zap.SugaredLogger, error) {
	enc, err := c.buildEncoder()
	if err != nil {
		return nil, err
	}

	sink, errDispSink, errSink, err := c.openSinks()
	if err != nil {
		return nil, err
	}

	baseCore := zapcore.NewCore(enc, sink, c.Level)
	errCore := zapcore.NewCore(enc, errDispSink, c.Level)
	errorDispatcher := NewErrorDispatcher(baseCore, errCore)
	l := zap.New(
		errorDispatcher,
		c.buildOptions(errSink)...,
	)
	if len(options) > 0 {
		l = l.WithOptions(options...)
	}
	return l.Sugar(), nil
}

func (c *ZapLogConfig) buildEncoder() (encoder zapcore.Encoder, err error) {
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

func (c *ZapLogConfig) openSinks() (zapcore.WriteSyncer, zapcore.WriteSyncer, zapcore.WriteSyncer, error) {
	sink, closeOut, err := zap.Open(c.OutputPaths...)
	if err != nil {
		closeOut()
		return nil, nil, nil, err
	}
	errDispSink, closeErrDisp, err := zap.Open(c.ErrorDispatcherPaths...)
	if err != nil {
		closeOut()
		closeErrDisp()
		return nil, nil, nil, err
	}
	errSink, closeErr, err := zap.Open(c.ErrorOutputPaths...)
	if err != nil {
		closeOut()
		closeErrDisp()
		closeErr()
		return nil, nil, nil, err
	}
	return sink, errDispSink, errSink, nil
}

func (c *ZapLogConfig) buildOptions(errSink zapcore.WriteSyncer) []zap.Option {
	options := []zap.Option{zap.ErrorOutput(errSink)}

	if c.Development {
		options = append(options, zap.Development())
	}

	if !c.DisableCaller {
		options = append(options, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if c.Development {
		stackLevel = zap.WarnLevel
	}
	if !c.DisableStacktrace {
		options = append(options, zap.AddStacktrace(stackLevel))
	}

	if c.Sampling != nil {
		options = append(options, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
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
		options = append(options, zap.Fields(fs...))
	}

	return options
}
