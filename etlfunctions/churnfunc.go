package etlfunctions

import (
	"context"
	"log"

	"github.com/pandemicsyn/netlify/churnprofiles"
)

// GCSEvent is the stock struct for GCS events
type GCSEvent struct {
	Bucket         string `json:"bucket"`
	Name           string `json:"name"`
	Metageneration string `json:"metageneration"`
	ResourceState  string `json:"resourceState"`
}

var successBucketName = "netlify-churncsv-success"

// ChurnTransform runs the etl process to clean/transfrom churn data for
// long term storage when the raw csv arrives at the interim bucket.
func ChurnTransform(ctx context.Context, e GCSEvent) error {
	if e.ResourceState == "not_exists" {
		log.Printf("File %v removed unexpectedly", e.Name)
		return nil
	}
	if e.Metageneration == "1" {
		log.Printf("New batch file %v/%v created", e.Bucket, e.Name)
		churnprofiles.Transform(e.Bucket, successBucketName, e.Name)
		return nil
	}
	// TODO: if processing of the file fails we could place a temporary hold
	// on the object until it can be reprocessed. This will make sure data isn't purged
	// from the bucket until we've successfully ingested it.
	return nil
}
