package enrichment

import (
	"bufio"
	"encoding/json"
	"io"
	"math/rand"

	"github.com/pandemicsyn/netlify/churnprofiles"
	log "github.com/sirupsen/logrus"
)

// EnrichedProfile is the original ChurnProfile plus the ChurnScore property
type EnrichedProfile struct {
	churnprofiles.ChurnProfile
	ChurnScore int
}

// CalculateChurnRisk just returns a random int to mock a churn risk score for a given profile
func CalculateChurnRisk(p EnrichedProfile) int {
	return rand.Intn(100)
}

func EnrichAndStore(r io.Reader) error {
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
		// TODO: insert into primary datastore
		log.Infof("Pretend row insert: %v", p)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
