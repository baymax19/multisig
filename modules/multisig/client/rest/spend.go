package rest

import (
	"net/http"

	"encoding/json"
	sdk "sentinel/modules/multisig/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	stypes "sentinel/modules/multisig/types"
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/types"
)

type MsgSpendFromMultiSig struct {
	Spend    string `json:"spend"`
	To       string     `json:"to"`
	From     string     `json:"from"`
	Amount   string     `json:"amount"`
	Password string     `json:"password"`
}

func multisignatureSpendHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg MsgSpendFromMultiSig
		var kb keys.Keybase
		var Txbytes stypes.StdtxSpend

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
		if (msg.Spend=="") {

			if(msg.Amount==""){
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Enter Amount to send ",
					},
				})
				return
			}

			if(msg.To==""){
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Enter Ato account",
					},
				})
				return
			}

			to,err:=types.AccAddressFromBech32(msg.To)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"address from bech64 string is failed ",
					},
				})
				return
			}

			coins,err:=types.ParseCoins(msg.Amount)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Error occurred parse coins ",
					},
				})
				return
			}
			bz, err := stypes.MsgSpendSignBytes(to,coins)
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
			signature, _, err := kb.Sign(msg.From, msg.Password, bz)
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

			Txbytes = stypes.StdtxSpend{
				To:     to,
				Amount: coins,
				Signature: append(Txbytes.Signature, signature),
			}

			data,err:=cdc.MarshalBinary(Txbytes)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Error occurred while marshal binary",
					},
				})
				return
			}

			result := hex.EncodeToString(data)

			w.Write([]byte(result))
			return

			}

		data,err:=hex.DecodeString(msg.Spend)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while hex decode string failed",
				},
			})
			return
		}

		cdc.UnmarshalBinary(data,&Txbytes)


		bz, err := stypes.MsgSpendSignBytes(Txbytes.To, Txbytes.Amount)

		signature, _, err := kb.Sign(msg.From, msg.Password, bz)
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

	Txbytes.Signature= append(Txbytes.Signature, signature)
		output := stypes.NewStdtxSpend(Txbytes.To, Txbytes.Amount, Txbytes.Signature)

		data,err =cdc.MarshalBinary(output)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while marshal binary",
				},
			})
			return
		}

		result := hex.EncodeToString(data)

		w.Write([]byte(result))
		return

	}
}
