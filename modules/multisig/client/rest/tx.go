package rest

import (
	"net/http"

	"encoding/json"
	sdk "sentinel/modules/multisig/types"

	"encoding/base64"
	"sentinel/modules/multisig"
	stypes "sentinel/modules/multisig/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type MultiSignature struct {
	Txbytes   string `json:"txbytes,omitempty"`
	MinKeys   *uint8 `json:"min_keys,omitempty"`
	TotalKeys *uint8 `json:"total_keys,omitempty"`
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
		if msg.Txbytes == "" {

			if msg.MinKeys == nil || msg.TotalKeys == nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Insufficient of data",
					},
				})
				return

			}

			if *msg.MinKeys < 1 {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Requir at least one key",
					},
				})
				return
			}

			if *msg.TotalKeys < 2 || *msg.MinKeys > *msg.TotalKeys || *msg.TotalKeys > 10 {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"min keys are grater than total keys or  max keys greaterthan upper bond 10	",
					},
				})
				return

			}

			bz, err := stypes.CreateSignBytes(*msg.MinKeys, msg.Order, *msg.TotalKeys)
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

			pubkeystr, err := types.Bech32ifyAccPub(pubkey)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Error while converting the pubkey to string",
					},
				})
				return
			}

			pubkey_encode := base64.StdEncoding.EncodeToString([]byte(pubkeystr))

			output := multisig.Stdtx{
				MinKeys:   *msg.MinKeys,
				TotalKeys: *msg.TotalKeys,
				Order:     msg.Order,
				Pubkey:    append(Txbytes.Pubkey, pubkey_encode),
				Counter:   1,
				Signature: append(Txbytes.Signature, signature),
			}

			data, err := cdc.MarshalBinary(output)
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

			data1 := base64.StdEncoding.EncodeToString(data)
			w.Write([]byte(data1))
			return
		}

		if msg.MinKeys != nil || msg.TotalKeys != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred Not allowed of minkeys , order ,totalkeys",
				},
			})
			return
		}

		data, err := base64.StdEncoding.DecodeString(msg.Txbytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while decode txbytes failed",
				},
			})
			return
		}

		err = cdc.UnmarshalBinary(data, &Txbytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while unmarshal txbytes",
				},
			})
			return

		}

		//check for required  no.of keys
		if uint8(len(Txbytes.Pubkey)) < Txbytes.TotalKeys {

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
			pubkey_encode := base64.StdEncoding.EncodeToString([]byte(pubkeystr))

			for _, value := range Txbytes.Pubkey {
				if value == pubkey_encode {
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

			Txbytes.Pubkey = append(Txbytes.Pubkey, pubkey_encode)
			Txbytes.Signature = append(Txbytes.Signature, signature)
			output := stypes.NewStdtx(Txbytes.Order, Txbytes.TotalKeys, Txbytes.MinKeys, Txbytes.Pubkey, 1, Txbytes.Signature)

			data, err := cdc.MarshalBinary(output)
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

			encode_data := base64.StdEncoding.EncodeToString(data)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
					Success: false,
					Error: sdk.Error{
						1,
						"Error occurred conversion of base64",
					},
				})
				return
			}
			w.Write([]byte(encode_data))
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
