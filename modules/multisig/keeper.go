package multisig

import ("github.com/cosmos/cosmos-sdk/types"
"github.com/cosmos/cosmos-sdk/x/auth"
"github.com/cosmos/cosmos-sdk/x/bank"
"github.com/cosmos/cosmos-sdk/wire"
mtypes "sentinel/modules/multisig/types"
	"strconv"
	"encoding/json"
	"github.com/tendermint/tendermint/crypto"
)
type Keeper struct {

	multiStoreKey types.StoreKey
	coinKeeper   bank.Keeper
	cdc          *wire.Codec
	codespace types.CodespaceType
	account   auth.AccountMapper
}

func NewKeeper(cdc *wire.Codec, key types.StoreKey, ck bank.Keeper, am auth.AccountMapper, codespace types.CodespaceType) Keeper {
	return Keeper{
		multiStoreKey: key,
		cdc:          cdc,
		coinKeeper:   ck,
		codespace:    codespace,
		account:      am,
	}

}

func(k Keeper)CreateMultiSigAddress(ctx types.Context, msg MsgCreateMultiSigAddress)(types.AccAddress,types.Error){

	var err error

	store:=ctx.KVStore(k.multiStoreKey)

	sequence,err:=k.account.GetSequence(ctx,msg.Address)
	if err!=nil{
			return nil,mtypes.ErrInvalidSequence("Invalid sequence")
	}

	addressbytes:=[]byte(msg.Address.String()+""+strconv.Itoa(int(sequence)))
	addressgen:=crypto.Sha256(addressbytes)[:20]
	address:=types.AccAddress(addressgen)
	account :=k.account.NewAccountWithAddress(ctx,address)
	k.account.SetAccount(ctx,account)

	for i:=0; i<int(msg.Txbytes.TotalKeys);i++{

		pubkey,err:=types.GetAccPubKeyBech32(msg.Txbytes.Pubkey[i])
		if err!=nil{
			return nil,mtypes.ErrInvalidPubKey("convert of bech32 pubkey failed")
		}

		bz,err:=mtypes.CreateSignBytes(msg.Txbytes.MinKeys,msg.Txbytes.Order,msg.Txbytes.TotalKeys)
		if err!=nil{
			return nil,mtypes.ErrMarshal("convert of marshal failed")
		}
		if(!pubkey.VerifyBytes(bz,msg.Txbytes.Signature[i])){
			return nil,mtypes.ErrSignatureVerfication("signature verification failed")
		}
	}

	bz, err := json.Marshal(msg.Txbytes)
	if err!=nil{
		return nil, mtypes.ErrMarshal("marshal bytes failed")
	}
	store.Set([]byte(address),bz)
	return address,nil
}

func (k Keeper)FundMultiSig(ctx types.Context, msg MsgFundMultiSig)(types.AccAddress,types.Error){
	var txbytes Stdtx
	var count int64

	store:=ctx.KVStore(k.multiStoreKey)
	data:=store.Get([]byte(msg.To))
	if data==nil{
		return nil,mtypes.ErrDataFromKVStore("Failed to get data from KVStore")

	}

	err:=json.Unmarshal(data,&txbytes)
	if err!=nil{
		return nil,mtypes.ErrUnMarshal("Unmarshal of byte failed")
	}

	pubkey,err:=k.account.GetPubKey(ctx,msg.Address)
	if err!=nil{
		return nil, mtypes.ErrInvalidPubKey("Retrive of pubkey is failed from account")
	}

	pubkeystr,_:=types.Bech32ifyAccPub(pubkey)
	for _, value := range txbytes.Pubkey{
		if value==pubkeystr{
			count++
		}
	}
	if count!=1{
		return nil,mtypes.ErrInvalidMultiSigAccount("account not associated with this multisig to fund")
	}

	k.coinKeeper.SubtractCoins(ctx,msg.Address,msg.Amount)
	k.coinKeeper.AddCoins(ctx,msg.To,msg.Amount)
	return msg.To,nil
}

func (k Keeper) SendFromMultiSig(ctx types.Context,msg MsgSendFromMultiSig) (types.AccAddress,types.Error){
	var txbytes Stdtx
	var count int64

	store:=ctx.KVStore(k.multiStoreKey)

	data:=store.Get([]byte(msg.MultiSigAddress))
	if data==nil{
		return nil,mtypes.ErrDataFromKVStore("Failed to get data from KVStore")
	}

	err:=json.Unmarshal(data,&txbytes)
	if err!=nil{
		return nil,mtypes.ErrUnMarshal("Unmarshal of byte failed")
	}

	for i:=0; i<int(txbytes.TotalKeys);i++{

		pubkey,err:=types.GetAccPubKeyBech32(txbytes.Pubkey[i])
		if err!=nil{
			return nil,mtypes.ErrInvalidPubKey("convert of bech32 pubkey failed")
		}

		bz,err:=mtypes.MsgSpendSignBytes(msg.Txbytes.To,msg.Txbytes.Amount)
		if err!=nil{
			return nil,mtypes.ErrMarshal("convert of marshal failed")
		}
		if(txbytes.Order==true){
			if(!pubkey.VerifyBytes(bz,msg.Txbytes.Signature[i])){
				return nil,mtypes.ErrSignatureVerfication("signature verification failed")
			}
		}
		if(txbytes.Order==false){
		for _,value := range msg.Txbytes.Signature{
			if(pubkey.VerifyBytes(bz,value)){
				count++;
			}
		}
		}
	}

	if count >= int64(txbytes.MinKeys){
		k.coinKeeper.AddCoins(ctx,msg.Txbytes.To,msg.Txbytes.Amount)
		k.coinKeeper.SubtractCoins(ctx,msg.MultiSigAddress,msg.Txbytes.Amount)
	}else {
		return nil,mtypes.ErrSigners("required minimun no of signers")
	}
	return msg.Txbytes.To,nil
}