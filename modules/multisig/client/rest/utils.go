package rest

import (
	"encoding/json"
)

//For sign the give  txn
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
	To     string
	Amount string
}

func MsgSpendSignBytes(to string, amount string) ([]byte, error) {
	bz, err := json.Marshal(SendTxBody{
		To:     to,
		Amount: amount,
	})

	if err != nil {
		return nil, err
	}
	return bz, nil
}

