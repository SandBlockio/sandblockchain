package keeper

import (
	"github.com/gosimple/slug"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sandblockio/sandblockchain/x/surprise/internal/types"
)

// NewQuerier creates a new querier for surprise clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGetBrandedToken:
			return queryGetBrandedToken(ctx, path[1:], k)

		case types.QueryListBrandedTokens:
			return queryListBrandedTokens(ctx, k)

		case types.QueryGetTotalSupply:
			return queryTotalSupply(ctx, k)

		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown surprise query endpoint")
		}
	}
}

func queryTotalSupply(ctx sdk.Context, k Keeper) ([]byte, error){
	var supply sdk.Int

	// Acquire the iterator and loop
	iterator := k.GetBrandedTokensIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		token, _ := k.GetBrandedToken(ctx, string(iterator.Key()))
		supply.Add(token.Amount)
	}

	// Convert and return
	res, err := codec.MarshalJSONIndent(k.cdc, supply)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryListBrandedTokens(ctx sdk.Context, k Keeper) ([]byte, error){
	var tokens types.QueryResFetch

	// Acquire the iterator and loop
	iterator := k.GetBrandedTokensIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		tokens = append(tokens, string(iterator.Key()))
	}

	// Convert and return
	res, err := codec.MarshalJSONIndent(k.cdc, tokens)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryGetBrandedToken(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	tokenSlug := slug.Make(path[0])

	// Ensure the branded token exists
	if !k.HasBrandedToken(ctx, tokenSlug) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The branded token does not exist")
	}

	// Fetch the entity
	brandedToken, err := k.GetBrandedToken(ctx, tokenSlug)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Unable to fetch the branded token")
	}

	// Convert and return
	res, err := codec.MarshalJSONIndent(k.cdc, brandedToken)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}