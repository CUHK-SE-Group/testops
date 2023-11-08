package log_agent

import (
	slogmulti "github.com/samber/slog-multi"
	"log/slog"
	"os"
)

func NewBaseLogger(filename string) (*slog.Logger, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger := slog.New(
		slogmulti.Fanout(
			slog.NewTextHandler(file, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo,
			}),
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			}),
		))
	return logger, err
}

func NewMQLogger(filename string, brokers []string, topic string) (*slog.Logger, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger := slog.New(
		slogmulti.Fanout(
			slog.NewTextHandler(file, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo,
			}),
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			}),
			slog.NewJSONHandler(NewProducer[any](brokers, topic), &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo,
			}),
		))
	return logger, err
}
