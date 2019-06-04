package churnprofiles

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"log"
	"strconv"
)

// ErrBadCsvCol is returned whenever we encounter a csv row with column thats malformed
var ErrBadCsvCol = errors.New("Encountered malformed csv column in ")

// setProfileValue parses and sets valid csv records in a ChurnProfile - gender and seniorCitizen
// are intentionally omitted. If phone number was present this is where we'd also anonymize the number
// by omitting the last 4 digits.
//
// TODO: return useful error types
func setProfileValue(c *ChurnProfile, field string, value string) error {
	switch field {
	case "customerID":
		c.CustomerID = value
		return nil
	case "Partner":
		c.Partner = value
		return nil
	case "Dependents":
		c.Dependents = value
		return nil
	case "tenure":
		if len(value) == 0 || value == " " {
			c.Tenure = 0
			return nil
		}
		var err error
		c.Tenure, err = strconv.Atoi(value)
		if err != nil {
			return errors.New("Error parsing tenure to int")
		}
		return nil
	case "PhoneService":
		c.PhoneService = value
		return nil
	case "MultipleLines":
		c.MultipleLines = value
		return nil
	case "InternetService":
		c.InternetService = value
		return nil
	case "OnlineSecurity":
		c.OnlineSecurity = value
		return nil
	case "OnlineBackup":
		c.OnlineBackup = value
		return nil
	case "DeviceProtection":
		c.DeviceProtection = value
		return nil
	case "TechSupport":
		c.TechSupport = value
		return nil
	case "StreamingTV":
		c.StreamingTV = value
		return nil
	case "StreamingMovies":
		c.StreamingMovies = value
		return nil
	case "Contract":
		c.Contract = value
		return nil
	case "PaperlessBilling":
		c.PaperlessBilling = value
		return nil
	case "PaymentMethod":
		c.PaymentMethod = value
		return nil
	case "MonthlyCharges":
		if len(value) == 0 || value == " " {
			c.MonthlyCharges = 0
			return nil
		}
		var err error
		c.MonthlyCharges, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.New("Error parsing MonthlyCharges to float")
		}
		return nil
	case "TotalCharges":
		if len(value) == 0 || value == " " {
			c.TotalCharges = 0
			return nil
		}
		var err error
		c.TotalCharges, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.New("Error parsing TotalCharges to float")
		}
		return nil
	}
	return nil
}

func loadProfile(p *ChurnProfile, record []string, headerFields []string) error {
	for n, value := range record {
		err := setProfileValue(p, headerFields[n], value)
		if err != nil {
			// TODO: handle bad csv rows, whether logging, adding to skip file
			// emitting metrics/notification etc, for now just skip this row
			log.Printf("Encountered bad record: %v", record)
			log.Printf("Record err: %v", err)
			return ErrBadCsvCol
		}
	}
	return nil
}

// CsvToJSON converts a stream of churn model records to json
func CsvToJSON(csvSrc io.Reader, jsonDst io.Writer) error {
	enc := json.NewEncoder(jsonDst)
	r := csv.NewReader(csvSrc)
	n := 0
	// TODO: should probably compare header to expected schema
	var headerFields []string
	for {
		p := &ChurnProfile{}
		var err error
		var record []string
		if record, err = r.Read(); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if n == 0 {
			headerFields = record
		} else {
			err := loadProfile(p, record, headerFields)
			if err == ErrBadCsvCol {
				// TODO: handle bad csv rows, whether logging, adding to skip file
				// emitting metrics/notification etc, for now just skip this row
				log.Println("Encountered bad column, skipping row")
				continue
			}
			if err := enc.Encode(&p); err != nil {
				return err
			}
		}
		n++
		log.Println("seen", n)
	}
	//TODO: emit metrics around lines processed, skipped, etc
	return nil
}
