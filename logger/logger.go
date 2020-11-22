// Copyright 2020 Scott Underwood.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package logger provides logging capability for an application or service.
// TODO: make this its own go module outside of cleanpg
package logger

import (
	"log"
	"os"
)

// MessageType holds the log level of the message.
// Possible values:
// INFO | NOTICE | WARNING | ERROR | FATAL
type MessageType int

const (
	// INFO indicates generally useful information
	INFO MessageType = 0
	// NOTICE indicates program state changes that are not abnormal
	NOTICE MessageType = iota
	// WARNING indicates application oddities that are recoverable
	WARNING
	// ERROR indicates condition fatal to the operation
	// but not to the application
	ERROR
	// FATAL indicates condition is fatal to the application or service
	// and will force a shutdown
	FATAL
)

var (
	logger       []*log.Logger // slice of loggers for each level
	stderrLogger *log.Logger   // stderr logger
)

var (
	logFileName string   = "log.txt" // holds name of log file
	logToStderr bool     = false     // flag to indicate whether loggint to stderr
	logFileFD   *os.File             // log file descriptor
)

// createLogFile is called from the logWriter if the log file is not open
func createLogFile(logFileFD *os.File) (*os.File, error) {
	logFileFD, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(logFileFD)

	return logFileFD, nil
}

// LogToStderr determines whether log messages will print to stderr
// as well as the log file
func LogToStderr(flag bool) {
	logToStderr = flag
}

// SetLogFile sets the name of the log file.
// If not set, the default filename is "log.txt"
func SetLogFile(fileName string) {
	closeLogFile()
	logFileName = fileName
	initLoggers()
}

// Truncate is used to truncate the log file to zero length
func Truncate() error {
	// If file doesn't exist, no need to truncate
	_, err := os.Stat(logFileName)
	if os.IsNotExist(err) {
		return nil
	}

	// Truncate it
	err = os.Truncate(logFileName, 0)
	if err != nil {
		log.Fatalf("Could not truncate log file: [%s]", err)
		return err
	}

	return nil
}

// Write is a function which writes a variable length string message to the log file
func Write(messageType MessageType, format string, a ...interface{}) {

	if logToStderr {
		stderrLogger.SetPrefix(logger[messageType].Prefix())
		stderrLogger.Printf(format, a...)
	}

	logger[messageType].Printf(format, a...)

}

// initLoggers initializes a log file and loggers for each level
func initLoggers() {

	var err error
	logFileFD, err = createLogFile(logFileFD)
	if err != nil {
		log.Fatalf("Could not create log file: [%s]", err)
	}

	// Logger flags
	const lflags int = log.Ldate | log.Ltime | log.Lmsgprefix

	// Build slice of loggers for each level
	logger = append(logger, log.New(logFileFD, "INFO: ", lflags))
	logger = append(logger, log.New(logFileFD, "NOTICE: ", lflags))
	logger = append(logger, log.New(logFileFD, "WARNING: ", lflags))
	logger = append(logger, log.New(logFileFD, "ERROR: ", lflags))
	logger = append(logger, log.New(logFileFD, "FATAL: ", lflags))

	// Special logger to handle output to stderr
	stderrLogger = log.New(os.Stderr, "", lflags)
}

// closeLogFile closes the current fd and removes the logfile if zero length
func closeLogFile() {
	// Close it
	if err := logFileFD.Close(); err != nil {
		log.Fatalf("Could not close log file [%s]: [%s]", logFileName, err)
		return
	}

	// Get the length of the file via stat
	info, err := os.Stat(logFileName)
	if err != nil {
		log.Fatalf("Could not stat log file [%s]: [%s]", logFileName, err)
		return
	}

	// If it's zero length, remove it
	if info.Size() == 0 {
		err := os.Remove(logFileName)
		if err != nil {
			log.Fatalf("Could not remove [%s]: [%s]", logFileName, err)
			return
		}
	}
}

func init() {
	initLoggers()
}
