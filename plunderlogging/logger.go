package plunderlogging

import "fmt"

// Logger - is a stuct that manages the verious types of logger available
type Logger struct {
	json JSONLogger
	file FileLogger
}

// EnableJSONLogging - will enable logging through JSON
func (l *Logger) EnableJSONLogging(e bool) {
	l.json.enabled = e
	l.json.initJSONLogger()
}

// EnableFileLogging - will enable logging to a file
func (l *Logger) EnableFileLogging(e bool) {
	l.file.enabled = e
}

// InitLogFile - will initialise file based logging
func (l *Logger) InitLogFile(path string) error {
	if l.file.enabled != true {
		return l.file.initFileLogger(path)
	}
	// Dont re-initialise the file
	return nil

}

// InitJSON - will start/initialise the JSON logging functionality
func (l *Logger) InitJSON() {
	// Dont re-initialise the json

	if l.json.enabled != true {
		l.json.initJSONLogger()
	}

}

// target - the entity we're affecting
// entry - the results of the operation on the target

// WriteLogEntry will capture what is transpiring and where
func (l *Logger) WriteLogEntry(target, task, entry, err string) {
	if l.file.enabled {
		l.file.writeEntry(target, entry)
	}
	if l.json.enabled {
		l.json.writeEntry(target, task, entry, err)
	}

	// A logging system shouldnt break anything so any errors are just outputed to STDOUT

}

// SetLoggingState - currently a NOOP (TODO)
func (l *Logger) SetLoggingState(target, state string) {
	if l.file.enabled {
		l.file.setLoggingState(target, state)
	}
	if l.json.enabled {
		l.json.setLoggingState(target, state)
	}

	// A logging system shouldnt break anything so any errors are just outputed to STDOUT

}

// GetJSONLogs - returns a pointer to the current JSON Logs
func (l *Logger) GetJSONLogs(target string) (*JSONLog, error) {
	if l.json.logger == nil {
		return nil, fmt.Errorf("JSON Logging hasn't been enabled")
	}
	// Check if the logger exists
	existingLog, ok := l.json.logger[target]
	if ok {
		return existingLog, nil
	}
	return nil, fmt.Errorf("No Logs for Target [%s] exist", target)
}

// DeleteLogs - will remove logs for a particular target
func (l *Logger) DeleteLogs(target string) error {
	if l.json.logger == nil {
		return nil
	}
	return l.json.deleteLog(target)

}
