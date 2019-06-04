package profiles

import (
	"errors"
	"math/rand"
	"strings"
	"testing"
)

func TestCalculateChurnRisk(t *testing.T) {
	e := EnrichedProfile{}
	_ = CalculateChurnRisk(e)
}

var ErrTestErr = errors.New("Test Error")
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type TestEProfileStore struct {
	Error error
}

func (l *TestEProfileStore) BulkSave(profiles *[]EnrichedProfile) error {
	return l.Error
}

func TestEnrichAndStore(t *testing.T) {
	testJSON := `{"CustomerID":"9237-HQITU","Partner":"No","Dependents":"No","Tenure":2,"PhoneService":"Yes","MultipleLines":"No","InternetService":"Fiber optic","OnlineSecurity":"No","OnlineBackup":"No","DeviceProtection":"No","TechSupport":"No","StreamingTV":"No","StreamingMovies":"No","Contract":"Month-to-month","PaperlessBilling":"Yes","PaymentMethod":"Electronic check","MonthlyCharges":70.7,"TotalCharges":151.65}`
	r := strings.NewReader(testJSON)
	estore := &TestEProfileStore{nil}
	err := EnrichAndStore(r, estore)
	if err != nil {
		t.Fatal(err)
	}

	estore = &TestEProfileStore{ErrTestErr}
	err = EnrichAndStore(r, estore)
	if err != nil {
		t.Fatal(err)
	}

	badJSON := `"CustomerID":"9237-HQITU","Partner":"No","Dependents":"No","Tenure":2,"PhoneService":"Yes","MultipleLines":"No","InternetService":"Fiber optic","OnlineSecurity":"No","OnlineBackup":"No","DeviceProtection":"No","TechSupport":"No","StreamingTV":"No","StreamingMovies":"No","Contract":"Month-to-month","PaperlessBilling":"Yes","PaymentMethod":"Electronic check","MonthlyCharges":70.7,"TotalCharges":151.65}`
	r = strings.NewReader(badJSON)
	estore = &TestEProfileStore{nil}
	err = EnrichAndStore(r, estore)
	if err == nil {
		t.Fatal("Bad json should have generated error")
	}
}
