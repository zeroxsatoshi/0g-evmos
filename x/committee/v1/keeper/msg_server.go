// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/evmos/evmos/v16/x/committee/v1/types"
)

var _ types.MsgServer = &Keeper{}

// Register handles MsgRegister messages
func (k Keeper) Register(goCtx context.Context, msg *types.MsgRegister) (*types.MsgRegisterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.Voter)
	if err != nil {
		return nil, err
	}

	fmt.Printf("valAddr: %s\n", valAddr.String())

	_, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	if err := k.AddVoter(ctx, valAddr, msg.Key); err != nil {
		return nil, err
	}

	return &types.MsgRegisterResponse{}, nil
}

// Vote handles MsgVote messages
func (k Keeper) Vote(goCtx context.Context, msg *types.MsgVote) (*types.MsgVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	voter, err := sdk.ValAddressFromBech32(msg.Voter)
	if err != nil {
		return nil, err
	}

	if err := k.AddVote(ctx, msg.CommitteeID, voter, msg.Ballots); err != nil {
		return nil, err
	}

	// ctx.EventManager().EmitEvent(
	// 	sdk.NewEvent(
	// 		sdk.EventTypeMessage,
	// 		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	// 		sdk.NewAttribute(sdk.AttributeKeySender, msg.Voter),
	// 	),
	// )

	return &types.MsgVoteResponse{}, nil
}
