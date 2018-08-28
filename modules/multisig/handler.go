package multisig

import (
	"github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) types.Handler {
	return func(ctx types.Context, msg types.Msg) types.Result {
		switch msg := msg.(type) {
		case MsgCreateMultiSigAddress:
			return handleMsgCreateMultiSigAddress(ctx, k, msg)
		case MsgFundMultiSig:
			return handleMsgFundMultiSig(ctx, k, msg)
		case MsgSendFromMultiSig:
			return handleMsgSendFromMultiSig(ctx, k, msg)
		default:
			errMsg := "Unrecognized gov msg type"
			return types.ErrUnknownRequest(errMsg).Result()

		}
	}
}

func handleMsgCreateMultiSigAddress(ctx types.Context, keeper Keeper, msg MsgCreateMultiSigAddress) types.Result {
	address, err := keeper.CreateMultiSigAddress(ctx, msg)
	if err != nil {
		return err.Result()
	}
	tag := types.NewTags("multisig adddress", []byte(address.String()))
	return types.Result{
		Data: address,
		Tags: tag,
	}
}

func handleMsgFundMultiSig(ctx types.Context, keeper Keeper, msg MsgFundMultiSig) types.Result {
	address, err := keeper.FundMultiSig(ctx, msg)
	if err != nil {
		return err.Result()
	}
	tag := types.NewTags("multisig adddress", []byte(address.String()))
	return types.Result{
		Data: address,
		Tags: tag,
	}
}
func handleMsgSendFromMultiSig(ctx types.Context, keeper Keeper, msg MsgSendFromMultiSig) types.Result {
	address, err := keeper.SendFromMultiSig(ctx, msg)
	if err != nil {
		return err.Result()
	}
	tag := types.NewTags("multisig adddress", []byte(msg.Txbytes.MultiSigAddress.String())).AppendTag("reciver address", []byte(address.String()))
	return types.Result{
		Data: address,
		Tags: tag,
	}
}
