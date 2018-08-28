package rest

import (
	"encoding/json"
	"net/http"
	"sentinel/modules/multisig/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	context2 "github.com/cosmos/cosmos-sdk/x/auth/client/context"

	"sentinel/modules/multisig"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgFungMultiSigAddr struct {
	To            string `json:"to"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	Chainid       string `json:"chainid"`
	AccountNumber int64  `json:"account_number"`
	Gas           int64  `json:"gas"`
	Amount        string `json:"amount"`
}

func multisignatureFundFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg MsgFungMultiSigAddr
		var err error

		// Decoinding the Request
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MsgFungMultiSigAddrResponse{
				Success: false,
				Error: types.Error{
					1,
					"Error occurred while decoding the request body",
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
			ChainID:       msg.Chainid,
			Sequence:      sequence,
			AccountNumber: msg.AccountNumber,
			Gas:           msg.Gas,
		}

		coins, err := sdk.ParseCoins(msg.Amount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MultiSigAddrCreateResponse{
				Success: false,
				Error: types.Error{
					1,
					"Error occurred while coins parsing is failed",
				},
			})
			return

		}

		to, err := sdk.AccAddressFromBech32(msg.To)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(types.MultiSigAddrCreateResponse{
				Success: false,
				Error: types.Error{
					1,
					"Error occurred while getting bech32 address from string",
				},
			})
			return

		}

		message := multisig.NewMsgFundMultiSig(to, address, coins)

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
