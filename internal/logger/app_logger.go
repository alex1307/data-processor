package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Config represents the logger configuration
type CustomTextFormatter struct {
	SessionID string
}

func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Get the file and line number where the log was called
	_, filename, line, _ := runtime.Caller(7)

	// Get the script name from the full file path
	scriptName := filepath.Base(filename)

	// Format the log message
	message := fmt.Sprintf("[%s] [%s] [%s] [%s:%d] %s\n",
		entry.Time.Format("2006-01-02 15:04:05"), // Date-time
		entry.Level.String(),                     // Log level
		f.SessionID,                              // Unique session ID
		scriptName,                               // Script name
		line,                                     // Line number
		entry.Message,                            // Log message
	)

	return []byte(message), nil
}

// Generate a unique session ID
func GenerateSessionID() string {
	randomUUID := uuid.New()
	return strings.Replace(randomUUID.String(), "-", "", -1)
}

// Configure sets up the global logger based on the provided configuration
