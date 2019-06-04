package enrichment

import (
	"math/rand"
	"strings"
	"testing"
)

func TestCalculateChurnRisk(t *testing.T) {
	e := EnrichedProfile{}
	_ = CalculateChurnRisk(e)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestEnrichAndStore(t *testing.T) {
	testJSON := `{"CustomerID":"9237-HQITU","Partner":"No","Dependents":"No","Tenure":2,"PhoneService":"Yes","MultipleLines":"No","InternetService":"Fiber optic","OnlineSecurity":"No","OnlineBackup":"No","DeviceProtection":"No","TechSupport":"No","StreamingTV":"No","StreamingMovies":"No","Contract":"Month-to-month","PaperlessBilling":"Yes","PaymentMethod":"Electronic check","MonthlyCharges":70.7,"TotalCharges":151.65}`
	r := strings.NewReader(testJSON)
	err := EnrichAndStore(r)
	if err != nil {
		t.Fatal(err)
	}

	badJSON := `"CustomerID":"9237-HQITU","Partner":"No","Dependents":"No","Tenure":2,"PhoneService":"Yes","MultipleLines":"No","InternetService":"Fiber optic","OnlineSecurity":"No","OnlineBackup":"No","DeviceProtection":"No","TechSupport":"No","StreamingTV":"No","StreamingMovies":"No","Contract":"Month-to-month","PaperlessBilling":"Yes","PaymentMethod":"Electronic check","MonthlyCharges":70.7,"TotalCharges":151.65}`
	r = strings.NewReader(badJSON)
	err = EnrichAndStore(r)
	if err == nil {
		t.Fatal("Bad json should have generated error")
	}
}
