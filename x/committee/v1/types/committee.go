package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"sigs.k8s.io/yaml"
)

type Proposals []Proposal

// NewProposal instantiates a new instance of Proposal
func NewProposal(id uint64, startingHeight uint64, votingStartHeight uint64, votingEndHeight uint64) (Proposal, error) {
	return Proposal{
		ID:                   id,
		StartingHeight:       startingHeight,
		VotingStartingHeight: votingStartHeight,
		VotingEndHeight:      votingEndHeight,
	}, nil
}

// String implements the fmt.Stringer interface.
func (p Proposal) String() string {
	bz, _ := yaml.Marshal(p)
	return string(bz)
}

// HasExpiredBy calculates if the proposal will have expired by a certain height.
// All votes must be cast before deadline, those cast at time == deadline are not valid
func (p Proposal) HasExpiredBy(height int64) bool {
	return height >= int64(p.VotingEndHeight)
}

// NewVote instantiates a new instance of Vote
func NewVote(proposalID uint64, voter sdk.ValAddress, ballots []*Ballot) Vote {
	return Vote{
		ProposalID: proposalID,
		Voter:      voter,
		Ballots:    ballots,
	}
}

// Validates Vote fields
func (v Vote) Validate() error {
	if v.Voter.Empty() {
		return fmt.Errorf("voter address cannot be empty")
	}

	// TODO: validate ballots
	return nil
}
