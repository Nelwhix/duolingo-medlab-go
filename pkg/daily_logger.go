package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type DailyLogger struct {
	LogDir string
}

func (d *DailyLogger) Write(p []byte) (n int, err error) {
	fileName := time.Now().Format("2006-01-02") + ".log"
	filePath := filepath.Join(d.LogDir, fileName)

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	return f.Write(p)
}
