package enrichment

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type LogEntry interface {
	CreateOrFail(*datastore.Key, *datastore.Client) error
	Finalize(*datastore.Key, *datastore.Client) error
}

type DatastoreLogEntry struct {
	Success   bool
	Started   time.Time
	Completed time.Time
}

func NewDatastoreEntry() LogEntry {
	return &DatastoreLogEntry{}
}

// CreateOrFail creates a churn profile batch file log entry or fails if the entry
// already exists (i.e. someone else already started to process it)
func (l *DatastoreLogEntry) CreateOrFail(key *datastore.Key, client *datastore.Client) error {
	log.Println(key.String())
	ctx := context.Background()
	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
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
func (l *DatastoreLogEntry) Finalize(key *datastore.Key, client *datastore.Client) error {
	ctx := context.Background()
	tx, err := client.NewTransaction(ctx)
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
