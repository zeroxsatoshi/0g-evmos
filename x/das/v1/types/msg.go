package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _, _ sdk.Msg = &MsgRequestDAS{}, &MsgReportDASResult{}

func NewMsgRequestDAS(fromAddr sdk.AccAddress, hash string, numBlobs uint32) *MsgRequestDAS {
	return &MsgRequestDAS{
		Requester:       fromAddr.String(),
		BatchHeaderHash: hash,
		NumBlobs:        numBlobs,
	}
}

func (msg MsgRequestDAS) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Requester)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (msg MsgRequestDAS) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Requester)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	return nil
}

func (msg *MsgReportDASResult) GetSigners() []sdk.AccAddress {
	sampler, err := sdk.AccAddressFromBech32(msg.Sampler)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sampler}
}

func (msg *MsgReportDASResult) ValidateBasic() error {
	return nil
}
