// pkg/logger/slog_logger.go
package logger

import (
	"log/slog"
	"os"
)

// slogLogger คือนักข่าวสายโปรดักชัน ที่รายงานทุกอย่างเป็น JSON
type slogLogger struct {
	logger *slog.Logger
}

// NewSlogLogger คือโรงงานสร้างนักข่าวสายโปรดักชัน
// มันจะสร้าง Logger ที่พิมพ์ JSON ออกไปที่ Standard Output
func NewSlogLogger() Logger {
	return &slogLogger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

// --- Implementation of Logger interface ---

func (l *slogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *slogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *slogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *slogLogger) Error(msg string, err error, args ...any) {
	// slog จะฉลาดพอที่จะรู้ว่าถ้ามี key ชื่อ "err" มันจะแสดงผลให้สวยงามเป็นพิเศษ
	allArgs := append(args, "err", err)
	l.logger.Error(msg, allArgs...)
}

func (l *slogLogger) Success(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *slogLogger) Print(msg string) {
	l.logger.Debug(msg)
}

func (l *slogLogger) Dump(data interface{}) {
	l.logger.Debug("dump_data", data)
}

func (l *slogLogger) Dumpf(level string, msg string, data interface{}) {
	switch level {
	case LevelDebug:
		l.logger.Debug(msg, "dumpf_data", data)
	case LevelInfo:
		l.logger.Info(msg, "dumpf_data", data)
	case LevelWarn:
		l.logger.Warn(msg, "dumpf_data", data)
	case LevelError:
		l.logger.Error(msg, "dumpf_data", data)
	case LevelSuccess:
		l.logger.Info(msg, "dumpf_data", data)
	}
}
