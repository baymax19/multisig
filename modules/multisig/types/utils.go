package types



//For add signature and publickeys to structure
type Stdtx struct {
	MinKeys   uint8    `json:"min_keys"`
	TotalKeys uint8    `json:"total_keys"`
	Order     bool     `json:"order"`
	Pubkey    []string `json:"pubkey"`
	Counter   int64    `json:"counter"`
}

func NewStdtx(order bool, totalkeys uint8, minkeys uint8, pubkey []string, count int64) Stdtx {
	data := Stdtx{
		MinKeys:   minkeys,
		TotalKeys: totalkeys,
		Order:     order,
		Pubkey:    pubkey,
		Counter:   count,
	}
	return data
}


//For add publikeys to structure
type StdtxSpend struct {
	To     string   `json:"to"`
	Amount string   `json:"amount"`
	Pubkey []string `json:"pubkey"`
}

func NewStdtxSpend(to string, amount string, pubkey []string) StdtxSpend {
	data := StdtxSpend{
		To:     to,
		Amount: amount,
		Pubkey: pubkey,
	}
	return data
}
