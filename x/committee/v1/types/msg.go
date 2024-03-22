package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _, _, _ sdk.Msg = &MsgRegister{}, &MsgPropose{}, &MsgVote{}

// GetSigners returns the expected signers for a MsgRegister message.
func (msg *MsgRegister) GetSigners() []sdk.AccAddress {
	initiator, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{initiator}
}

// ValidateBasic does a sanity check of the provided data
func (msg *MsgRegister) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if _, err := sdk.ValAddressFromBech32(msg.Voter); err != nil {
		return ErrInvalidValidatorAddress
	}
	if len(msg.Key) != vrf.PublicKeySize {
		return ErrInvalidPublicKey
	}
	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgRegister) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgPropose message.
func (msg *MsgPropose) GetSigners() []sdk.AccAddress {
	initiator, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{initiator}
}

// ValidateBasic does a sanity check of the provided data
func (msg *MsgPropose) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if _, err := sdk.ValAddressFromBech32(msg.Proposer); err != nil {
		return ErrInvalidValidatorAddress
	}
	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgPropose) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgVote message.
func (msg *MsgVote) GetSigners() []sdk.AccAddress {
	initiator, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{initiator}
}

// ValidateBasic does a sanity check of the provided data
func (msg *MsgVote) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if _, err := sdk.ValAddressFromBech32(msg.Voter); err != nil {
		return ErrInvalidValidatorAddress
	}
	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgVote) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}
