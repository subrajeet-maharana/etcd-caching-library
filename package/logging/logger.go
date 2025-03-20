package logging

import (
	"log"

	"github.com/natefinch/lumberjack"
)

type Logger struct {
	logger *lumberjack.Logger
}

func NewLogger(filename string, maxSize, maxBackups, maxAge int, compress bool) *Logger {
	logger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,    // megabytes
		MaxBackups: maxBackups, // number of backups
		MaxAge:     maxAge,     // days
		Compress:   compress,   // compress the backups
	}
	log.SetOutput(logger)
	log.Println("Log rotation initialized")
	return &Logger{logger: logger}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	log.Printf("ERROR: "+format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	log.Printf("WARNING: "+format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	log.Printf("DEBUG: "+format, v...)
}
