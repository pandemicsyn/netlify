package enrichment

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
)

var (
	//ErrLogEntryExists occurr's if a chunfile log entry is already present in log entry datastore
	ErrLogEntryExists = errors.New("churnfile log entry exists")
)

// LogEntry is the basic mechanism we use to track whether .json churn files have already been processed
type LogEntry interface {
	CreateOrFail(bucket, object string) error
	Finalize(bucket, object string) error
}

// DatastoreLog is a LogEntry backed by Google Cloud Datstore
type DatastoreLog struct {
	client *datastore.Client
}

// DatastoreLogEntry is the behind scene's schema we use in Google Cloud Datastore
type DatastoreLogEntry struct {
	Success   bool      `datastore:"success:`
	Started   time.Time `datastore:"started"`
	Completed time.Time `datastore:"completed"`
}

func NewDatastoreLog(client *datastore.Client) LogEntry {
	return &DatastoreLog{client: client}
}

func (l *DatastoreLog) genKey(bucket, object string) *datastore.Key {
	return datastore.NameKey(ks, fmt.Sprintf("%s/%s", bucket, object), nil)
}

// CreateOrFail creates a churn profile batch file log entry or fails if the entry
// already exists (i.e. someone else already started to process it)
func (l *DatastoreLog) CreateOrFail(bucket, object string) error {
	key := l.genKey(bucket, object)
	ctx := context.Background()
	_, err := l.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var empty DatastoreLogEntry
		err := tx.Get(key, &empty)
		if err == nil {
			return ErrLogEntryExists
		}
		if err == datastore.ErrNoSuchEntity {
			_, err = tx.Put(key, &DatastoreLogEntry{
				Success: false,
				Started: time.Now(),
			})
			return err
		}
		return err
	})
	return err
}

// Finalize a churn profile batch file log entry thats been successfully processed
func (l *DatastoreLog) Finalize(bucket, object string) error {
	key := l.genKey(bucket, object)
	ctx := context.Background()
	tx, err := l.client.NewTransaction(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get transaction to update log entry")
	}
	var entry DatastoreLogEntry
	if err := tx.Get(key, &entry); err != nil {
		return errors.Wrap(err, "error getting log entry for update")
	}
	entry.Success = true
	entry.Completed = time.Now()
	if _, err := tx.Put(key, &entry); err != nil {
		return errors.Wrap(err, "error putting updated log entry")
	}
	if _, err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error on log entry update commit")
	}
	return nil
}
