package voting

import "log"

type Registry interface {
	AddCandidate(name string) error
	RemoveCandidate(name string) error
}

type CandidatesSvc struct {
	registry Registry
}

func NewCandidates(r Registry) *CandidatesSvc {
	return &CandidatesSvc{
		registry: r,
	}
}

func (c *CandidatesSvc) Add(name string) error {
	err := c.registry.AddCandidate(name)
	if err != nil {
		log.Println("Candidate add failed, error:", err)
	}

	return err
}

func (c *CandidatesSvc) Del(name string) error {
	err := c.registry.RemoveCandidate(name)
	if err != nil {
		log.Println("Candidate was not deleted, error:", err)
	}

	return err
}
