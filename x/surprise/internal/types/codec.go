package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateBrandedToken{}, "surprise/CreateBrandedToken", nil)
	cdc.RegisterConcrete(MsgTransferBrandedTokenOwnership{}, "surprise/TransferBrandedTokenOwnership", nil)
	cdc.RegisterConcrete(MsgBurnBrandedToken{}, "surprise/BurnBrandedToken", nil)
	cdc.RegisterConcrete(MsgMintBrandedToken{}, "surprise/MintBrandedToken", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
