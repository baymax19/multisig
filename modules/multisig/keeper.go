package multisig

import (
	"encoding/json"
	mtypes "sentinel/modules/multisig/types"
	"strconv"

	"encoding/base64"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/crypto"
)

type Keeper struct {
	multiStoreKey types.StoreKey
	coinKeeper    bank.Keeper
	cdc           *wire.Codec
	codespace     types.CodespaceType
	account       auth.AccountMapper
}

func NewKeeper(cdc *wire.Codec, key types.StoreKey, ck bank.Keeper, am auth.AccountMapper, codespace types.CodespaceType) Keeper {
	return Keeper{
		multiStoreKey: key,
		cdc:           cdc,
		coinKeeper:    ck,
		codespace:     codespace,
		account:       am,
	}

}

func (k Keeper) CreateMultiSigAddress(ctx types.Context, msg MsgCreateMultiSigAddress) (types.AccAddress, types.Error) {
	var err error

	store := ctx.KVStore(k.multiStoreKey)

	sequence, err := k.account.GetSequence(ctx, msg.Address)
	if err != nil {
		return nil, mtypes.ErrInvalidSequence("Invalid sequence")
	}

	if len(msg.Txbytes.Signature) != int(msg.Txbytes.TotalKeys) {
		return nil, mtypes.ErrSignatureVerfication("signatures or pubkeys are not equal to  " + strconv.Itoa(int(msg.Txbytes.TotalKeys)) + "total keys")
	}

	for i := 0; i < int(msg.Txbytes.TotalKeys); i++ {

		pubkey_decode, err := base64.StdEncoding.DecodeString(msg.Txbytes.Pubkey[i])
		if err != nil {
			return nil, mtypes.ErrInvalidPubKey("pubkey decoding is failed")
		}

		pubkey, err := types.GetAccPubKeyBech32(string(pubkey_decode))
		if err != nil {
			return nil, mtypes.ErrInvalidPubKey("convert of bech32 pubkey failed")
		}

		bz, err := mtypes.CreateSignBytes(msg.Txbytes.MinKeys, msg.Txbytes.Order, msg.Txbytes.TotalKeys)
		if err != nil {
			return nil, mtypes.ErrMarshal("convert of marshal failed")
		}

		if !pubkey.VerifyBytes(bz, msg.Txbytes.Signature[i]) {
			return nil, mtypes.ErrSignatureVerfication("signature verification failed")
		}

	}

	addressbytes := []byte(msg.Address.String() + "" + strconv.Itoa(int(sequence)))
	addressgen := crypto.Sha256(addressbytes)[:20]
	address := types.AccAddress(addressgen)
	account := k.account.NewAccountWithAddress(ctx, address)
	k.account.SetAccount(ctx, account)

	bz, err := json.Marshal(msg.Txbytes)
	if err != nil {
		return nil, mtypes.ErrMarshal("marshal bytes failed")
	}

	store.Set([]byte(address), bz)
	return address, nil
}

func (k Keeper) FundMultiSig(ctx types.Context, msg MsgFundMultiSig) (types.AccAddress, types.Error) {
	var txbytes Stdtx

	store := ctx.KVStore(k.multiStoreKey)
	data := store.Get([]byte(msg.To))
	if data == nil {
		return nil, mtypes.ErrDataFromKVStore("Failed to get data from KVStore")

	}

	err := json.Unmarshal(data, &txbytes)
	if err != nil {
		return nil, mtypes.ErrUnMarshal("Unmarshal of byte failed")
	}

	_, _, err = k.coinKeeper.SubtractCoins(ctx, msg.Address, msg.Amount)
	if err != nil {
		return nil, mtypes.ErrInsufficientCoins("Insufficient funds from account ")
	}
	_, _, err = k.coinKeeper.AddCoins(ctx, msg.To, msg.Amount)
	if err != nil {
		return nil, mtypes.ErrInsufficientCoins("Insuffiecient funds")
	}
	return msg.To, nil
}

func (k Keeper) SendFromMultiSig(ctx types.Context, msg MsgSendFromMultiSig) (types.AccAddress, types.Error) {
	var txbytes Stdtx
	var count int64
	var index int64

	store := ctx.KVStore(k.multiStoreKey)

	data := store.Get([]byte(msg.Txbytes.MultiSigAddress))
	if data == nil {
		return nil, mtypes.ErrDataFromKVStore("Failed to get data from KVStore")
	}

	err := json.Unmarshal(data, &txbytes)
	if err != nil {
		return nil, mtypes.ErrUnMarshal("Unmarshal of byte failed")
	}

	if len(msg.Txbytes.Signature) < int((txbytes.MinKeys)) {

		return nil, mtypes.ErrSigners("required minimun no of signers ")
	}

	for i := 0; i < len(txbytes.Pubkey); i++ {

		pukey_decode, err := base64.StdEncoding.DecodeString(txbytes.Pubkey[i])
		if err != nil {
			return nil, mtypes.ErrInvalidPubKey("decode of  pubkey failed")
		}

		pubkey, err := types.GetAccPubKeyBech32(string(pukey_decode))
		if err != nil {
			return nil, mtypes.ErrInvalidPubKey("convert of bech32 pubkey failed")
		}

		bz, err := mtypes.MsgSpendSignBytes(msg.Txbytes.To, msg.Txbytes.Amount, msg.Txbytes.MultiSigAddress)
		if err != nil {
			return nil, mtypes.ErrMarshal("convert of marshal failed")
		}

		if txbytes.Order == true && count < int64((txbytes.MinKeys)) {

			if !pubkey.VerifyBytes(bz, msg.Txbytes.Signature[index]) {
				count = 0
				continue
			}
			index++
			count++

		}
		if txbytes.Order == false {

			for _, value := range msg.Txbytes.Signature {
				if pubkey.VerifyBytes(bz, value) {
					count++
				}
			}
		}

	}

	if txbytes.Counter != msg.Txbytes.TxNumber {
		return nil, mtypes.ErrTxNumber("Tx number is invalid got this  " + strconv.Itoa(int(msg.Txbytes.TxNumber)) + " expected this  " + strconv.Itoa(int(txbytes.Counter)))

	}

	if count < int64(txbytes.MinKeys) {
		return nil, mtypes.ErrSigners("required minimun no of signers or  order is not same")
	}

	_, _, err = k.coinKeeper.SubtractCoins(ctx, msg.Txbytes.MultiSigAddress, msg.Txbytes.Amount)
	if err != nil {
		return nil, mtypes.ErrInsufficientCoins("Insufficient funds from multisig wallet")
	}
	_, _, err = k.coinKeeper.AddCoins(ctx, msg.Txbytes.To, msg.Txbytes.Amount)
	if err != nil {
		return nil, mtypes.ErrInsufficientCoins("Insufficient funds")
	}

	txbytes.Counter++

	bz, err := json.Marshal(txbytes)
	if err != nil {
		return nil, mtypes.ErrMarshal("marshal bytes failed")
	}

	store.Set([]byte(msg.Txbytes.MultiSigAddress), bz)
	return msg.Txbytes.To, nil
}
