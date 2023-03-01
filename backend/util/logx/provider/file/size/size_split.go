package size

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type splitFileLogger struct {
	BaseDir    string
	Prefix     string
	DateFormat string
	UTCTime    bool
	file       *os.File
	mu         sync.Mutex
}

func (s *splitFileLogger) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	logFile := s.generateLogFile()

	if s.file == nil {
		if err = s.openExistingOrNew(logFile); err != nil {
			return 0, err
		}
	}
	if s.file.Name() != logFile {
		if err := s.rotate(logFile); err != nil {
			return 0, err
		}
	}
	return s.file.Write(p)
}

func (s *splitFileLogger) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.close()
}

func (s *splitFileLogger) close() error {
	if s.file == nil {
		return nil
	}
	err := s.file.Close()
	s.file = nil
	return err
}

func (s *splitFileLogger) generateLogFile() string {
	t := time.Now()
	if s.UTCTime {
		t = t.UTC()
	}
	timestamp := t.Format(s.DateFormat)
	return filepath.Join(s.BaseDir, fmt.Sprintf("%s.%s", s.Prefix, timestamp))
}

func (s *splitFileLogger) openExistingOrNew(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return s.openNew(filePath)
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return s.openNew(filePath)
	}
	s.file = file
	return nil
}

func (s *splitFileLogger) openNew(filePath string) error {
	err := os.MkdirAll(s.BaseDir, 0744)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	mode := os.FileMode(0644)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	s.file = f
	return nil
}

func (s *splitFileLogger) rotate(logFile string) error {
	if err := s.close(); err != nil {
		return err
	}
	if err := s.openNew(logFile); err != nil {
		return err
	}
	return nil
}
