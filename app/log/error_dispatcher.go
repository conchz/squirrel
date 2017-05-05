package log

import "go.uber.org/zap/zapcore"

func NewErrorDispatcher(base zapcore.Core, err zapcore.Core) zapcore.Core {
	return &errorDispatcher{
		Core: base,
		err:  err,
	}
}

type errorDispatcher struct {
	zapcore.Core
	err zapcore.Core
}

func (e *errorDispatcher) With(fields []zapcore.Field) zapcore.Core {
	clone := e.clone()
	clone.Core = e.Core.With(fields)
	clone.err = e.err.With(fields)
	return clone
}

func (e *errorDispatcher) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if ent.Level >= zapcore.ErrorLevel {
		return e.err.Check(ent, ce)
	}
	return e.Core.Check(ent, ce)
}

func (e *errorDispatcher) Sync() error {
	if err := e.err.Sync(); err != nil {
		return err
	}
	return e.Core.Sync()
}

func (e *errorDispatcher) clone() *errorDispatcher {
	return &errorDispatcher{
		Core: e.Core,
		err:  e.err,
	}
}
