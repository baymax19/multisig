package types
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)
type Error struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}

type MultiSignatureResponse struct {
	Success bool   `json:"success"`
	Error   Error  `json:"error,omitempty"`
	Data    []byte `json:"data,omitempty"`
}

type MsgSpendFromMultiSigResponse struct {
	Success bool   `json:"success"`
	Error   Error  `json:"error,omitempty"`
	Data    []byte `json:"data,omitempty"`
}


type MultiSigAddrCreateResponse struct {
	Success bool   `json:"success"`
	Error   Error  `json:"error,omitempty"`
	Data    []byte `json:"data,omitempty"`
}

type MsgFungMultiSigAddrResponse struct {
	Success bool   `json:"success"`
	Error   Error  `json:"error,omitempty"`
	Data    []byte `json:"data,omitempty"`
}

const (
	DefaultCodespace sdk.CodespaceType = 19

	CodeInvalidSequence sdk.CodeType = 191
	CodeMarshalInterface          sdk.CodeType=192
	CodeInvalidAddres           sdk.CodeType=193
	CodeCreateMultiSig          sdk.CodeType=194
	CodeDataFromKVStore         sdk.CodeType=195
	CodeUnMarshal               sdk.CodeType=196
	CodeInvalidPubKey sdk.CodeType=197
	CodeInvalidMultiSigAccount   sdk.CodeType=198
	CodeUnknownRequest  sdk.CodeType = sdk.CodeUnknownRequest
)


func ErrInvalidSequence(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidSequence,msg)
}

func ErrMarshal(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeMarshalInterface,msg)
}
func ErrInvalidAddres(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidAddres,msg)
}
func ErrCreateMultiSig(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeCreateMultiSig,msg)
}
func ErrDataFromKVStore(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeDataFromKVStore,msg)
}
func ErrUnMarshal(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeUnMarshal,msg)
}

func ErrInvalidPubKey(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidPubKey,msg)
}
func ErrInvalidMultiSigAccount(msg string) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidMultiSigAccount,msg)
}
// -------------------------
// Helpers

// nolint: unparam
func newError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	return sdk.NewError(codespace, code, msg)
}
