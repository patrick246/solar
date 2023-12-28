package listener

import (
	"context"
	"fmt"
	"log/slog"
)

type SLogLogger struct {
	logger *slog.Logger
	level  slog.Level
}

func NewLogger(logger *slog.Logger, level slog.Level) *SLogLogger {
	return &SLogLogger{
		logger: logger,
		level:  level,
	}
}

func (s *SLogLogger) Println(v ...interface{}) {
	s.logger.Log(context.Background(), s.level, fmt.Sprint(v...))
}

func (s *SLogLogger) Printf(format string, v ...interface{}) {
	s.logger.Log(context.Background(), s.level, fmt.Sprintf(format, v...))
}
