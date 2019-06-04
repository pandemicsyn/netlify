package etlfunctions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/pandemicsyn/netlify/pkg/events"
	"github.com/pandemicsyn/netlify/pkg/utils"
	"github.com/pandemicsyn/netlify/transform"
)

// GCSEvent is the stock struct for GCS events
type GCSEvent struct {
	Bucket         string `json:"bucket"`
	Name           string `json:"name"`
	Metageneration string `json:"metageneration"`
	ResourceState  string `json:"resourceState"`
}

var (
	successBucketName = "netlify-churncsv-success"
	client            *pubsub.Client
	topic             *pubsub.Topic
	project           = os.Getenv("GOOGLE_CLOUD_PROJECT")
)

func init() {
	var err error
	ctx := context.Background()
	client, err = pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Could not create pubsub client: %v", err)
	}
	topic, err = utils.CreateTopicIfNotExists(client, utils.DefaultTopic)
	if err != nil {
		log.Fatalf("Could not create or acquire pubsub topic: %v", err)
	}
}

func convert(e GCSEvent, successBucketName string, successObjectName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return err
	}
	rc, err := client.Bucket(e.Bucket).Object(e.Name).NewReader(ctx)
	if err != nil {
		log.Printf("Failed to acquire Reader on csv: %v", err)
		return err
	}
	defer rc.Close()
	// TODO: write compressed file ?
	wc := client.Bucket(successBucketName).Object(successObjectName).NewWriter(ctx)
	wc.ContentType = "application/json"
	err = transform.Transform(rc, wc)
	if err != nil {
		return err
	}
	// gcloud storage streaming writes most often error on close when object is finalized
	if err := wc.Close(); err != nil {
		//TODO: attempt to upload to failure bucket/preform recovery/notification etc
		log.Printf("Error on object write close - likely failed: %v", err)
		return err
	}
	return nil
}

// ChurnTransform runs the etl process to clean/transfrom churn data for
// long term storage when the raw csv arrives at the interim bucket.
func ChurnTransform(ctx context.Context, e GCSEvent) error {
	if e.ResourceState == "not_exists" {
		log.Printf("File %v removed unexpectedly", e.Name)
		return nil
	}
	if e.Metageneration == "1" {
		log.Printf("New batch file %v/%v created", e.Bucket, e.Name)
		now := time.Now()
		successObjectName := fmt.Sprintf("%s/%s-churn-profiles.json", now.Format("2006-01-02"), now.Format("0304"))
		err := convert(e, successBucketName, successObjectName)
		payload, err := json.Marshal(events.FileEvent{successBucketName, successObjectName, events.StatusCreated, 1})
		if err != nil {
			log.Printf("Failed to encode FileEvent json: %v", err)
			// TODO: track differently than a regular failure - since new object already exists
			return nil
		}
		r := topic.Publish(ctx, &pubsub.Message{Data: payload})
		id, err := r.Get(ctx)
		if err != nil {
			log.Printf("Failed to publish FileEvent: %v", err)
			// TODO: track differently than a regular failure - since new object already exists
			return nil
		}
		log.Printf("Sent FileEvent: %s", id)
		return nil
	}
	return nil
}
