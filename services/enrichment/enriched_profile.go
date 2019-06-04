package enrichment

import (
	"bufio"
	"encoding/json"
	"io"
	"math/rand"
	"time"

	"github.com/pandemicsyn/netlify/churnprofiles"
	"github.com/pkg/errors"
)

// EnrichedProfile is the original ChurnProfile plus the ChurnScore property
// an EnrichedAt field is also included so that you can see changes for the profiles over time
type EnrichedProfile struct {
	churnprofiles.ChurnProfile
	ChurnScore int
	EnrichedAt time.Time
}

// CalculateChurnRisk just returns a random int to mock a churn risk score for a given profile
func CalculateChurnRisk(p EnrichedProfile) int {
	return rand.Intn(100)
}

//EnrichAndStore reads in a .json blog of ChurnProfile's and then enriches them with a ChurnScore
// and stores them in Enriched Profile Store
func EnrichAndStore(r io.Reader, store EProfileStore) error {
	profiles := make([]EnrichedProfile, 0)
	scanner := bufio.NewScanner(r)
	// TODO: should try to batch these
	for scanner.Scan() {
		var p EnrichedProfile
		err := json.Unmarshal(scanner.Bytes(), &p)
		if err != nil {
			return err
		}
		p.ChurnScore = CalculateChurnRisk(p)
		if p.CustomerID != "" {
			profiles = append(profiles, p)
		}
		p.EnrichedAt = time.Now()
		profiles = append(profiles, p)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	err := store.BulkSave(&profiles)
	if err != nil {
		errors.Wrap(err, "failed to bulk save profiles")
	}
	return nil
}
