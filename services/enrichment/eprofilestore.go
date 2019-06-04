package enrichment

import (
	"database/sql"

	"github.com/lib/pq"
)

//EProfileStore is how we store the EnrichedProfiles
type EProfileStore interface {
	BulkSave(*[]EnrichedProfile) error
}

//PGEProfileStore is a Postgres backed EProfileStore
type PGEProfileStore struct {
	db *sql.DB
}

func NewPGEProfileStore(db *sql.DB) EProfileStore {
	return &PGEProfileStore{db}
}

// BulkSave uses a postgres Copy In call to try an efficiently store our enriched profiles
// I'm assuming that its ok to have multiple records for the same customer (to see their change over time)
// If we needed to have a unique constraint on CustomerID we'd have to switch models. Perhaps to smaller batched
// transactions where can first see if CustomerID exists and then insert or upsert accordingly.
func (s *PGEProfileStore) BulkSave(profiles *[]EnrichedProfile) error {
	txn, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := txn.Prepare(pq.CopyIn("enriched_profiles", "CustomerID",
		"Partner",
		"Dependents",
		"Tenure",
		"PhoneService",
		"MultipleLines",
		"InternetService",
		"OnlineSecurity",
		"OnlineBackup",
		"DeviceProtection",
		"TechSupport",
		"StreamingTV",
		"StreamingMovies",
		"Contract",
		"PaperlessBilling",
		"PaymentMethod",
		"MonthlyCharges",
		"TotalCharges",
		"ChurnScore"))

	if err != nil {
		return err
	}

	for _, p := range *profiles {
		_, err = stmt.Exec(p.CustomerID,
			p.Partner,
			p.Dependents,
			p.Tenure,
			p.PhoneService,
			p.MultipleLines,
			p.InternetService,
			p.OnlineSecurity,
			p.OnlineBackup,
			p.DeviceProtection,
			p.TechSupport,
			p.StreamingTV,
			p.StreamingMovies,
			p.Contract,
			p.PaperlessBilling,
			p.PaymentMethod,
			p.MonthlyCharges,
			p.TotalCharges,
			p.ChurnScore,
		)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}
