package enrichment

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

// ProfileStore is our the io.Reader provider for our stored .json churn files
type ProfileStore interface {
	Reader(bucket, object string) (io.Reader, error)
}

// GCSProfileStore is Google Cloud Storage implementation of our ProfileStore
type GCSProfileStore struct {
	c *storage.Client
}

func NewGCSProfileStore(c *storage.Client) ProfileStore {
	return &GCSProfileStore{c}
}

// Reader provides GCS backed io.Reader
func (s *GCSProfileStore) Reader(bucket, object string) (io.Reader, error) {
	ctx := context.Background()
	return s.c.Bucket(bucket).Object(object).NewReader(ctx)
}
