package churnprofiles

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
)

// Transform a churn profile csv batch file on upload to a short term bucket
// into a cleaned json file stored in long term storage. Destination files use a path format of:
// longtermBucket/2006-01-02/0304-churn-profiles.json
func Transform(readBucket, successBucket, objectName string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return err
	}
	rc, err := client.Bucket(readBucket).Object(objectName).NewReader(ctx)
	if err != nil {
		log.Printf("Failed to acquire Reader on csv: %v", err)
		return err
	}
	defer rc.Close()

	now := time.Now()
	// TODO: write compressed file ?
	wc := client.Bucket(successBucket).Object(fmt.Sprintf("%s/%s-churn-profiles.json", now.Format("2006-01-02"), now.Format("0304"))).NewWriter(ctx)
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
