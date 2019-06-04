package transform

import (
	"io"
	"log"
)

// Transform a churn profile csv batch file on upload to a short term bucket
// into a cleaned json file stored in long term storage.
func Transform(src io.Reader, dst io.Writer) error {
	err := CsvToJSON(src, dst)
	if err != nil {
		//TODO: attempt to quarantine original file, perform notifications etc
		log.Printf("Error converting csv to json: %v", err)
		return err
	}
	return nil
}
