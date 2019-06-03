package enrichment

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/pandemicsyn/netlify/utils"
	log "github.com/sirupsen/logrus"
)

func TestHandleFileEvent(t *testing.T) {
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project == "" {
		project = "netlify-242319"
	}
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, project)
	if err != nil {
		t.Fatal(err)
	}

	w := Worker{
		ds: client,
	}

	fe := utils.FileEvent{}
	data, _ := json.Marshal(fe)
	log.Println(data, w)
}
