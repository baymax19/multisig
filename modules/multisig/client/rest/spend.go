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
	"github.com/cosmos/cosmos-sdk/types"
	"encoding/base64"
)

type MsgSpendFromMultiSig struct {
	Spend    string `json:"spend"`
	To       string     `json:"to"`
	MultiSigAddress string `json:"multi_sig_address"`
	From     string     `json:"from"`
	Amount   string     `json:"amount"`
	Password string     `json:"password"`
	TxNumber *int64    `json:"tx_number"`
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
					2,
					"Error occurred while Fetching the keybase",
				},
			})
			return
		}

		//check the  existing txbytes
		if (msg.Spend=="") {

			if(msg.Amount==""){
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						3,
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
						3,
						"Enter to account",
					},
				})
				return
			}

			if(msg.MultiSigAddress==""){
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						3,
						"Enter multisig Address",
					},
				})
				return
			}

			if(*msg.TxNumber<1){
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						3,
						"Error tx number should be starts with 1",
					},
				})
				return

			}

			if(msg.To==msg.MultiSigAddress){
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Error multisig address and to address will not be same",
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

			maddress,err:=types.AccAddressFromBech32(msg.MultiSigAddress)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"multisig address from bech64 string is failed ",
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

			bz, err := stypes.MsgSpendSignBytes(to,coins,maddress)
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
				MultiSigAddress:maddress,
				Amount: coins,
				Signature: append(Txbytes.Signature, signature),
				TxNumber:*msg.TxNumber,
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

			result := base64.StdEncoding.EncodeToString(data)

			w.Write([]byte(result))
			return

			}
		if(msg.Amount!="" || msg.MultiSigAddress!=""|| msg.TxNumber!=nil){
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					" amount , multisig address or tx_number data is not allowed to enter when you have spen tx bytes",
				},
			})
			return

		}
		data,err:=base64.StdEncoding.DecodeString(msg.Spend)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while  decode tx  string failed",
				},
			})
			return
		}

		err=cdc.UnmarshalBinary(data,&Txbytes)
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while unmarshal failed",
				},
			})
			return
		}


		bz, err := stypes.MsgSpendSignBytes(Txbytes.To, Txbytes.Amount,Txbytes.MultiSigAddress)
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while sign bytes is failed",
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
					"Error occurred while creating sign bytes ",
				},
			})
			return
		}

		for _,value :=range Txbytes.Signature{
			if string(value)==string(signature){
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Error occurred signature is already exist ",
					},
				})
				return

			}
		}

	    Txbytes.Signature= append(Txbytes.Signature, signature)
		output := stypes.NewStdtxSpend(Txbytes.To, Txbytes.MultiSigAddress,Txbytes.Amount, Txbytes.Signature,Txbytes.TxNumber)

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

		result := base64.StdEncoding.EncodeToString(data)
		w.Write([]byte(result))
		return

	}
}
