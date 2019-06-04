package enrichment

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
)

func TestCreateOrFail(t *testing.T) {
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project == "" {
		project = "netlify-242319"
	}
	ctx := context.Background()
	tmpKey := time.Now().String()
	client, err := datastore.NewClient(ctx, project)
	if err != nil {
		t.Fatal(err)
	}
	logEntry := NewDatastoreLog(client)

	err = logEntry.CreateOrFail(tmpKey, tmpKey)
	if err != nil {
		t.Fatal(err)
	}

	//key should exist
	err = logEntry.CreateOrFail(tmpKey, tmpKey)
	if err != nil {
		if err != ErrLogEntryExists {
			t.Fatal(err)
		}
	}
	key := datastore.NameKey(ks, fmt.Sprintf("%s/%s", tmpKey, tmpKey), nil)
	_ = client.Delete(ctx, key)
}

func TestFinalize(t *testing.T) {
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project == "" {
		project = "netlify-242319"
	}
	ctx := context.Background()
	tmpKey := time.Now().String()
	client, err := datastore.NewClient(ctx, project)
	if err != nil {
		t.Fatal(err)
	}
	logEntry := NewDatastoreLog(client)

	//key should not yet exist, so this should fail
	err = logEntry.Finalize(tmpKey, tmpKey)
	if err != nil {
		if errors.Cause(err) != datastore.ErrNoSuchEntity {
			t.Fatal("Unexpected error finalizing entry, expect ErrNoSuchEntity, got:", err)
		}
	}
	if err == nil {
		t.Fatal("Expected Finalize LogEntry key to not yet exist, but it did!")
	}

	// test successful creation and finalization
	err = logEntry.CreateOrFail(tmpKey, tmpKey)
	if err != nil {
		t.Fatal(err)
	}

	//key should exist - so we should be able to finalize the entry
	err = logEntry.Finalize(tmpKey, tmpKey)
	if err != nil {
		if err != ErrLogEntryExists {
			t.Fatal(err)
		}
	}

	//cleanup
	key := datastore.NameKey(ks, fmt.Sprintf("%s/%s", tmpKey, tmpKey), nil)
	_ = client.Delete(ctx, key)
}
