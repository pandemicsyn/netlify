package transform

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/pandemicsyn/netlify/pkg/profiles"
)

func fileSrc() (*bufio.Reader, *os.File) {
	f, err := os.Open("Telco-Customer-Churn.csv")
	if err != nil {
		panic(err)
	}
	bf := bufio.NewReader(f)
	return bf, f
}

func stdOutSrc() io.Writer {
	return os.Stdout
}

func TestSetProfileValue(t *testing.T) {

	wantStr := "good"

	//since some column names are subtly different from the struct
	//support maping column to struct field.

	testStringCols := map[string]string{
		"customerID":       "CustomerID",
		"Partner":          "Partner",
		"Dependents":       "Dependents",
		"PhoneService":     "PhoneService",
		"MultipleLines":    "MultipleLines",
		"InternetService":  "InternetService",
		"OnlineSecurity":   "OnlineSecurity",
		"OnlineBackup":     "OnlineBackup",
		"DeviceProtection": "DeviceProtection",
		"TechSupport":      "TechSupport",
		"StreamingTV":      "StreamingTV",
		"StreamingMovies":  "StreamingMovies",
		"Contract":         "Contract",
		"PaperlessBilling": "PaperlessBilling",
		"PaymentMethod":    "PaymentMethod",
	}

	for field, structField := range testStringCols {
		c := &profiles.ChurnProfile{}
		err := setProfileValue(c, field, wantStr)
		if err != nil {
			t.Fatal(err)
		}
		r := reflect.ValueOf(c)
		got := reflect.Indirect(r).FieldByName(structField).String()
		if !reflect.DeepEqual(wantStr, got) {
			t.Fatalf("expected: '%v' for field '%v', got: '%v'", wantStr, structField, got)
		}
	}

	intStrOk := "42"
	intOk := int64(42)
	intZero := int64(0)
	intStrEmpty := ""
	intStrSpace := " "
	intBad := "."

	testIntCols := map[string]string{
		"tenure": "Tenure",
	}

	// test int columns with valid int
	for field, structField := range testIntCols {
		c := &profiles.ChurnProfile{}
		err := setProfileValue(c, field, intStrOk)
		if err != nil {
			t.Fatal(err)
		}
		r := reflect.ValueOf(c)
		got := reflect.Indirect(r).FieldByName(structField).Int()
		if !reflect.DeepEqual(intOk, got) {
			t.Fatalf("expected: '%v' for field '%v', got: '%v'", intOk, structField, got)
		}
	}

	// test int columns with empty fields
	for field, structField := range testIntCols {
		c := &profiles.ChurnProfile{}
		err := setProfileValue(c, field, intStrEmpty)
		if err != nil {
			t.Fatal(err)
		}
		r := reflect.ValueOf(c)
		got := reflect.Indirect(r).FieldByName(structField).Int()
		if !reflect.DeepEqual(intZero, got) {
			t.Fatalf("expected: '%v' for field '%v', got: '%v'", 0, structField, got)
		}
	}

	// test int columns with space as field
	for field, structField := range testIntCols {
		c := &profiles.ChurnProfile{}
		err := setProfileValue(c, field, intStrSpace)
		if err != nil {
			t.Fatal(err)
		}
		r := reflect.ValueOf(c)
		got := reflect.Indirect(r).FieldByName(structField).Int()
		if !reflect.DeepEqual(intZero, got) {
			t.Fatalf("expected: '%v' for field '%v', got: '%v'", 0, structField, got)
		}
	}

	// test int with string content
	for field, structField := range testIntCols {
		c := &profiles.ChurnProfile{}
		err := setProfileValue(c, field, intBad)
		if err == nil {
			t.Fatalf("expected error for field '%v', sent string", structField)
		}
	}

	floatStrOk := "42.0"
	floatOk := float64(42.0)
	floatZeroOk := float64(0.0)
	floatStrEmpty := ""
	floatStrSpace := " "
	floatBad := "."

	testFloatCols := map[string]string{
		"TotalCharges":   "TotalCharges",
		"MonthlyCharges": "MonthlyCharges",
	}

	// test float columns with valid floats
	for field, structField := range testFloatCols {
		c := &profiles.ChurnProfile{}
		err := setProfileValue(c, field, floatStrOk)
		if err != nil {
			t.Fatal(err)
		}
		r := reflect.ValueOf(c)
		log.Println(r)
		got := reflect.Indirect(r).FieldByName(structField).Float()
		if !reflect.DeepEqual(floatOk, got) {
			t.Fatalf("expected: '%v' for field '%v', got: '%v'", floatOk, structField, got)
		}
	}

	// test float columns with empty fields
	for field, structField := range testFloatCols {
		c := &profiles.ChurnProfile{}
		err := setProfileValue(c, field, floatStrEmpty)
		if err != nil {
			t.Fatal(err)
		}
		r := reflect.ValueOf(c)
		log.Println(r)
		got := reflect.Indirect(r).FieldByName(structField).Float()
		if !reflect.DeepEqual(floatZeroOk, got) {
			t.Fatalf("expected: '%v' for field '%v', got: '%v'", floatZeroOk, structField, got)
		}
	}

	// test float columns with space as field
	for field, structField := range testFloatCols {
		c := &profiles.ChurnProfile{}
		err := setProfileValue(c, field, floatStrSpace)
		if err != nil {
			t.Fatal(err)
		}
		r := reflect.ValueOf(c)
		log.Println(r)
		got := reflect.Indirect(r).FieldByName(structField).Float()
		if !reflect.DeepEqual(floatZeroOk, got) {
			t.Fatalf("expected: '%v' for field '%v', got: '%v'", floatZeroOk, structField, got)
		}
	}

	// test float with string content
	for field, structField := range testFloatCols {
		c := &profiles.ChurnProfile{}
		err := setProfileValue(c, field, floatBad)
		if err == nil {
			t.Fatalf("expected error for field '%v', sent string", structField)
		}
	}

	// test non existent field
	c := &profiles.ChurnProfile{}
	err := setProfileValue(c, "notpresent", "novalue")
	if err != nil {
		t.Fatal(err)
	}
}

