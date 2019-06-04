package enrichment

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/pandemicsyn/netlify/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const FileCreated = "created"

// handleFileEvent does the bulk of the work when a new file event is received from our ETL process.
// right now it doesn't look utils.FileEvent.Status at all. In the future we have more than just 'created'
// it could switch on status and handle the different FileEvent status accordingly. i.e. reprocess, delete, etc
func (w *Worker) handleFileEvent(msg *pubsub.Message) error {
	var e utils.FileEvent
	if err := json.Unmarshal(msg.Data, &e); err != nil {
		return errors.Wrap(err, "Could not decode FileEvent data")
	}
	log.Debugf("Received notification: %v", e)

	// We create a log entry in google datastore when we start to process the file
	// if theres already an entry for the file we skip it as another process serviced it
	err := w.logEntry.CreateOrFail(e.Bucket, e.Object)
	if err != nil {
		if err == ErrLogEntryExists {
			log.Info("File already seen")
			return nil
		}
		return errors.Wrap(err, "Error creating log entry")
	}

	// next we snag an io.Reader of the file in storage
	pio, err := w.profileStore.Reader(e.Bucket, e.Object)
	if err != nil {
		return errors.Wrap(err, "Error getting object store stream")
	}

	// EnrichAndStore reads the json file out of storage, converting them to
	// EnrichedProfiles with a mock ChurnScore, and stores the profiles in our primary datastore
	err = EnrichAndStore(pio, w.eprofileStore)
	if err != nil {
		return errors.Wrap(err, "Error enriching or storing profiles")
	}

	// lastly we finalize the log entry indicating it was completed successfully.
	err = w.logEntry.Finalize(e.Bucket, e.Object)
	if err != nil {
		return errors.Wrap(err, "Error finalizing log entry")
	}
	return nil
}

// Receive just pull's file events off of the churn profile notification topic
func (w *Worker) Receive() error {
	ctx := context.Background()
	err := w.sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		err := w.handleFileEvent(msg)
		if err != nil {
			log.Warnf("Error handling file event: %v", err)
			//TODO: emit error metrics, alerting, etc
		} else {
			//TODO: emit processed metrics
			msg.Ack()
		}
	})
	if err != nil {
		return err
	}
	return nil
}
