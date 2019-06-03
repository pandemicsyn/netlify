package enrichment

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
)

func TestCreateOrFail(t *testing.T) {
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project == "" {
		project = "netlify-242319"
	}
	testEntry := NewDatastoreEntry()
	ctx := context.Background()
	tmpKey := time.Now().String()
	key := datastore.NameKey(ks, tmpKey, nil)
	client, err := datastore.NewClient(ctx, project)
	err = testEntry.CreateOrFail(key, client)
	if err != nil {
		t.Fatal(err)
	}

	//key should exist
	err = testEntry.CreateOrFail(key, client)
	if err != nil {
		if err != ErrLogEntryExists {
			t.Fatal(err)
		}
	}

	_ = client.Delete(ctx, key)
}

func TestFinalize(t *testing.T) {
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project == "" {
		project = "netlify-242319"
	}
	testEntry := NewDatastoreEntry()
	ctx := context.Background()
	tmpKey := time.Now().String()
	key := datastore.NameKey(ks, tmpKey, nil)
	client, err := datastore.NewClient(ctx, project)
	err = testEntry.CreateOrFail(key, client)
	if err != nil {
		t.Fatal(err)
	}

	//key should exist
	err = testEntry.Finalize(key, client)
	if err != nil {
		if err != ErrLogEntryExists {
			t.Fatal(err)
		}
	}

	_ = client.Delete(ctx, key)
}
