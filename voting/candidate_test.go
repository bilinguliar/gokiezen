package voting

import "testing"

type NastyError error

type MockedRegistry struct {
	AddCandidateFunc    func(name string) error
	RemoveCandidateFunc func(name string) error
}

func (mr *MockedRegistry) AddCandidate(name string) error {
	return mr.AddCandidateFunc(name)
}

func (mr *MockedRegistry) RemoveCandidate(name string) error {
	return mr.RemoveCandidateFunc(name)
}

func TestNewCandidatesCreatesInstanceWithGivenRegistry(t *testing.T) {
	registry := &MockedRegistry{}

	service := NewCandidates(registry)

	if service == nil {
		t.Error("New instance is nil.")
	}

	if service.registry != registry {
		t.Error("Registry is not the one passed to constructor.")
	}
}

func TestAddTrigersRegistryAddCandidate(t *testing.T) {
	var (
		expectedName = "Max"
		added        = false
	)
	registry := &MockedRegistry{
		AddCandidateFunc: func(name string) error {
			if name == expectedName {
				added = true
			}
			return nil
		},
	}

	svc := NewCandidates(registry)
	err := svc.Add(expectedName)
	if err != nil {
		t.Error("Error occurred did not expect that.")
	}

	if !added {
		t.Error("Candidate was not added.")
	}
}

func TestDelTrigersRegistryRemoveCandidate(t *testing.T) {
	var (
		expectedName = "Jonas"
		removed      = false
	)
	registry := &MockedRegistry{
		RemoveCandidateFunc: func(name string) error {
			if name == expectedName {
				removed = true
			}
			return nil
		},
	}

	svc := NewCandidates(registry)
	err := svc.Del(expectedName)
	if err != nil {
		t.Error("Error occurred did not expect that.")
	}

	if !removed {
		t.Error("Candidate was not removed.")
	}
}

func TestAddReturnsErrorIfRegistryDid(t *testing.T) {
	var errWithRegistry NastyError

	registry := &MockedRegistry{
		AddCandidateFunc: func(name string) error {
			return errWithRegistry
		},
	}

	svc := NewCandidates(registry)

	err := svc.Add("anything")

	if err != errWithRegistry {
		t.Error("Error differs from the one we expect.")
	}
}

func TestDelReturnsErrorIfRegistryDid(t *testing.T) {
	var errWithRegistry NastyError

	registry := &MockedRegistry{
		RemoveCandidateFunc: func(name string) error {
			return errWithRegistry
		},
	}

	svc := NewCandidates(registry)

	err := svc.Del("anything")

	if err != errWithRegistry {
		t.Error("Error differs from the one we expect.")
	}
}
