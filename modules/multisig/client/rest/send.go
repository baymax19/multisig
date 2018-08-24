package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"net/http"
)

type MsgSendFromMultiSig struct {
	To string `json:"to"`

}


func multisignatureFundFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg MsgFungMultiSigAddr