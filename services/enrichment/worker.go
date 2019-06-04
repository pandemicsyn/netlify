package enrichment

import (
	"context"
	"database/sql"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/pandemicsyn/netlify/pkg/profiles"
	"github.com/pandemicsyn/netlify/pkg/utils"
	"github.com/pkg/errors"
)

const sub = "enrichment-worker-test"
const topic = "churn-enrichment-test"
const ks = "churnfile"

// Worker is our enrichment worker
type Worker struct {
	project       string
	client        *pubsub.Client
	topic         *pubsub.Topic
	sub           *pubsub.Subscription
	ds            *datastore.Client
	logEntry      LogEntry
	os            *storage.Client
	profileStore  profiles.ProfileStore
	db            *sql.DB
	eprofileStore profiles.EProfileStore
}

// New enrichment worker
func New(project string, db *sql.DB) (*Worker, error) {
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
	w.logEntry = NewDatastoreLog(w.ds)

	w.os, err = storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating storage client failed")
	}
	w.profileStore = profiles.NewGCSProfileStore(w.os)

	w.db = db
	w.eprofileStore = profiles.NewPGEProfileStore(w.db)

	return w, nil
}
