package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _, _ sdk.Msg = &MsgRequestDAS{}, &MsgReportDASResult{}

func (msg *MsgRequestDAS) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgRequestDAS) ValidateBasic() error {
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
