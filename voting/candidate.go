package voting

import "log"

// Registry stores all existing candidates. Supports add and delete operations.
type Registry interface {
	AddCandidate(name string) error
	RemoveCandidate(name string) error
}

// CandidatesSvc provides API for candidates.
type CandidatesSvc struct {
	registry Registry
}

// NewCandidates creates new instance with given registry.
func NewCandidates(r Registry) *CandidatesSvc {
	return &CandidatesSvc{
		registry: r,
	}
}

// Add stores single candidate. If already exists - this is not an error.
func (c *CandidatesSvc) Add(name string) error {
	err := c.registry.AddCandidate(name)
	if err != nil {
		log.Println("Candidate add failed, error:", err)
	}

	return err
}

// Del - removes candidate with given name.
func (c *CandidatesSvc) Del(name string) error {
	err := c.registry.RemoveCandidate(name)
	if err != nil {
		log.Println("Candidate was not deleted, error:", err)
	}

	return err
}
