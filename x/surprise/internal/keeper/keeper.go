package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sandblockio/sandblockchain/x/surprise/internal/types"
)

// Keeper of the surprise store
type Keeper struct {
	CoinKeeper bank.Keeper
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramspace types.ParamSubspace
}

// NewKeeper creates a surprise keeper
func NewKeeper(coinKeeper bank.Keeper, cdc *codec.Codec, key sdk.StoreKey, paramspace types.ParamSubspace) Keeper {
	keeper := Keeper{
		CoinKeeper: coinKeeper,
		storeKey:   key,
		cdc:        cdc,
		paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetBrandedToken return a branded token by its slug if it exist, return a blank one if not. Should never be exposed publicly
func (k Keeper) GetBrandedToken(ctx sdk.Context, key string) (types.BrandedToken, error) {
	store := ctx.KVStore(k.storeKey)

	// If it does not exists we return an empty one
	if !store.Has([]byte(key)){
		return types.NewBrandedToken(), nil
	}

	token := types.BrandedToken{}
	bz := store.Get([]byte(key))

	// If there is an error we return an empty one with the error
	err := k.cdc.UnmarshalBinaryBare(bz, &token)
	if err != nil {
		return token, err
	}

	// Everything went fine, we return
	return token, nil
}

// SetBrandedToken update a new branded token struct inside the datastore. Should never be exposed publicly
func (k Keeper) SetBrandedToken(ctx sdk.Context, key string, value types.BrandedToken) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(value))
}

// HasBrandedToken return a bool depending the given branded token exists or not
func (k Keeper) HasBrandedToken(ctx sdk.Context, key string) bool {
	return ctx.KVStore(k.storeKey).Has([]byte(key))
}

// DeleteBrandedToken delete the corresponding branded token
func (k Keeper) DeleteBrandedToken(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(key))
}

// GetBrandedTokensIterator return an iterator over all tokens
func (k Keeper) GetBrandedTokensIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}