package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/evmos/v16/x/committee/v1/types"
)

var _ types.QueryServer = Keeper{}

// Period returns the current period of the inflation module.
func (k Keeper) NextProposalID(
	c context.Context,
	_ *types.QueryNextProposalIDRequest,
) (*types.QueryNextProposalIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	nextProposalID, err := k.GetNextProposalID(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryNextProposalIDResponse{NextProposalID: nextProposalID}, nil
}
