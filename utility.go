package main

import (
	// standard library
	"io"
	"log"
	"os"

	// external packages
	"bitbucket.org/kardianos/osext"
)

// changeDirectoryToExecutable switches the current working directory to the executable.
// This is necessary for relative paths.
func changeDirectoryToExecutable() error {
	path, err := osext.ExecutableFolder()
	if err != nil {
		return err
	}
	return os.Chdir(path)
}

// initializeLogger sets the logging to a file and to stdout.
func initializeLogger() (err error) {
	if err = os.MkdirAll("../log", os.ModeDir|os.ModePerm); err != nil {
		log.Fatalf("Error creating log directory: %v", err)
		return
	}
	logfile, err := os.OpenFile("../log/sparkdo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return
	}
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	return
}
