package infrastructure

import (
"log"
"time"

"go.uber.org/zap"
"go.uber.org/zap/zapcore"
)

type LogService interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Critical(args ...interface{})
}

type LogLevel string

const (
	InfoLogLevel = "info"
	DebugLogLevel = "debug"
	CriticalLogLevel = "critical"
)

type zapFactory struct { }

func newZapFactory() *zapFactory {
	return &zapFactory{}
}

func (*zapFactory) getConfig(level zapcore.Level) zap.Config {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = func(t time.Time, e zapcore.PrimitiveArrayEncoder) {
		e.AppendString(t.Format("15:04:05.000"))
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		DisableCaller:    true,
	}

	return config
}

func (*zapFactory) getLevel(level LogLevel) zapcore.Level {
	if level == InfoLogLevel {
		return zap.InfoLevel
	} else if level == CriticalLogLevel {
		return zap.ErrorLevel
	} else if level == DebugLogLevel {
		return zap.DebugLevel
	}
	return zap.InfoLevel
}

func (z *zapFactory) Build(level LogLevel) (*zap.SugaredLogger, error) {
	zapLevel := z.getLevel(level)
	conf := z.getConfig(zapLevel)
	logger, err := conf.Build()
	if err != nil {
		log.Fatalf("could not build logger %s", err)
	}
	return logger.Sugar(), nil
}

type logService struct {
	logger   *zap.SugaredLogger
}

func NewLogService(logLevel LogLevel) (*logService, error){
	zapFact := newZapFactory()
	logger, err := zapFact.Build(logLevel)
	if err != nil {
		return nil, err
	}
	return &logService{logger: logger}, nil
}

func (l *logService) Info(args ...interface{}) {
	l.logger.Info(args)
}

func (l *logService) Debug(args ...interface{}) {
	l.logger.Debug(args)
}

func (l *logService) Critical(args ...interface{}) {
	l.logger.Error(args)
}

