package multisig

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/types"
	types2 "sentinel/modules/multisig/types"
)

type Stdtx struct {
	MinKeys   uint8    `json:"min_keys"`
	TotalKeys uint8    `json:"total_keys"`
	Order     bool     `json:"order"`
	Pubkey    []string `json:"pubkey"`
	Counter   int64    `json:"counter"`
	Signature [][]byte  `json:"signature"`
}

type MsgCreateMultiSigAddress struct {
	Txbytes Stdtx `json:"txbytes"`
	Address  types.AccAddress `json:"from"`
}

func NewMsgCreateMultiSigAddress(txbytes Stdtx,address types.AccAddress) MsgCreateMultiSigAddress{
	return MsgCreateMultiSigAddress{
		Txbytes:txbytes,
		Address:address,
	}
}



func (msg MsgCreateMultiSigAddress) Type() string {
	return "multisig"
}

func (msg MsgCreateMultiSigAddress) GetSignBytes() []byte {
	byte_format, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return byte_format
}

func (msg MsgCreateMultiSigAddress) ValidateBasic() types.Error {

	return nil
}
func (msg MsgCreateMultiSigAddress) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.Address}
}


type MsgFundMultiSig struct{
	To types.AccAddress
	Address types.AccAddress
	Amount types.Coins

}

func NewMsgFundMultiSig(to types.AccAddress,from types.AccAddress,coins types.Coins) MsgFundMultiSig{
	return MsgFundMultiSig{
		To:to,
		Address:from,
		Amount:coins,
	}
}

func (msg MsgFundMultiSig) Type() string {
	return "multisig"
}

func (msg MsgFundMultiSig) GetSignBytes() []byte {
	byte_format, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return byte_format
}

func (msg MsgFundMultiSig) ValidateBasic() types.Error {

	return nil
}
func (msg MsgFundMultiSig) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.Address}
}




type MsgSendFromMultiSig struct {
	Txbytes types2.StdtxSpend
	Address types.AccAddress
}

func NewMsgSendFromMultiSig(txbytes types2.StdtxSpend,address types.AccAddress) MsgSendFromMultiSig{
		return MsgSendFromMultiSig{
			Txbytes:txbytes,
			Address:address,
		}
}

func (msg MsgSendFromMultiSig) Type() string {
	return "multisig"
}

func (msg MsgSendFromMultiSig) GetSignBytes() []byte {
	byte_format, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return byte_format
}

func (msg MsgSendFromMultiSig) ValidateBasic() types.Error {

	return nil
}
func (msg MsgSendFromMultiSig) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.Address}
}
