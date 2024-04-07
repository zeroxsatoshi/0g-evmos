package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/evmos/v16/x/das/v1/types"
)

var _ types.MsgServer = &Keeper{}

// RequestDAS handles MsgRequestDAS messages
func (k Keeper) RequestDAS(
	goCtx context.Context, msg *types.MsgRequestDAS,
) (*types.MsgRequestDASResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	requestID, err := k.StoreNewDASRequest(ctx, msg.BatchHeaderHash, msg.NumBlobs)
	if err != nil {
		return nil, err
	}

	return &types.MsgRequestDASResponse{
		RequestID: requestID,
	}, nil
}

// ReportDASResult handles MsgReportDASResult messages
func (k Keeper) ReportDASResult(
	goCtx context.Context, msg *types.MsgReportDASResult,
) (*types.MsgReportDASResultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sampler, err := sdk.ValAddressFromBech32(msg.Sampler)
	if err != nil {
		return nil, err
	}

	if err := k.StoreNewDASResponse(ctx, msg.RequestID, sampler, msg.Success); err != nil {
		return nil, err
	}

	return &types.MsgReportDASResultResponse{}, nil
}
