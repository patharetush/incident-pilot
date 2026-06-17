package logging

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Filename      string
	MaxFileSizeMB int
	MaxBackups    int
	MaxAgeDays    int
}

var DefaultConfig = Config{
	Filename:      "/tmp/app.log",
	MaxFileSizeMB: 10,
	MaxBackups:    5,
	MaxAgeDays:    5,
}

var (
	mu          sync.Mutex
	initialized bool
	defaultLog  = zerolog.New(os.Stderr).With().Timestamp().Logger()
	active      atomic.Pointer[zerolog.Logger]
	rotator     *lumberjack.Logger
)

func init() {
	active.Store(&defaultLog)
}

func Init() error {
	return InitWithConfig(DefaultConfig)
}

func InitWithConfig(cfg Config) error {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return nil
	}

	if cfg.Filename == "" {
		return errors.New("log filename cannot be empty")
	}
	if err := os.MkdirAll(filepath.Dir(cfg.Filename), 0o755); err != nil {
		L().Error().Err(err).Str("path", cfg.Filename).Msg("failed to create log directory, falling back to stderr")
		zerolog.TimeFieldFormat = time.RFC3339Nano
		initialized = true
		return nil
	}

	writer := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxFileSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
	}
	logger := zerolog.New(writer).
		With().
		Timestamp().
		Logger()
	active.Store(&logger)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	rotator = writer
	initialized = true

	return nil
}

func L() *zerolog.Logger {
	if logger := active.Load(); logger != nil {
		return logger
	}
	return &defaultLog
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()
	if rotator != nil {
		err := rotator.Close()
		rotator = nil
		return err
	}
	return nil
}
