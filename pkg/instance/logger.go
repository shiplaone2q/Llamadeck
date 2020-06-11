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

