package enrichment

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"github.com/pandemicsyn/netlify/utils"
	log "github.com/sirupsen/logrus"
)

func (w *Worker) handleFileEvent(msg *pubsub.Message) error {
	var e utils.FileEvent
	if err := json.Unmarshal(msg.Data, &e); err != nil {
		log.Printf("could not decode message data: %#v", msg)
		msg.Ack()
		return nil
	}
	log.Infof("Received notification: %v", e)
	key := datastore.NameKey(ks, fmt.Sprintf("%s/%s", e.Bucket, e.Object), nil)
	entry := NewDatastoreEntry()

	// We create a log entry in google datastore when we start to process the file
	// if theres already an entry for the file we skip it as another process serviced it
	err := entry.CreateOrFail(key, w.ds)
	if err != nil {
		if err == ErrLogEntryExists {
			log.Info("File already seen")
			msg.Ack()
			return nil
		}
		log.Warnf("Error creating log entry: %v", err)
		return nil
	}

	// next we snag an io.Reader of the file in storage
	pio, err := w.profilesReader(e.Bucket, e.Object)
	if err != nil {
		log.Warnf("Error getting object storage stream: %v", err)
	}

	// EnrichAndStore reads the json file out of storage, converting them to
	// EnrichedProfiles with a mock ChurnScore, and stores the profiles in our primary datastore
	err = EnrichAndStore(pio)
	if err != nil {
		log.Warnf("Error enriching or storing profiles: %v", err)
	}

	// lastly we finalize the log entry indicating it was completed successfully.
	err = entry.Finalize(key, w.ds)
	if err != nil {
		log.Warnf("Failed to finalize log entry: %v", err)
	}
	log.Println("done")
	msg.Ack()
	return nil
}

// Receive pull's file events off of churn profile notification topic
// when a new event is received it retrieves the object, enriches the contained
// profiles, and then stores them in our primary db
func (w *Worker) Receive() error {
	ctx := context.Background()
	err := w.sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		w.handleFileEvent(msg)
	})
	if err != nil {
		return err
	}
	return nil
}
