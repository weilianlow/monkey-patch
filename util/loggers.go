package util

import (
	"context"
	"reflect"
	"unsafe"
)

var (
	globalLoggers      = LoggerAdapter{}
	globalIsNoopLogger = true
)

type LoggerAdapter struct {
	ulog   *ulog.Logger
	wlogv2 wsaloggerv2.Logger
	wlog   wsalogger.Logger
}

func (l LoggerAdapter) Enabled(level wsaloggerv2.Level) bool {
	if l.ulog != nil {
		// method n/a
	} else if l.wlogv2 != nil {
		return l.wlogv2.Enabled(level)
	} else if l.wlog != nil {
		return l.wlog.Enabled(wsalogger.Level(level))
	}
	return false
}

func (l LoggerAdapter) Debug(msg string) {
	if l.ulog != nil {
		l.ulog.Debug(msg)
	} else if l.wlogv2 != nil {
		l.wlogv2.Debug(msg)
	} else if l.wlog != nil {
		l.wlog.Debug(msg)
	}
}

func (l LoggerAdapter) Info(msg string) {
	if l.ulog != nil {
		l.ulog.Info(msg)
	} else if l.wlogv2 != nil {
		l.wlogv2.Info(msg)
	} else if l.wlog != nil {
		l.wlog.Info(msg)
	}
}

func (l LoggerAdapter) Data(msg string) {
	if l.ulog != nil {
		l.ulog.Data(msg)
	} else if l.wlogv2 != nil {
		l.wlogv2.Data(msg)
	} else if l.wlog != nil {
		l.wlog.Data(msg)
	}
}

func (l LoggerAdapter) Access(msg string) {
	if l.ulog != nil {
		// method n/a
	} else if l.wlogv2 != nil {
		l.wlogv2.Access(msg)
	} else if l.wlog != nil {
		l.wlog.Access(msg)
	}
}

func (l LoggerAdapter) Warn(msg string) {
	if l.ulog != nil {
		l.ulog.Warn(msg)
	} else if l.wlogv2 != nil {
		l.wlogv2.Warn(msg)
	} else if l.wlog != nil {
		l.wlog.Warn(msg)
	}
}

func (l LoggerAdapter) Error(msg string) {
	if l.ulog != nil {
		l.ulog.Error(msg)
	} else if l.wlogv2 != nil {
		l.wlogv2.Error(msg)
	} else if l.wlog != nil {
		l.wlog.Error(msg)
	}
}

func (l LoggerAdapter) WithField(key string, value interface{}) wsaloggerv2.Logger {
	if l.ulog != nil {
		l.ulog.WithField(key, value)
	} else if l.wlogv2 != nil {
		l.wlogv2.WithField(key, value)
	} else if l.wlog != nil {
		l.wlog.WithField(key, value)
	}
	return l
}

func (l LoggerAdapter) WithFields(fields wsaloggerv2.Fields) wsaloggerv2.Logger {
	if l.ulog != nil {
		l.ulog.WithFields(ulog.Fields(fields))
	} else if l.wlogv2 != nil {
		l.wlogv2.WithFields(fields)
	} else if l.wlog != nil {
		l.wlog.WithFields(wsalogger.Fields(fields))
	}
	return l
}

func (l LoggerAdapter) Withs(args ...interface{}) wsaloggerv2.Logger {
	if l.ulog != nil {
		l.ulog.Withs(args)
	} else if l.wlogv2 != nil {
		l.wlogv2.Withs(args)
	} else if l.wlog != nil {
		l.wlog.Withs(args)
	}
	return l
}

func (l LoggerAdapter) WithTypedField(typedField wsaloggerv2.TypedField) wsaloggerv2.Logger {
	if l.ulog != nil {
		var typedFields []ulog.TypedField
		typedFields = append(typedFields, ulog.Any(typedField.Key, typedField.Value()))
		l.ulog.WithTypedFields(typedFields...)
	} else if l.wlogv2 != nil {
		l.wlogv2.WithTypedField(typedField)
	} else if l.wlog != nil {
		// method n/a
	}
	return l
}

