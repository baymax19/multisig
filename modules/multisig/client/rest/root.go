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
	r.HandleFunc("/multisig", multisignatureCreateSignHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/multisig/spend", multisignatureSpendHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/multisig/create", multisignatureCreateAddressFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/multisig/fund", multisignatureFundFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/multisig/transfer", multisignatureSendFn(cdc, cliCtx)).Methods("POST")



}
