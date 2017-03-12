package voting

import (
	"log"
	"strings"
)

const (
	unresolved = "N/A"
)

// StatItem holds counter name and current read.
type StatItem struct {
	Name  string
	Value int
}

// Stats holds collections of StatItems.
type Stats struct {
	Candidates []StatItem
	Countries  []StatItem
}

// Voting is a service that holds all business logic required to run voting.
type Voting struct {
	messenger Messenger
	enquirer  Enquirer
	scoreKpr  ScoreKeeper
	event     string
}

// Meesenger is used to send text messages.
type Messenger interface {
	RequestSMS(sender, msisdn, text string)
}

// Enquirer is used to resolve Country by MSISDN.
type Enquirer interface {
	Lookup(msisdn string) (string, error)
}

// ScoreKeeper persists score and stats, returns results.
type ScoreKeeper interface {
	AddPoint(participant string) error
	AddCountry(name string) error
	GetAllCandidates() ([]string, error)
	GetAllCountries() ([]string, error)
	Get(key string) (int, error)
}

// New constructs Voting service instance initialized with all dependencies.
func New(m Messenger, en Enquirer, sk ScoreKeeper, ev string) *Voting {
	return &Voting{
		messenger: m,
		enquirer:  en,
		scoreKpr:  sk,
		event:     ev,
	}
}

// RegisterVote increments votes counter for participant and also keeps track of number of votes for each country.
func (s *Voting) RegisterVote(msisdn, cand string) error {
	log.Printf("Got new message: %q from MSISDN: %q", cand, msisdn)
	var (
		country string
		err     error
	)

	cand = strings.TrimSpace(cand)
	if cand == "" {
		log.Println("Voter sent blank SMS, score not changed.")
		s.messenger.RequestSMS(s.event, msisdn, "Please specify candidate's name to actually vote.")
		return nil
	}

	err = s.scoreKpr.AddPoint(cand)
	if err != nil {
		log.Println("Point was not added to participant's score, error:", err)
		return err
	}

	country, err = s.enquirer.Lookup(msisdn)
	if err != nil {
		log.Printf("Country lookup failed for MSISDN: %q, error: %q", msisdn, err)
		country = unresolved
	}

	err = s.scoreKpr.AddCountry(country)
	if err != nil {
		log.Printf("Failed to assure contry with code: %q is present in countries set. Error: %q", country, err)
	}

	// Here just recording stats, so if failed - no big deal, candidates's vote is there already.
	if err = s.scoreKpr.AddPoint(country); err != nil {
		log.Println("Country counter was not incremented, error:", err)
	}

	s.messenger.RequestSMS(s.event, msisdn, "Thanks for your vote!")

	return nil
}

// GetStats returns voting statistics for each participant and distribution by countries.
func (s *Voting) GetStats() (Stats, error) {
	candidates, err := s.scoreKpr.GetAllCandidates()
	if err != nil {
		log.Println("Failed to retrieve set of all candidates, error:", err)
		return Stats{}, err
	}

	countries, err := s.scoreKpr.GetAllCountries()
	if err != nil {
		log.Println("Failed to retrieve set of all countries, error:", err)
		return Stats{}, err
	}

	return Stats{
		Candidates: s.populateStatItems(candidates),
		Countries:  s.populateStatItems(countries),
	}, nil
}

func (s *Voting) populateStatItems(keys []string) []StatItem {
	results := make([]StatItem, 0, len(keys))

	for _, k := range keys {
		v, err := s.scoreKpr.Get(k)
		if err != nil {
			log.Printf("Failed to get score for key: %q, error: %q", k, err)
			// handle -1 as temporary unresolvable on client,
			// most likely we will get proper value during next update.
			v = -1
		}
		results = append(results, StatItem{Name: k, Value: v})
	}

	return results
}