var goodCsv = `customerID,gender,SeniorCitizen,Partner,Dependents,tenure,PhoneService,MultipleLines,InternetService,OnlineSecurity,OnlineBackup,DeviceProtection,TechSupport,StreamingTV,StreamingMovies,Contract,PaperlessBilling,PaymentMethod,MonthlyCharges,TotalCharges,Churn
7590-VHVEG,Female,0,Yes,No,1,No,No phone service,DSL,No,Yes,No,No,No,No,Month-to-month,Yes,Electronic check,29.85,29.85,No
5575-GNVDE,Male,0,No,No,34,Yes,No,DSL,Yes,No,Yes,No,No,No,One year,No,Mailed check,56.95,1889.5,No
3668-QPYBK,Male,0,No,No,2,Yes,No,DSL,Yes,Yes,No,No,No,No,Month-to-month,Yes,Mailed check,53.85,108.15,Yes
7795-CFOCW,Male,0,No,No,45,No,No phone service,DSL,Yes,No,Yes,Yes,No,No,One year,No,Bank transfer (automatic),42.3,1840.75,No
9237-HQITU,Female,0,No,No,2,Yes,No,Fiber optic,No,No,No,No,No,No,Month-to-month,Yes,Electronic check,70.7,151.65,Yes
`

var badCsv = `customerID,gender,SeniorCitizen,Partner,Dependents,tenure,PhoneService,MultipleLines,InternetService,OnlineSecurity,OnlineBackup,DeviceProtection,TechSupport,StreamingTV,StreamingMovies,Contract,PaperlessBilling,PaymentMethod,MonthlyCharges,TotalCharges,Churn
GOOD1,Female,0,Yes,No,1,No,No phone service,DSL,No,Yes,No,No,No,No,Month-to-month,Yes,Electronic check,29.85,29.85,No
BAD,Male,0,No,No,34,Yes,No,DSL,Yes,No,Yes,No,No,No,One year,No,Mailed check,BAD,BAD,No
GOOD2,Female,0,Yes,No,1,No,No phone service,DSL,No,Yes,No,No,No,No,Month-to-month,Yes,Electronic check,29.85,29.85,No
`

func getProfiles(r io.Reader) ([]profiles.ChurnProfile, error) {
	p := make([]profiles.ChurnProfile, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var cp profiles.ChurnProfile
		err := json.Unmarshal(scanner.Bytes(), &cp)
		if err != nil {
			return p, err
		}
		p = append(p, cp)
	}
	if err := scanner.Err(); err != nil {
		return p, err
	}
	return p, nil
}

func TestCsvToJSON(t *testing.T) {
	var b bytes.Buffer
	r := strings.NewReader(goodCsv)

	err := CsvToJSON(r, &b)
	if err != nil {
		t.Fatal(err)
	}
	p, err := getProfiles(bytes.NewReader(b.Bytes()))
	if err != nil {
		t.Fatal("Decoding json failed for some reason:", err)
	}
	if len(p) != 5 {
		t.Fatal("Expected 5 records but got", len(p))
	}
	
	//make sure the header row is not present
	if p[0].CustomerID == "customerID" {
		t.Fatal("First record in json batch appears to be csv header")
	}

	b.Reset()
	r = strings.NewReader(badCsv)

	err = CsvToJSON(r, &b)
	if err != nil {
		t.Fatal(err)
	}
	p, err = getProfiles(bytes.NewReader(b.Bytes()))
	if err != nil {
		log.Println(string(b.Bytes()))
		t.Fatal("Decoding json failed for some reason:", err)
	}
	if len(p) != 2 {
		log.Println("bytes:", string(b.Bytes()))
		t.Fatal("Expected only 2 valid records but got", len(p))
	}
	//make sure the bad row is not present
	for i := range p {
		if p[i].CustomerID == "customerID" {
			t.Fatalf("CSV Header row present as json for some reason: index %d", i)
		}
		if p[i].CustomerID == "BAD" {
			t.Fatalf("Malformed csv record present in json for some reason: index %d", i)
		}
	}
}
