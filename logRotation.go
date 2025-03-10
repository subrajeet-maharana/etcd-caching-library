package main

import (
	"log"

	"github.com/natefinch/lumberjack"
)

func SetupLogRotation() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "events.log",
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	})
	log.Println("Log Rotation intialized")
}
