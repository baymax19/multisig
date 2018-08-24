package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"net/http"
	"sentinel/modules/multisig/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"encoding/json"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	context2 "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sentinel/modules/multisig"
	"encoding/hex"
)

type MsgSendFromMultiSig struct {
	Txbytes string `json:"txbytes"`
	MultiSigAddress string  `json:"multi_sig_address"`
	Name string `json:"name"`
	Password string `json:"password"`
	ChainId string `json:"chain_id"`
	AccountNumber int64 `json:"account_number"`
	Gas int64 `json:"gas"`

}


func multisignatureSendFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
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

		data,err:=hex.DecodeString(msg.Txbytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MsgSendFromMultiSig{
				Success: false,
				Error: types.Error{
					1,
					"Error occurred while hex decode string failed",
				},
			})
			return
		}

		cdc.UnmarshalBinary(data,&Txbytes)


		cliCtx=cliCtx.WithFromAddressName(msg.Name)
		cliCtx=cliCtx.WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

		address,err:=cliCtx.GetFromAddress()
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
		account,err := cliCtx.GetAccount(address)
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

		sequence:=account.GetSequence()

		txcontext:=context2.TxContext{
			Codec:cdc,
			ChainID:msg.ChainId,
			Sequence:sequence,
			AccountNumber:msg.AccountNumber,
			Gas:msg.Gas,
		}

		multisigaddress,err:=sdk.AccAddressFromBech32(msg.MultiSigAddress)
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MultiSigAddrCreateResponse{
				Success: false,
				Error: types.Error{
					1,
					"Error occurred while getting bech32 address multisg string",
				},
			})
			return

		}
		message:=multisig.NewMsgSendFromMultiSig(multisigaddress,Txbytes,address)

		txbytes,err:=txcontext.BuildAndSign(msg.Name,msg.Password,[]sdk.Msg{message})
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

		res,err:=cliCtx.BroadcastTx(txbytes)
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