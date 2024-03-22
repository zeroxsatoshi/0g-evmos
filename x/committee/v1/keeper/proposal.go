package keeper

import (
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/evmos/v16/x/committee/v1/types"
)

func (k Keeper) RegisterVoter(ctx sdk.Context, voter sdk.ValAddress, key []byte) error {
	if len(key) != vrf.PublicKeySize {
		return types.ErrInvalidPublicKey
	}

	k.SetVoter(ctx, voter, vrf.PublicKey(key))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterVoter,
			sdk.NewAttribute(types.AttributeKeyVoter, voter.String()),
			// TODO: types.AttributeKeyPublicKey
		),
	)

	return nil
}

func (k Keeper) AddProposal(ctx sdk.Context, proposer sdk.ValAddress, startingHeight uint64, votingStartHeight uint64, votingEndHeight uint64) (uint64, error) {
	// Get a new ID and store the proposal
	proposalID, err := k.StoreNewProposal(ctx, startingHeight, votingStartHeight, votingEndHeight)
	if err != nil {
		return 0, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalSubmit,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyVotingStartHeight, strconv.FormatUint(votingStartHeight, 10)),
			sdk.NewAttribute(types.AttributeKeyVotingEndHeight, strconv.FormatUint(votingEndHeight, 10)),
		),
	)
	return proposalID, nil
}

// AddVote submits a vote on a proposal.
func (k Keeper) AddVote(ctx sdk.Context, proposalID uint64, voter sdk.ValAddress, ballots []*types.Ballot) error {
	// Validate
	pr, found := k.GetProposal(ctx, proposalID)
	if !found {
		return errorsmod.Wrapf(types.ErrUnknownProposal, "%d", proposalID)
	}
	if pr.HasExpiredBy(ctx.BlockHeight()) {
		return errorsmod.Wrapf(types.ErrProposalExpired, "%d ≥ %d", ctx.BlockHeight(), pr.VotingEndHeight)
	}

	// TODO: verify if the voter is registered
	// TODO: verify whether ballots are valid or not

	// Store vote, overwriting any prior vote
	k.SetVote(ctx, types.NewVote(proposalID, voter, ballots))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalVote,
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", pr.ID)),
			sdk.NewAttribute(types.AttributeKeyVoter, voter.String()),
			// TODO: types.AttributeKeyBallots
		),
	)

	return nil
}
