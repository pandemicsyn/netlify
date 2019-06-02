package churnprofiles

import (
	"bufio"
	"io"
	"log"
	"os"
	"reflect"
	"testing"
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
		c := &ChurnProfile{}
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
		c := &ChurnProfile{}
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
		c := &ChurnProfile{}
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
		c := &ChurnProfile{}
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
		c := &ChurnProfile{}
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
		c := &ChurnProfile{}
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
		c := &ChurnProfile{}
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
		c := &ChurnProfile{}
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
		c := &ChurnProfile{}
		err := setProfileValue(c, field, floatBad)
		if err == nil {
			t.Fatalf("expected error for field '%v', sent string", structField)
		}
	}

	// test non existent field
	c := &ChurnProfile{}
	err := setProfileValue(c, "notpresent", "novalue")
	if err != nil {
		t.Fatal(err)
	}
}
