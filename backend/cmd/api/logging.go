package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strings"
)

// logLineBridge envía líneas del paquete `log` al logger slog por defecto (misma salida JSON/texto).
type logLineBridge struct{}

func (*logLineBridge) Write(p []byte) (n int, err error) {
	msg := strings.TrimRight(string(p), "\r\n")
	if msg == "" {
		return len(p), nil
	}
	slog.InfoContext(context.Background(), msg)
	return len(p), nil
}

func initLogging() {
	level := slog.LevelInfo
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL"))) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	jsonOut := strings.EqualFold(strings.TrimSpace(os.Getenv("LOG_FORMAT")), "json") ||
		strings.EqualFold(strings.TrimSpace(os.Getenv("GIN_MODE")), "release")

	opts := &slog.HandlerOptions{Level: level}
	var h slog.Handler
	if jsonOut {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = slog.NewTextHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(h))
}

func bridgeStdLog() {
	log.SetFlags(0)
	log.SetOutput(&logLineBridge{})
}
