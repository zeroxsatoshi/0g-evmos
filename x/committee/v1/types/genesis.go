package types

import (
	"fmt"
)

// DefaultNextProposalID is the starting point for proposal IDs.
const DefaultNextProposalID uint64 = 1

// NewGenesisState returns a new genesis state object for the module.
func NewGenesisState(nextProposalID uint64, proposals Proposals, votes []Vote) *GenesisState {
	return &GenesisState{
		NextProposalID: nextProposalID,
		Proposals:      proposals,
		Votes:          votes,
	}
}

// DefaultGenesisState returns the default genesis state for the module.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		DefaultNextProposalID,
		Proposals{},
		[]Vote{},
	)
}

// Validate performs basic validation of genesis data.
func (gs GenesisState) Validate() error {
	// validate proposals
	proposalMap := make(map[uint64]bool, len(gs.Proposals))
	for _, p := range gs.Proposals {
		// check there are no duplicate IDs
		if _, ok := proposalMap[p.ID]; ok {
			return fmt.Errorf("duplicate proposal ID found in genesis state; id: %d", p.ID)
		}
		proposalMap[p.ID] = true

		// validate next proposal ID
		if p.ID >= gs.NextProposalID {
			return fmt.Errorf("NextProposalID is not greater than all proposal IDs; id: %d", p.ID)
		}
	}

	// validate votes
	for _, v := range gs.Votes {
		// validate committee
		if err := v.Validate(); err != nil {
			return err
		}

		// check proposal exists
		if !proposalMap[v.ProposalID] {
			return fmt.Errorf("vote refers to non existent proposal; vote: %+v", v)
		}
	}
	return nil
}
