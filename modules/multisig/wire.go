package multisig

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCreateMultiSigAddress{}, "sentinel/multisig_create_address", nil)
	cdc.RegisterConcrete(MsgFundMultiSig{}, "sentinel/fund_to _multisig_address", nil)
	cdc.RegisterConcrete(MsgSendFromMultiSig{}, "sentinel/send_from _multisig_address", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
