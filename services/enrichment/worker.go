package enrichment

import (
	"context"
	"io"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/pandemicsyn/netlify/utils"
	"github.com/pkg/errors"
)

const sub = "enrichment-worker-test"
const topic = "churn-enrichment-test"
const ks = "churnfile"

var (
	ErrLogEntryExists = errors.New("churn file log entry exists")
)

type Worker struct {
	project string
	client  *pubsub.Client
	topic   *pubsub.Topic
	sub     *pubsub.Subscription
	ds      *datastore.Client
	os      *storage.Client
}

// New enrichment worker
func New(project string) (*Worker, error) {
	w := &Worker{
		project: project,
	}
	var err error
	ctx := context.Background()
	w.client, err = pubsub.NewClient(ctx, w.project)
	if err != nil {
		return nil, errors.Wrap(err, "pubsub client create failed")
	}

	w.topic, err = utils.CreateTopicIfNotExists(w.client, utils.DefaultTopic)
	if err != nil {
		return nil, err
	}

	w.sub, err = utils.CreateSub(w.client, sub, w.topic)
	if err != nil {
		return nil, errors.Wrap(err, "creating subscription failed")
	}

	w.ds, err = datastore.NewClient(ctx, w.project)
	if err != nil {
		return nil, errors.Wrap(err, "creating datastore client failed")
	}

	w.os, err = storage.NewClient(ctx)
	return w, nil
}

func (w *Worker) profilesReader(bucket, object string) (io.Reader, error) {
	ctx := context.Background()
	return w.os.Bucket(bucket).Object(object).NewReader(ctx)
}
