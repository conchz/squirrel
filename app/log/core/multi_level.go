package core

import "go.uber.org/zap/zapcore"

func NewMultiLevel(base zapcore.Core) *MultiLevel {
	return &MultiLevel{
		Core:  base,
		cores: map[zapcore.Level]zapcore.Core{},
	}
}

type MultiLevel struct {
	zapcore.Core
	cores map[zapcore.Level]zapcore.Core
}

func (m *MultiLevel) SetCore(level zapcore.Level, core zapcore.Core) {
	m.cores[level] = core
}

func (m *MultiLevel) With(fields []zapcore.Field) zapcore.Core {
	clone := m.clone()
	clone.Core = m.Core.With(fields)
	for level, core := range m.cores {
		clone.cores[level] = core.With(fields)
	}
	return clone
}

func (m *MultiLevel) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if core := m.cores[ent.Level]; core != nil {
		return core.Check(ent, ce)
	}
	return m.Core.Check(ent, ce)
}

func (m *MultiLevel) Sync() error {
	for _, core := range m.cores {
		if err := core.Sync(); err != nil {
			return err
		}
	}
	return m.Core.Sync()
}

func (m *MultiLevel) clone() *MultiLevel {
	newCores := make(map[zapcore.Level]zapcore.Core, len(m.cores))
	for level, core := range m.cores {
		newCores[level] = core
	}
	return &MultiLevel{
		Core:  m.Core,
		cores: newCores,
	}
}
