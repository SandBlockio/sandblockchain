package surprise

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sandblockio/sandblockchain/x/surprise/internal/types"

	"github.com/gosimple/slug"
)

// NewHandler creates an sdk.Handler for all the surprise type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateBrandedToken:
			return handleMsgCreateBrandedToken(ctx, k, msg)

		case types.MsgTransferBrandedTokenOwnership:
			return handleMsgTransferBrandedTokenOwnership(ctx, k, msg)

		case types.MsgMintBrandedToken:
			return handleMsgMintBrandedToken(ctx, k, msg)

		case types.MsgBurnBrandedToken:
			return handleMsgBurnBrandedToken(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreateBrandedToken(ctx sdk.Context, k Keeper, msg types.MsgCreateBrandedToken) (*sdk.Result, error) {
	// Construct a slug from the name
	tokenSlug := slug.Make(msg.Name)

	// Ensure the branded token does not exists
	if k.HasBrandedToken(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Branded Token with that name already exists")
	}

	// Create the branded token
	newBrandedToken, _ := k.GetBrandedToken(ctx, tokenSlug)
	newBrandedToken.Coin = sdk.NewCoin(msg.Name, msg.InitialSupply)
	newBrandedToken.Owner = msg.FromAddress
	k.SetBrandedToken(ctx, tokenSlug, newBrandedToken)

	// Add the coin to the coin keeper
	_, err = k.CoinKeeper.AddCoins(ctx, newBrandedToken.GetOwner(), sdk.NewCoins(newBrandedToken.Coin))
	if err != nil {
		// Delete the persisted coin and return
		k.DeleteBrandedToken(ctx, tokenSlug)
		return nil, sdkerrors.Wrap(sdkerrors.ErrPanic, "Failure when setting the coins on the coinKeeper")
	}

	// Emit the log-events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.InitialSupply.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgTransferBrandedTokenOwnership(ctx sdk.Context, k Keeper, msg types.MsgTransferBrandedTokenOwnership) (*sdk.Result, error) {
	// Construct a slug from the name
	tokenSlug := slug.Make(msg.Name)

	// Ensure the branded token exists
	if !k.HasBrandedToken(ctx, tokenSlug) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The given branded token does not exists")
	}

	// Fetch the entity from keeper
	brandedToken, err := k.GetBrandedToken(ctx, tokenSlug)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to fetch the branded token from kvstore")
	}

	// Ensure the initiator is the owner
	if !brandedToken.GetOwner().Equals(msg.FromAddress) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "You are not the owner of that BrandedToken")
	}

	// Finally change the owner and update the entity
	brandedToken.Owner = msg.NewOwner
	k.SetBrandedToken(ctx, tokenSlug, brandedToken)

	// Emit the log-events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgMintBrandedToken(ctx sdk.Context, k Keeper, msg types.MsgMintBrandedToken) (*sdk.Result, error) {
	// Construct a slug from the name
	tokenSlug := slug.Make(msg.Name)

	// Ensure the branded token exists
	if !k.HasBrandedToken(ctx, tokenSlug) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The given branded token does not exists")
	}

	// Fetch the entity from keeper
	brandedToken, err := k.GetBrandedToken(ctx, tokenSlug)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to fetch the branded token from kvstore")
	}

	// Ensure the initiator is the owner
	if !brandedToken.GetOwner().Equals(msg.FromAddress) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "You are not the owner of that BrandedToken")
	}

	// Update the coin keeper
	_, err = k.CoinKeeper.AddCoins(ctx, msg.FromAddress, sdk.NewCoins(sdk.NewCoin(brandedToken.GetName(), msg.Amount)))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPanic, "Failure when adding the coins on the coinKeeper")
	}

	//  Update and persist the entity
	brandedToken.Amount.Add(msg.Amount)
	k.SetBrandedToken(ctx, tokenSlug, brandedToken)

	// Emit the log-event and return
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgBurnBrandedToken(ctx sdk.Context, k Keeper, msg types.MsgBurnBrandedToken) (*sdk.Result, error) {
	// Construct a slug from the name
	tokenSlug := slug.Make(msg.Name)

	// Ensure the branded token exists
	if !k.HasBrandedToken(ctx, tokenSlug) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The given branded token does not exists")
	}

	// Fetch the entity from keeper
	brandedToken, err := k.GetBrandedToken(ctx, tokenSlug)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to fetch the branded token from kvstore")
	}

	// Ensure the initiator is the owner
	if !brandedToken.GetOwner().Equals(msg.FromAddress) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "You are not the owner of that BrandedToken")
	}

	// Update the coin keeper
	_, err = k.CoinKeeper.SubtractCoins(ctx, msg.FromAddress, sdk.NewCoins(sdk.NewCoin(brandedToken.GetName(), msg.Amount)))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPanic, "Failure when substracting the coins on the coinKeeper")
	}

	//  Update and persist the entity
	brandedToken.Amount.Sub(msg.Amount)
	k.SetBrandedToken(ctx, tokenSlug, brandedToken)

	// Emit the log-event and return
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}