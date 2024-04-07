package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/evmos/v16/x/das/v1/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
}

// NewKeeper creates a new das Keeper instance
func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) SetNextRequestID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextRequestIDKey, types.GetKeyFromID(id))
}

func (k Keeper) GetNextRequestID(ctx sdk.Context) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextRequestIDKey)
	if bz == nil {
		return 0, errorsmod.Wrap(types.ErrInvalidGenesis, "next request ID not set at genesis")
	}
	return types.Uint64FromBytes(bz), nil
}

func (k Keeper) IncrementNextRequestID(ctx sdk.Context) error {
	id, err := k.GetNextRequestID(ctx)
	if err != nil {
		return err
	}
	k.SetNextRequestID(ctx, id+1)
	return nil
}

func (k Keeper) GetDASRequest(ctx sdk.Context, requestID uint64) (types.DASRequest, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RequestKeyPrefix)
	bz := store.Get(types.GetKeyFromID(requestID))
	if bz == nil {
		return types.DASRequest{}, false
	}
	var req types.DASRequest
	k.cdc.MustUnmarshal(bz, &req)
	return req, true
}

func (k Keeper) SetDASRequest(ctx sdk.Context, req types.DASRequest) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RequestKeyPrefix)
	bz := k.cdc.MustMarshal(&req)
	store.Set(types.GetKeyFromID(req.ID), bz)
}

func (k Keeper) IterateDASRequest(ctx sdk.Context, cb func(req types.DASRequest) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.RequestKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var req types.DASRequest
		k.cdc.MustUnmarshal(iterator.Value(), &req)
		if cb(req) {
			break
		}
	}
}

func (k Keeper) GetDASRequests(ctx sdk.Context) []types.DASRequest {
	results := []types.DASRequest{}
	k.IterateDASRequest(ctx, func(req types.DASRequest) bool {
		results = append(results, req)
		return false
	})
	return results
}

func (k Keeper) StoreNewDASRequest(
	ctx sdk.Context,
	batchHeaderHash string,
	numBlobs uint32) (uint64, error) {
	requestID, err := k.GetNextRequestID(ctx)
	if err != nil {
		return 0, err
	}

	req := types.DASRequest{
		ID:              requestID,
		BatchHeaderHash: []byte(batchHeaderHash),
		NumBlobs:        numBlobs,
	}
	k.SetDASRequest(ctx, req)

	return requestID, nil
}

func (k Keeper) GetDASResponse(
	ctx sdk.Context, requestID uint64, sampler sdk.ValAddress,
) (types.DASResponse, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ResponseKeyPrefix)
	bz := store.Get(types.GetResponseKey(requestID, sampler))
	if bz == nil {
		return types.DASResponse{}, false
	}
	var vote types.DASResponse
	k.cdc.MustUnmarshal(bz, &vote)
	return vote, true
}

func (k Keeper) SetDASResponse(ctx sdk.Context, resp types.DASResponse) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ResponseKeyPrefix)
	bz := k.cdc.MustMarshal(&resp)
	store.Set(types.GetResponseKey(resp.ID, resp.Sampler), bz)
}

func (k Keeper) IterateDASResponse(ctx sdk.Context, cb func(resp types.DASResponse) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.ResponseKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var resp types.DASResponse
		k.cdc.MustUnmarshal(iterator.Value(), &resp)
		if cb(resp) {
			break
		}
	}
}

func (k Keeper) GetDASResponses(ctx sdk.Context) []types.DASResponse {
	results := []types.DASResponse{}
	k.IterateDASResponse(ctx, func(resp types.DASResponse) bool {
		results = append(results, resp)
		return false
	})
	return results
}

func (k Keeper) StoreNewDASResponse(
	ctx sdk.Context, requestID uint64, sampler sdk.ValAddress, success bool) error {
	_, found := k.GetDASRequest(ctx, requestID)
	if !found {
		return errorsmod.Wrapf(types.ErrUnknownRequest, "%d", requestID)
	}

	// TODO: verify sampler

	k.SetDASResponse(ctx, types.DASResponse{
		ID:      requestID,
		Sampler: sampler,
		Success: success,
	})

	return nil
}
