package rest

import (
	"net/http"

	"encoding/json"
	sdk "sentinel/modules/multisig/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	stypes "sentinel/modules/multisig/types"
	"encoding/hex"
	"sentinel/modules/multisig"
)

type MultiSignature struct {
	Txbytes   string  `json:"txbytes,omitempty"`
	MinKeys   uint8  `json:"min_keys,omitempty"`
	TotalKeys uint8  `json:"total_keys,omitempty"`
	Order     bool   `json:"order,omitempty"`
	Name      string `json:"name"`
	Password  string `json:"password"`
}

func multisignatureCreateSignHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg MultiSignature
		var kb keys.Keybase
		var err error
		var Txbytes multisig.Stdtx

		// Decoding the Request
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while decoding the request body",
				},
			})
			return
		}



		// Fetching the keybase
		kb, err = ckeys.GetKeyBase()
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

		//Check the txbytes exist or not
		if   msg.Txbytes == "" {

			if(msg.MinKeys<1){
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Reuiree at least one key",
					},
				})
				return
			}

			if(msg.TotalKeys<2 || msg.MinKeys>msg.TotalKeys) {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"total keys grater than min keys",
					},
				})
				return

			}


			bz, err := stypes.CreateSignBytes(msg.MinKeys, msg.Order, msg.TotalKeys)
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

			signature, pubkey, err := kb.Sign(msg.Name, msg.Password, bz)
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
			output := multisig.Stdtx{
				MinKeys:   msg.MinKeys,
				TotalKeys: msg.TotalKeys,
				Order:     msg.Order,
				Pubkey:    append(Txbytes.Pubkey, pubkeystr),
				Counter:   0,
				Signature: append(Txbytes.Signature, signature),
			}

			data,err:=cdc.MarshalBinary(output)
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

		data,err:=hex.DecodeString(msg.Txbytes)
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

		//check for required  no.of keys
		if uint8(Txbytes.Counter) < Txbytes.TotalKeys-1 {

			bz, err := stypes.CreateSignBytes(Txbytes.MinKeys, Txbytes.Order, Txbytes.TotalKeys)
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

			signature, pubkey, _ := kb.Sign(msg.Name, msg.Password, bz)
			pubkeystr, _ := types.Bech32ifyAccPub(pubkey)

			for _, value := range Txbytes.Pubkey {
				if value == pubkeystr {
					json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
						Success: false,
						Error: sdk.Error{
							1,
							"Error occurred the given signature and publickey already exist",
						},
					})
					return
				}
			}

			Txbytes.Pubkey = append(Txbytes.Pubkey, pubkeystr)
			Txbytes.Signature= append(Txbytes.Signature, signature)
			output := stypes.NewStdtx(Txbytes.Order, Txbytes.TotalKeys, Txbytes.MinKeys, Txbytes.Pubkey, Txbytes.Counter+1,Txbytes.Signature)

			data,err:=cdc.MarshalBinary(output)
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

		json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
			Success: false,
			Error: sdk.Error{
				1,
				"Error occurred No of Publickeys and signatures  are more ",
			},
		})

		return

	}
}
