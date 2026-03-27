package instance

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	timber "github.com/DeRuina/timberjack"
)

// LogRotationConfig contains log rotation settings for instances
type LogRotationConfig struct {
	Enabled  bool
	MaxSize  int
	Compress bool
}

type logger struct {
	name        string
	logDir      string
	logFile     *timber.Logger
	logFilePath string
	mu          sync.RWMutex
	cfg         *LogRotationConfig
}

func newLogger(name, logDir string, cfg *LogRotationConfig) *logger {
	return &logger{
		name:   name,
		logDir: logDir,
		cfg:    cfg,
	}
}

func (l *logger) create() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logDir == "" {
		return fmt.Errorf("logDir empty for instance %s", l.name)
	}

	if err := os.MkdirAll(l.logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logPath := fmt.Sprintf("%s/%s.log", l.logDir, l.name)
	l.logFilePath = logPath

	// Build the timber logger
	t := &timber.Logger{
		Filename:   logPath,
		MaxSize:    l.cfg.MaxSize,
		MaxBackups: 0, // No limit on backups
		// Compression: "gzip" if Compress is true, else "none"
		Compression: func() string {
			if l.cfg.Compress {
				return "gzip"
			}
			return "none"
		}(),
		FileMode:  0644,
		LocalTime: true,
	}

	// If rotation is disabled, set MaxSize to 0 so no rotation occurs
	if !l.cfg.Enabled {
		t.MaxSize = 0
	}

	l.logFile = t

