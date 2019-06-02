package churnprofiles

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
)

// Transform a churn profile csv batch file on upload to a short term bucket
// into a cleaned json file stored in long term storage. Destination files use a path format of:
// longtermBucket/2006-01-02/0304-churn-profiles.json
func Transform(readBucket, readObject, successBucket, successObject string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return err
	}
	rc, err := client.Bucket(readBucket).Object(readObject).NewReader(ctx)
	if err != nil {
		log.Printf("Failed to acquire Reader on csv: %v", err)
		return err
	}
	defer rc.Close()
	// TODO: write compressed file ?
	wc := client.Bucket(successBucket).Object(successObject).NewWriter(ctx)
	wc.ContentType = "application/json"

	err = CsvToJSON(rc, wc)
	if err != nil {
		//TODO: attempt to quarantine original file, perform notifications etc
		log.Printf("Error parsing csv: %v", err)
		return err
	}

	// gcloud storage streaming writes most likely error on close when object is finalized
	if err := wc.Close(); err != nil {
		//TODO: attempt to upload to failure bucket/preform recovery/notification etc
		log.Printf("Error on object write close - likely failed: %v", err)
		return err
	}
	log.Println("successfully wrote new object:", wc.Attrs())
	return nil
}
