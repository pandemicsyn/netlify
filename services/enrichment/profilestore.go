package enrichment

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

type ProfileStore interface {
	Reader(bucket, object string) (io.Reader, error)
}

type GCSProfileStore struct {
	c *storage.Client
}

func NewGCSProfileStore(c *storage.Client) ProfileStore {
	return &GCSProfileStore{c}
}

func (s *GCSProfileStore) Reader(bucket, object string) (io.Reader, error) {
	ctx := context.Background()
	return s.c.Bucket(bucket).Object(object).NewReader(ctx)
}
