package logger

import "log"

// Info loguea mensajes informativos.
func Info(msg string) { log.Printf("INFO: %s", msg) }

// Error loguea mensajes de error.
func Error(msg string) { log.Printf("ERROR: %s", msg) }
