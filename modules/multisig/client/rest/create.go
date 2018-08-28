package rest

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sentinel/modules/multisig"
	sdk "sentinel/modules/multisig/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	context2 "github.com/cosmos/cosmos-sdk/x/auth/client/context"
)

type MultiSigAddrCreate struct {
	Txbytes       string `json:"txbytes"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	ChainId       string `json:"chain_id"`
	AccountNumber int64  `json:"account_number"`
	Gas           int64  `json:"gas"`
}

func multisignatureCreateAddressFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg MultiSigAddrCreate
		var err error
		var Txbytes multisig.Stdtx

		// Decoinding the Request
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSigAddrCreateResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while decoding the request body",
				},
			})
			return
		}
		if msg.Txbytes == "" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSigAddrCreateResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while fetching tx bytes",
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
					"Error occurred while  decode string failed",
				},
			})
			return
		}

		err=cdc.UnmarshalBinary(data, &Txbytes)
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSignatureResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while unmarshl string failed",
				},
			})
			return
		}
		cliCtx = cliCtx.WithFromAddressName(msg.Name)
		cliCtx = cliCtx.WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

		address, err := cliCtx.GetFromAddress()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSigAddrCreateResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while retrive the address",
				},
			})
			return
		}

		account, err := cliCtx.GetAccount(address)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSigAddrCreateResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while retrive the account",
				},
			})
			return
		}

		sequence := account.GetSequence()
		txcontext := context2.TxContext{
			Codec:         cdc,
			ChainID:       msg.ChainId,
			Sequence:      sequence,
			AccountNumber: msg.AccountNumber,
			Gas:           msg.Gas,
		}
		message := multisig.NewMsgCreateMultiSigAddress(Txbytes, address)

		txbytes, err := txcontext.BuildAndSign(msg.Name, msg.Password, []types.Msg{message})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSigAddrCreateResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while retrive the txbytes",
				},
			})
			return
		}

		res, err := cliCtx.BroadcastTx(txbytes)
		if err != nil {
			panic(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSigAddrCreateResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while txbytes broadcast failed",
				},
			})
			return
		}

		output, err := wire.MarshalJSONIndent(cdc, res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(sdk.MultiSigAddrCreateResponse{
				Success: false,
				Error: sdk.Error{
					1,
					"Error occurred while retrive data",
				},
			})
			return
		}
		w.Write(output)

	}
}
