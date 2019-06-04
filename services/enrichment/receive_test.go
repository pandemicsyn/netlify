package enrichment

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/pandemicsyn/netlify/pkg/events"
	"github.com/pandemicsyn/netlify/pkg/profiles"
	"github.com/pkg/errors"
)

type TestLogEntry struct {
	CreateOrFailResult error
	FinalizeResult     error
}

var ErrTestErr = errors.New("Test Error")

func (l *TestLogEntry) CreateOrFail(k, v string) error {
	return l.CreateOrFailResult
}

func (l *TestLogEntry) Finalize(k, v string) error {
	return l.FinalizeResult
}

type TestProfileStore struct {
	R     io.Reader
	Error error
}

func (l *TestProfileStore) Reader(k, v string) (io.Reader, error) {
	return l.R, l.Error
}

type TestEProfileStore struct {
	Error error
}

func (l *TestEProfileStore) BulkSave(profiles *[]profiles.EnrichedProfile) error {
	return l.Error
}

func TestHandleFileEvent(t *testing.T) {

	w := Worker{
		logEntry:      &TestLogEntry{ErrTestErr, nil},
		eprofileStore: &TestEProfileStore{nil},
	}

	fe := events.FileEvent{}
	data, _ := json.Marshal(fe)

	//test garbage json
	err := w.handleFileEvent(&pubsub.Message{Data: []byte("stuff")})
	if err == nil {
		t.Fatal("Expected unable to decode json error, but receive no error")
	}

	//test ok json but logEntry.CreateOrFail finds log entry already exists
	w.logEntry = &TestLogEntry{ErrLogEntryExists, nil}
	err = w.handleFileEvent(&pubsub.Message{Data: data})
	if err != nil {
		t.Fatal("Should have skipped this log as it should have already been seen")
	}

	//test ok json but logEntry.CreateOrFail throws other error
	w.logEntry = &TestLogEntry{ErrTestErr, nil}
	err = w.handleFileEvent(&pubsub.Message{Data: data})
	if err != nil {
		if !strings.HasPrefix(err.Error(), "Error creating log entry") {
			t.Fatal("Expected error creating log entry but got:", err)
		}
	} else {
		t.Fatal("Expected error creating log entry but got nil.")
	}

	//test error from object store stream
	w.logEntry = &TestLogEntry{nil, nil}
	w.profileStore = &TestProfileStore{nil, ErrTestErr}
	err = w.handleFileEvent(&pubsub.Message{Data: data})
	if err == nil {
		t.Fatal("Expected getting object store stream to fail")
	}

	//test ok enrichAndStore
	w.logEntry = &TestLogEntry{nil, nil}
	testJSON := `{"CustomerID":"9237-HQITU","Partner":"No","Dependents":"No","Tenure":2,"PhoneService":"Yes","MultipleLines":"No","InternetService":"Fiber optic","OnlineSecurity":"No","OnlineBackup":"No","DeviceProtection":"No","TechSupport":"No","StreamingTV":"No","StreamingMovies":"No","Contract":"Month-to-month","PaperlessBilling":"Yes","PaymentMethod":"Electronic check","MonthlyCharges":70.7,"TotalCharges":151.65}`
	r := strings.NewReader(testJSON)
	w.profileStore = &TestProfileStore{r, nil}
	err = w.handleFileEvent(&pubsub.Message{Data: data})
	if err != nil {
		t.Fatal("Expected no issues enriching profile but got:", err)
	}

	//test malformed profile in enrichAndStore
	w.logEntry = &TestLogEntry{nil, nil}
	testJSON = `"CustomerID":"9237-HQITU","Partner":"No","Dependents":"No","Tenure":2,"PhoneService":"Yes","MultipleLines":"No","InternetService":"Fiber optic","OnlineSecurity":"No","OnlineBackup":"No","DeviceProtection":"No","TechSupport":"No","StreamingTV":"No","StreamingMovies":"No","Contract":"Month-to-month","PaperlessBilling":"Yes","PaymentMethod":"Electronic check","MonthlyCharges":70.7,"TotalCharges":151.65}`
	r = strings.NewReader(testJSON)
	w.profileStore = &TestProfileStore{r, nil}
	err = w.handleFileEvent(&pubsub.Message{Data: data})
	if err == nil {
		t.Fatal("Expected err from EnrichAndStore due to malformed profile but got nil")
	}

	//test finalize log entry failed
	w.logEntry = &TestLogEntry{nil, ErrTestErr}
	testJSON = `{"CustomerID":"9237-HQITU","Partner":"No","Dependents":"No","Tenure":2,"PhoneService":"Yes","MultipleLines":"No","InternetService":"Fiber optic","OnlineSecurity":"No","OnlineBackup":"No","DeviceProtection":"No","TechSupport":"No","StreamingTV":"No","StreamingMovies":"No","Contract":"Month-to-month","PaperlessBilling":"Yes","PaymentMethod":"Electronic check","MonthlyCharges":70.7,"TotalCharges":151.65}`
	r = strings.NewReader(testJSON)
	w.profileStore = &TestProfileStore{r, nil}
	err = w.handleFileEvent(&pubsub.Message{Data: data})
	if err == nil {
		t.Fatal("Expected error finalizing log entry but got nil")
	}
}
