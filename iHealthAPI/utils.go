package ihealthapi

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

// Check if a file exists return true if it does or false if it doesn't
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return errors.Is(err, os.ErrNotExist)
}

// Write the file and return the file pointer
func writeToFile(filename string, rc io.Reader) (*os.File, bool) {
	// Get the contents
	d, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Errorf("writeToFile() -> Problem getting contents of rc body: %s", err)
		return nil, false
	}

	// Create the file
	f, err := os.Create(filename)
	if err != nil {
		log.Errorf("writeToFile() -> Problem creating temp file at %s: %s", TempLogFile, err)
		return nil, false
	}
	defer f.Close()

	//Write to the file
	byteswritten, err := f.Write(d)
	if err != nil {
		log.Errorf("writeToFile() -> Failed to write to ./Temp/%s: %s", filename, err)
		return nil, false
	}

	log.Infof("%d bytes written to ./Temp/%s", byteswritten, filename)
	return f, true
}
