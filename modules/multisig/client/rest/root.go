package rest

import (
	"sentinel/modules/multisig"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
)

var msgWireCdc = wire.NewCodec()

func init() {
	multisig.RegisterWire(msgWireCdc)
}

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec) {
	r.HandleFunc("/initiate", initiateWalletHandlerFn(cdc, cliCtx)).Methods("POST")   // to initiate txns to generate addr
	r.HandleFunc("/create", createWalletAddressHandleFn(cdc, cliCtx)).Methods("POST") //to get wallet address
	r.HandleFunc("/initiate_txn", initiateTxnHandlerFn(cdc, cliCtx)).Methods("POST")  // to sign the transaction in offchain
	r.HandleFunc("/send", sendHandleFn(cdc, cliCtx)).Methods("POST")                  // to send money to wallet
	r.HandleFunc("/transfer", transferHandleFn(cdc, cliCtx)).Methods("POST")

}
