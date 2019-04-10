package shell

import (
	"context"
)

// PeersList ...
type PeersList struct {
	Peers []string
}

// BootstrapAdd ...
func (s *Shell) BootstrapAdd(peers []string) ([]string, error) {
	var addOutput PeersList
	err := s.Request("bootstrap/add", peers...).Exec(context.Background(), &addOutput)
	return addOutput.Peers, err
}

// BootstrapAddDefault ...
func (s *Shell) BootstrapAddDefault() ([]string, error) {
	var addOutput PeersList
	err := s.Request("bootstrap/add/default").Exec(context.Background(), &addOutput)
	return addOutput.Peers, err
}

// BootstrapRmAll ...
func (s *Shell) BootstrapRmAll() ([]string, error) {
	var rmAllOutput PeersList
	err := s.Request("bootstrap/rm/all").Exec(context.Background(), &rmAllOutput)
	return rmAllOutput.Peers, err
}
