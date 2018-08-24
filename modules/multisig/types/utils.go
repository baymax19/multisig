package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	"encoding/json"
)

//For add signature and publickeys to structure
type Stdtx struct {
	MinKeys   uint8    `json:"min_keys"`
	TotalKeys uint8    `json:"total_keys"`
	Order     bool     `json:"order"`
	Pubkey    []string `json:"pubkey"`
	Counter   int64    `json:"counter"`
	Signature  [][]byte   `json:"signature"`
}

func NewStdtx(order bool, totalkeys uint8, minkeys uint8, pubkey []string, count int64,sign [][]byte) Stdtx {
	data := Stdtx{
		MinKeys:   minkeys,
		TotalKeys: totalkeys,
		Order:     order,
		Pubkey:    pubkey,
		Counter:   count,
		Signature:sign,
	}
	return data
}


//For add publikeys to structure
type StdtxSpend struct {
	To     types.AccAddress   `json:"to"`
	Amount types.Coins   `json:"amount"`
	Signature [][]byte `json:"signature"`
}

func NewStdtxSpend(to types.AccAddress, amount types.Coins, sign [][]byte) StdtxSpend {
	data := StdtxSpend{
		To:     to,
		Amount: amount,
		Signature:sign,
	}
	return data
}


func Intersection(txbytesfromchain,  txbytesfromuser []string) (c []string) {
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

type StdTxSend struct{
	To types.AccAddress `json:"to"`
	Amount types.Coins `json:"amount"`
	Pubkey []string `json:"pubkey"`
	Signature [][]byte `json:"signature"`
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
	To    types.AccAddress
	Amount types.Coins
}

func MsgSpendSignBytes(to types.AccAddress, amount types.Coins) ([]byte, error) {
	bz, err := json.Marshal(SendTxBody{
		To:     to,
		Amount: amount,
	})

	if err != nil {
		return nil, err
	}
	return bz, nil
}

