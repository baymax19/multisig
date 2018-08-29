package rest

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sentinel/modules/multisig"
	"sentinel/modules/multisig/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	context2 "github.com/cosmos/cosmos-sdk/x/auth/client/context"
)

type MsgSendFromMultiSig struct {
	Txbytes       string `json:"txbytes"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	ChainId       string `json:"chain_id"`
	AccountNumber int64  `json:"account_number"`
	Gas           int64  `json:"gas"`
}

/**
* @api {post} /transfer To send tokens from Multisig wallet.
* @apiName Transfer tokens from Multisig wallet
* @apiGroup MultisigWallet
* @apiParam {String} to To Address.
* @apiParam {String} name Name of Account.
* @apiParam {String} password Password for account.
* @apiParam {String} chain_id Chain Id.
* @apiParam {Number} account_number Account number.
* @apiParam {Number} gas Gas value.
* @apiParam {String} amount amount to send.
* @apiError AccountAlreadyExists AccountName is  already exists
* @apiErrorExample AccountAlreadyExists-Response:
* {
*   Account with name XXXXX... already exists.
* }
* @apiSuccessExample Response:
*{
*  "check_tx": {
*    "log": "Msg 0: ",
*    "gasWanted": "21000",
*    "gasUsed": "1209"
*  },
*  "deliver_tx": {
*    "data": "IGpXqASE+6AvgVVRO3NyNtCzqH4=",
*    "log": "Msg 0: ",
*    "gasWanted": "21000",
*    "gasUsed": "6670",
*    "tags": [
*      {
*        "key": "bXVsdGlzaWcgYWRkZHJlc3M=",
*        "value": "Y29zbW9zYWNjYWRkcjF5cDQ5MDJxeXNuYTZxdHVwMjRnbmt1bWp4bWd0ODJyN2g1a3V2Yw=="
*      }
*    ]
*  },
*  "hash": "CC78A0E5445A2EE945308F3A599EF96BD529A9AF",
*  "height": "14863"
*}
 */

func transferHandleFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg MsgSendFromMultiSig
		var err error
		var Txbytes types.StdtxSpend

		// Decoinding the Request
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MsgSendFromMultiSig{
				Success: false,
				Error: types.Error{
					1,
					"Error occurred while decoding the request body",
				},
			})
			return
		}

		if msg.Txbytes == "" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MsgSendFromMultiSig{
				Success: false,
				Error: types.Error{
					1,
					"Error occurred while tx bytes nil",
				},
			})
			return

		}

		data, err := base64.StdEncoding.DecodeString(msg.Txbytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MsgSendFromMultiSig{
				Success: false,
				Error: types.Error{
					1,
					"Error occurred while deocding tx bytes failed",
				},
			})
			return
		}

		err = cdc.UnmarshalBinary(data, &Txbytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MsgSendFromMultiSig{
				Success: false,
				Error: types.Error{
					1,
					"Error occurred while unmarshal data",
				},
			})
			return

		}

		cliCtx = cliCtx.WithFromAddressName(msg.Name)
		cliCtx = cliCtx.WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

		address, err := cliCtx.GetFromAddress()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MultiSigAddrCreateResponse{
				Success: false,
				Error: types.Error{
					1,
					"Error occurred while retrive the address",
				},
			})
			return
		}
		account, err := cliCtx.GetAccount(address)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MultiSigAddrCreateResponse{
				Success: false,
				Error: types.Error{
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

		message := multisig.NewMsgSendFromMultiSig(Txbytes, address)

		txbytes, err := txcontext.BuildAndSign(msg.Name, msg.Password, []sdk.Msg{message})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MultiSigAddrCreateResponse{
				Success: false,
				Error: types.Error{
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
			json.NewEncoder(w).Encode(types.MultiSigAddrCreateResponse{
				Success: false,
				Error: types.Error{
					1,
					"Error occurred while txbytes broadcast failed",
				},
			})
			return
		}

		output, err := wire.MarshalJSONIndent(cdc, res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MultiSigAddrCreateResponse{
				Success: false,
				Error: types.Error{
					1,
					"Error occurred while retrive data",
				},
			})
			return
		}
		w.Write(output)
	}
}
