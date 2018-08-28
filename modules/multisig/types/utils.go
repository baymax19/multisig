package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/types"
)

//For add signature and publickeys to structure
type Stdtx struct {
	MinKeys   uint8    `json:"min_keys"`
	TotalKeys uint8    `json:"total_keys"`
	Order     bool     `json:"order"`
	Pubkey    []string `json:"pubkey"`
	Counter   int64    `json:"counter"`
	Signature [][]byte `json:"signature"`
}

func NewStdtx(order bool, totalkeys uint8, minkeys uint8, pubkey []string, count int64, sign [][]byte) Stdtx {
	data := Stdtx{
		MinKeys:   minkeys,
		TotalKeys: totalkeys,
		Order:     order,
		Pubkey:    pubkey,
		Counter:   count,
		Signature: sign,
	}
	return data
}

//For add publikeys to structure
type StdtxSpend struct {
	To              types.AccAddress `json:"to"`
	MultiSigAddress types.AccAddress `json:"multi_sig_address`
	Amount          types.Coins      `json:"amount"`
	Signature       [][]byte         `json:"signature"`
	TxNumber        int64            `json:"tx_number"`
}

func NewStdtxSpend(to types.AccAddress, maddress types.AccAddress, amount types.Coins, sign [][]byte, tx_number int64) StdtxSpend {
	data := StdtxSpend{
		To:              to,
		MultiSigAddress: maddress,
		Amount:          amount,
		Signature:       sign,
		TxNumber:        tx_number,
	}
	return data
}

func Intersection(txbytesfromchain, txbytesfromuser []string) (c []string) {
	hash := make(map[string]bool)

	for _, value := range txbytesfromchain {
		hash[value] = true
	}

	for _, value := range txbytesfromuser {
		if _, ok := hash[value]; ok {
			c = append(c, value)
		}
	}
	return
}

type StdTxSend struct {
	To        types.AccAddress `json:"to"`
	Amount    types.Coins      `json:"amount"`
	Pubkey    []string         `json:"pubkey"`
	Signature [][]byte         `json:"signature"`
}

type StdSig struct {
	MinKeys   uint8 `json:"min_keys,omitempty"`
	TotalKeys uint8 `json:"total_keys,omitempty"`
	Order     bool  `json:"order,omitempty"`
}

func CreateSignBytes(minkey uint8, order bool, totalkeys uint8) ([]byte, error) {
	bz, err := json.Marshal(StdSig{
		MinKeys:   minkey,
		TotalKeys: totalkeys,
		Order:     order,
	})

	if err != nil {
		return nil, err
	}
	return bz, nil
}

//For sign the Spend txn
type SendTxBody struct {
	To              types.AccAddress
	Amount          types.Coins
	MultiSigAddress types.AccAddress
}

func MsgSpendSignBytes(to types.AccAddress, amount types.Coins, maddress types.AccAddress) ([]byte, error) {
	bz, err := json.Marshal(SendTxBody{
		To:              to,
		Amount:          amount,
		MultiSigAddress: maddress,
	})

	if err != nil {
		return nil, err
	}
	return bz, nil
}