func (l LoggerAdapter) WithTypedFields(typedFielder wsaloggerv2.TypedFielder) wsaloggerv2.Logger {
	if l.ulog != nil {
		var typedFields []ulog.TypedField
		fields := typedFielder.Fields()
		for _, field := range fields {
			typedFields = append(typedFields, ulog.Any(field.Key, field.Value()))
		}
		l.ulog.WithTypedFields(typedFields...)
	} else if l.wlogv2 != nil {
		l.wlogv2.WithTypedFields(typedFielder)
	} else if l.wlog != nil {
		// method n/a
	}
	return l
}

func (l LoggerAdapter) WithSentryField(key string, value interface{}) wsaloggerv2.Logger {
	if l.ulog != nil {
		// method n/a
	} else if l.wlogv2 != nil {
		l.wlogv2.WithSentryField(key, value)
	} else if l.wlog != nil {
		// method n/a
	}
	return l
}

func (l LoggerAdapter) WithSentryFields(args ...interface{}) wsaloggerv2.Logger {
	if l.ulog != nil {
		// method n/a
	} else if l.wlogv2 != nil {
		l.wlogv2.WithSentryFields(args)
	} else if l.wlog != nil {
		// method n/a
	}
	return l
}

func (l LoggerAdapter) WithError(err error) wsaloggerv2.Logger {
	if l.ulog != nil {
		l.ulog.Error(err.Error())
	} else if l.wlogv2 != nil {
		l.wlogv2.WithError(err)
	} else if l.wlog != nil {
		l.wlog.WithError(err)
	}
	return l
}

func (l LoggerAdapter) WithMetaField(key wsaloggerv2.MetaKey, value interface{}) wsaloggerv2.Logger {
	if l.ulog != nil {
		// method n/a
	} else if l.wlogv2 != nil {
		l.wlogv2.WithMetaField(key, value)
	} else if l.wlog != nil {
		l.wlog.WithMetaField(wsalogger.MetaKey(key), value)
	}
	return l
}

func (l LoggerAdapter) GetMeta() wsaloggerv2.Meta {
	var meta wsaloggerv2.Meta
	if l.ulog != nil {
		// method n/a
	} else if l.wlogv2 != nil {
		meta = l.wlogv2.GetMeta()
	} else if l.wlog != nil {
		meta = l.wlog.GetMeta()
	}
	return meta
}

func (l LoggerAdapter) WithSpan(span tracing.Span) wsaloggerv2.Logger {
	if l.ulog != nil {
		// method n/a
	} else if l.wlogv2 != nil {
		l.wlogv2.WithSpan(span)
	} else if l.wlog != nil {
		l.wlog.WithSpan(span)
	}
	return l
}

// logger without context
func InitULog(logger *ulog.Logger) {
	if !IsNoopLogger(logger) {
		globalLoggers.ulog = logger
	}
}

func InitWLogV2(logger wsaloggerv2.Logger) {
	if logger != nil {
		globalLoggers.wlogv2 = logger
	}
}

func InitWLog(logger wsalogger.Logger) {
	if logger != nil {
		globalLoggers.wlog = logger
	}
}

func GetLogger() wsaloggerv2.Logger {
	return globalLoggers
}

// logger with context
func GetLoggerFromContext(ctx context.Context) wsaloggerv2.Logger {
	logger := LoggerAdapter{}
	if l := ulog.DefaultLoggerFromContext(ctx); !IsNoopLogger(l) {
		logger.ulog = l
	} else if l := utils.GetLoggerV2(ctx); l != nil {
		logger.wlogv2 = l
	} else if l := utils.GetLogger(ctx); l != nil {
		logger.wlog = l
	}
	return logger
}

func IsNoopLogger(logger *ulog.Logger) bool {
	if !globalIsNoopLogger {
		return globalIsNoopLogger
	}
	if logger == nil {
		return true
	}
	var originalPtr, TargetPtr interface{} = logger, &ulog.NoopLogger{}
	e := reflect.ValueOf(originalPtr).Elem()
	varName := e.Type().Field(0).Name
	if varName == "baseLogger" {
		if reflect.TypeOf(GetUnexportedField(e.Field(0))) == reflect.TypeOf(TargetPtr) {
			return true
		}
	}
	globalIsNoopLogger = false
	return false
}

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}
