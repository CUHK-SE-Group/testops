package log_agent

import (
	"fmt"
	slogmulti "github.com/samber/slog-multi"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	logger3 := slog.New(
		slogmulti.Fanout(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}), // pass to first handler: logstash over tcp
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}), // then to second handler: stderr
		),
	)

	logger3.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now()),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("A message")

}

func TestNewLogger2(t *testing.T) {
	logger := slog.New(
		slogmulti.Failover()(
			slog.NewJSONHandler(os.Stdout, nil), // send to this instance first
			slog.NewJSONHandler(os.Stderr, nil), // then this instance in case of failure
		),
	)

	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now()),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("A message")
}

func TestNewBaseLogger(t *testing.T) {
	logger, _ := NewBaseLogger("/tmp/test.log")
	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now()),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("A message")
}

func TestNewMQLogger(t *testing.T) {
	t.Skip()
	logger, _ := NewMQLogger("/tmp/test2.log", []string{"100.96.5.143:19092"}, "chat-room")
	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now()),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("A message")
	time.Sleep(4 * time.Second)
}
func TestNewLogger3(t *testing.T) {
	t.Skip()
	logger := slog.New(slog.NewJSONHandler(NewProducer[testMessage]([]string{"100.96.5.143:19092"}, "chat-room"), nil))
	logger = logger.With("release", "v1.0.0")

	for i := 0; i < 10; i++ {
		logger.
			With(
				slog.Group("user",
					slog.String("id", "user-123"),
					slog.Time("created_at", time.Now()),
				),
			).
			With("error", fmt.Errorf("an error")).
			Error("a message")
	}
	time.Sleep(5 * time.Second) // wait for requests to send
}
