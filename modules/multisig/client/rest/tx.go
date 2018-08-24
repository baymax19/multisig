package rest

import (
	"net/http"

	"encoding/json"
	sdk "sentinel/modules/multisig/types"

	"reflect"

	"github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	stypes "sentinel/modules/multisig/types"
)

type MultiSignature struct {
	Txbytes   stypes.Stdtx  `json:"txbytes,omitempty"`
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
		var output stypes.Stdtx

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

		//Check the txbytes exist or not
		if reflect.DeepEqual(stypes.Stdtx{}, msg.Txbytes) {
			bz, err := CreateSignBytes(msg.MinKeys, msg.Order, msg.TotalKeys)
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

			_, pubkey, err := kb.Sign(msg.Name, msg.Password, bz)
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
			output = stypes.Stdtx{
				MinKeys:   msg.MinKeys,
				TotalKeys: msg.TotalKeys,
				Order:     msg.Order,
				Pubkey:    append(output.Pubkey, pubkeystr),
				Counter:   0,
			}

			data, _ := json.Marshal(output)
			w.Write(data)
			return
		}

		//check for required  no.of keys
		if uint8(msg.Txbytes.Counter) < msg.Txbytes.TotalKeys-1 {

			bz, err := CreateSignBytes(msg.Txbytes.MinKeys, msg.Txbytes.Order, msg.Txbytes.TotalKeys)
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

			_, pubkey, _ := kb.Sign(msg.Name, msg.Password, bz)
			pubkeystr, _ := types.Bech32ifyAccPub(pubkey)

			for _, value := range msg.Txbytes.Pubkey {
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

			msg.Txbytes.Pubkey = append(msg.Txbytes.Pubkey, pubkeystr)
			output = stypes.NewStdtx(msg.Txbytes.Order, msg.Txbytes.TotalKeys, msg.Txbytes.MinKeys, msg.Txbytes.Pubkey, msg.Txbytes.Counter+1)

			data, _ := json.Marshal(output)
			w.Write(data)
			return
		}

		json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
			Success: false,
			Error: sdk.Error{
				1,
				"Error occurred No of Publickeys are more ",
			},
		})
		return

	}
}
