package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/evmos/v16/x/committee/v1/types"
)

// Keeper of the inflation store
type Keeper struct {
	storeKey      storetypes.StoreKey
	cdc           codec.BinaryCodec
	stakingKeeper types.StakingKeeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	stakingKeeper types.StakingKeeper,
) Keeper {
	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		stakingKeeper: stakingKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// ------------------------------------------
//				Proposals
// ------------------------------------------

// SetNextProposalID stores an ID to be used for the next created proposal
func (k Keeper) SetNextProposalID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextProposalIDKey, types.GetKeyFromID(id))
}

// GetNextProposalID reads the next available global ID from store
func (k Keeper) GetNextProposalID(ctx sdk.Context) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextProposalIDKey)
	if bz == nil {
		return 0, errorsmod.Wrap(types.ErrInvalidGenesis, "next proposal ID not set at genesis")
	}
	return types.Uint64FromBytes(bz), nil
}

// IncrementNextProposalID increments the next proposal ID in the store by 1.
func (k Keeper) IncrementNextProposalID(ctx sdk.Context) error {
	id, err := k.GetNextProposalID(ctx)
	if err != nil {
		return err
	}
	k.SetNextProposalID(ctx, id+1)
	return nil
}

// StoreNewProposal stores a proposal, adding a new ID
func (k Keeper) StoreNewProposal(ctx sdk.Context, startingHeight uint64, votingStartHeight uint64, votingEndHeight uint64) (uint64, error) {
	newProposalID, err := k.GetNextProposalID(ctx)
	if err != nil {
		return 0, err
	}
	proposal, err := types.NewProposal(
		newProposalID,
		startingHeight,
		votingStartHeight,
		votingEndHeight,
	)
	if err != nil {
		return 0, err
	}

	k.SetProposal(ctx, proposal)

	err = k.IncrementNextProposalID(ctx)
	if err != nil {
		return 0, err
	}
	return newProposalID, nil
}

// GetProposal gets a proposal from the store.
func (k Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (types.Proposal, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefix)
	bz := store.Get(types.GetKeyFromID(proposalID))
	if bz == nil {
		return types.Proposal{}, false
	}
	var proposal types.Proposal
	k.cdc.MustUnmarshal(bz, &proposal)
	return proposal, true
}

// SetProposal puts a proposal into the store.
func (k Keeper) SetProposal(ctx sdk.Context, proposal types.Proposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefix)
	bz := k.cdc.MustMarshal(&proposal)
	store.Set(types.GetKeyFromID(proposal.ID), bz)
}

// DeleteProposal removes a proposal from the store.
func (k Keeper) DeleteProposal(ctx sdk.Context, proposalID uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefix)
	store.Delete(types.GetKeyFromID(proposalID))
}

// IterateProposals provides an iterator over all stored proposals.
// For each proposal, cb will be called. If cb returns true, the iterator will close and stop.
func (k Keeper) IterateProposals(ctx sdk.Context, cb func(proposal types.Proposal) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.ProposalKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proposal types.Proposal
		k.cdc.MustUnmarshal(iterator.Value(), &proposal)
		if cb(proposal) {
			break
		}
	}
}

// GetProposals returns all stored proposals.
func (k Keeper) GetProposals(ctx sdk.Context) types.Proposals {
	results := types.Proposals{}
	k.IterateProposals(ctx, func(prop types.Proposal) bool {
		results = append(results, prop)
		return false
	})
	return results
}

// DeleteProposalAndVotes removes a proposal and its associated votes.
func (k Keeper) DeleteProposalAndVotes(ctx sdk.Context, proposalID uint64) {
	votes := k.GetVotesByProposal(ctx, proposalID)
	k.DeleteProposal(ctx, proposalID)
	for _, v := range votes {
		k.DeleteVote(ctx, v.ProposalID, v.Voter)
	}
}

// ------------------------------------------
//				Votes
// ------------------------------------------

// GetVote gets a vote from the store.
func (k Keeper) GetVote(ctx sdk.Context, proposalID uint64, voter sdk.ValAddress) (types.Vote, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.VoteKeyPrefix)
	bz := store.Get(types.GetVoteKey(proposalID, voter))
	if bz == nil {
		return types.Vote{}, false
	}
	var vote types.Vote
	k.cdc.MustUnmarshal(bz, &vote)
	return vote, true
}

// SetVote puts a vote into the store.
func (k Keeper) SetVote(ctx sdk.Context, vote types.Vote) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.VoteKeyPrefix)
	bz := k.cdc.MustMarshal(&vote)
	store.Set(types.GetVoteKey(vote.ProposalID, vote.Voter), bz)
}

// DeleteVote removes a Vote from the store.
func (k Keeper) DeleteVote(ctx sdk.Context, proposalID uint64, voter sdk.ValAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.VoteKeyPrefix)
	store.Delete(types.GetVoteKey(proposalID, voter))
}

// IterateVotes provides an iterator over all stored votes.
// For each vote, cb will be called. If cb returns true, the iterator will close and stop.
func (k Keeper) IterateVotes(ctx sdk.Context, cb func(vote types.Vote) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.VoteKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		k.cdc.MustUnmarshal(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// GetVotes returns all stored votes.
func (k Keeper) GetVotes(ctx sdk.Context) []types.Vote {
	results := []types.Vote{}
	k.IterateVotes(ctx, func(vote types.Vote) bool {
		results = append(results, vote)
		return false
	})
	return results
}

// GetVotesByProposal returns all votes for one proposal.
func (k Keeper) GetVotesByProposal(ctx sdk.Context, proposalID uint64) []types.Vote {
	results := []types.Vote{}
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), append(types.VoteKeyPrefix, types.GetKeyFromID(proposalID)...))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		k.cdc.MustUnmarshal(iterator.Value(), &vote)
		results = append(results, vote)
	}

	return results
}

// ------------------------------------------
//				Voter
// ------------------------------------------

func (k Keeper) SetVoter(ctx sdk.Context, voter sdk.ValAddress, pk vrf.PublicKey) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.VoterKeyPrefix)
	store.Set(types.GetVoterKey(voter), pk)
}
