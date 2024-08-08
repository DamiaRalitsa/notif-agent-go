package log

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/DamiaRalitsa/notif-agent-go/config"

	log "github.com/sirupsen/logrus"
)

type Log struct {
	appName  string
	logLevel log.Level
	// logstash *logstash.Logstash
}

var (
	logger *Log
	once   sync.Once
)

func initLogger() {
	cfg := config.GetConfig()

	level, ok := map[string]log.Level{
		"ERROR": log.ErrorLevel,
		"DEBUG": log.DebugLevel,
		"INFO":  log.InfoLevel,
	}[cfg.LogLevel]
	if !ok {
		panic(fmt.Sprintf("invalid log level: %s", cfg.LogLevel))
	}

	logger = &Log{
		appName:  cfg.AppName,
		logLevel: level,
		// logstash: logstash.New(cfg.LogstashHost, cfg.LogstashPort, 10, logstash.NetDialer{}),
	}

	// if err := logger.logstash.Connect(); err != nil {
	// 	log.Println("Unable to connect to Logstash:", err)
	// }
}

func GetLogger() *Log {
	once.Do(initLogger)
	return logger
}

func (l *Log) log(level log.Level, fields log.Fields) {
	if l.logLevel >= level {
		_, file, line, _ := runtime.Caller(2)
		message := fmt.Sprintf("[%s][%s:%d][%s] %v: %v, Scope: %v, Meta: %v", time.Now().Format(time.RFC3339), file, line, fields["context"], fields["level"], fields["message"], fields["scope"], fields["meta"])
		// err := l.logstash.Writeln(message)
		// if err != nil {
		// 	log.Println("Unable to send log to Logstash:", err)
		// }

		switch level {
		case log.InfoLevel:
			log.Info(message)
		case log.WarnLevel:
			log.Warn(message)
		case log.ErrorLevel:
			log.Error(message)
		default:
			log.Info(message)
		}
	}
}

func (l *Log) LogWithContext(level log.Level, context, message, scope, meta string) {
	fields := log.Fields{
		"serviceName": "MESSAGING_SERVICE",
		"context":     context,
		"message":     message,
		"scope":       scope,
		"meta":        meta,
		"level":       level.String(),
		"label":       l.appName,
	}
	l.log(level, fields)
}

func (l *Log) Info(context, message, scope, meta string) {
	l.LogWithContext(log.InfoLevel, context, message, scope, meta)
}

func (l *Log) Error(context, message, scope, meta string) {
	l.LogWithContext(log.ErrorLevel, context, message, scope, meta)
}

func (l *Log) Slow(context, message, scope, meta string) {
	l.LogWithContext(log.WarnLevel, context, message, scope, meta) // Use Warn for Slow logs
}
