package profiles

import (
	"errors"
	"log"
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
	Error    error
	Profiles *[]EnrichedProfile
}

func (l *TestEProfileStore) BulkSave(profiles *[]EnrichedProfile) error {
	l.Profiles = profiles
	return l.Error
}

func TestEnrichAndStore(t *testing.T) {

	testJSON := `{"CustomerID":"9237-HQITU","Partner":"No","Dependents":"No","Tenure":2,"PhoneService":"Yes","MultipleLines":"No","InternetService":"Fiber optic","OnlineSecurity":"No","OnlineBackup":"No","DeviceProtection":"No","TechSupport":"No","StreamingTV":"No","StreamingMovies":"No","Contract":"Month-to-month","PaperlessBilling":"Yes","PaymentMethod":"Electronic check","MonthlyCharges":70.7,"TotalCharges":151.65}`
	r := strings.NewReader(testJSON)
	estore := TestEProfileStore{nil, nil}
	err := EnrichAndStore(r, &estore)
	if err != nil {
		t.Fatal(err)
	}
	if len(*estore.Profiles) != 1 {
		for k, v := range *estore.Profiles {
			log.Println(k)
			log.Println(v)
		}
		t.Fatal("expected only 1 Profile in BulkSave call but got", len(*estore.Profiles))
	}

	estore = TestEProfileStore{ErrTestErr, nil}
	err = EnrichAndStore(r, &estore)
	if err != nil {
		t.Fatal(err)
	}

	badJSON := `"CustomerID":"9237-HQITU","Partner":"No","Dependents":"No","Tenure":2,"PhoneService":"Yes","MultipleLines":"No","InternetService":"Fiber optic","OnlineSecurity":"No","OnlineBackup":"No","DeviceProtection":"No","TechSupport":"No","StreamingTV":"No","StreamingMovies":"No","Contract":"Month-to-month","PaperlessBilling":"Yes","PaymentMethod":"Electronic check","MonthlyCharges":70.7,"TotalCharges":151.65}`
	r = strings.NewReader(badJSON)
	estore = TestEProfileStore{nil, nil}
	err = EnrichAndStore(r, &estore)
	if err == nil {
		t.Fatal("Bad json should have generated error")
	}
}
