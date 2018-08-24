package rest

import (
	"net/http"

	"encoding/json"
	"reflect"
	sdk "sentinel/modules/multisig/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	stypes "sentinel/modules/multisig/types"
)

type MsgSpendFromMultiSig struct {
	Spend    stypes.StdtxSpend `json:"spend"`
	To       string     `json:"to"`
	From     string     `json:"from"`
	Amount   string     `json:"amount"`
	Password string     `json:"password"`
}

func multisignatureSpendHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg MsgSpendFromMultiSig
		var kb keys.Keybase
		var output stypes.StdtxSpend

		// Decoinding the Request
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MsgSpendFromMultiSigResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while decoding the request body",
				},
			})
			return
		}

		// Fetching the keybase
		kb, err := ckeys.GetKeyBase()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while Fetching the keybase",
				},
			})
			return
		}

		//check the  existing txxbytes
		if (reflect.DeepEqual(stypes.StdtxSpend{}, msg.Spend)) {
			bz, err := MsgSpendSignBytes(msg.To, msg.Amount)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Error occurred while marshal ",
					},
				})
				return
			}
			_, pubkey, err := kb.Sign(msg.From, msg.Password, bz)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Error occurred while creating sign bytes",
					},
				})
				return
			}

			pubkeystr, _ := types.Bech32ifyAccPub(pubkey)
			output = stypes.StdtxSpend{
				To:     msg.To,
				Amount: msg.Amount,
				Pubkey: append(output.Pubkey, pubkeystr),
			}

			data, _ := json.Marshal(output)
			w.Write(data)
			return
		}

		bz, err := MsgSpendSignBytes(msg.To, msg.Amount)

		_, pubkey, err := kb.Sign(msg.From, msg.Password, bz)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while creating sign bytes ",
				},
			})
			return
		}

		pubkeystr, _ := types.Bech32ifyAccPub(pubkey)

		for _, value := range msg.Spend.Pubkey {
			if value == pubkeystr {
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Error occurred the given publickey already exist",
					},
				})
				return
			}
		}

		msg.Spend.Pubkey = append(msg.Spend.Pubkey, pubkeystr)
		output = stypes.NewStdtxSpend(msg.Spend.To, msg.Spend.Amount, msg.Spend.Pubkey)

		data, _ := json.Marshal(output)
		w.Write(data)
		return

	}
}
