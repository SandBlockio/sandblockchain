package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const CreateBrandedTokenConst = "CreateBrandedToken"
const MsgTransferBrandedTokenOwnershipConst = "TransferBrandedTokenOwnership"
const MsgMintBrandedTokenConst = "MintBrandedToken"
const MsgBurnBrandedTokenConst = "BurnBrandedToken"

// MsgCreateBrandedToken
type MsgCreateBrandedToken struct {
	Name          string         `json:"name"`
	InitialSupply sdk.Int        `json:"supply"`
	FromAddress   sdk.AccAddress `json:"from_address"`
}

var _ sdk.Msg = &MsgCreateBrandedToken{}

func NewMsgCreateBrandedToken(name string, supply sdk.Int, creator sdk.AccAddress) MsgCreateBrandedToken {
	return MsgCreateBrandedToken{
		Name:          name,
		InitialSupply: supply,
		FromAddress:   creator,
	}
}

func (msg MsgCreateBrandedToken) Route() string { return RouterKey }
func (msg MsgCreateBrandedToken) Type() string  { return CreateBrandedTokenConst }
func (msg MsgCreateBrandedToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.FromAddress)}
}
func (msg MsgCreateBrandedToken) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
func (msg MsgCreateBrandedToken) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "from_address can't be empty")
	}
	if len(msg.Name) <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name can't be empty")
	}
	if msg.InitialSupply.LT(sdk.NewInt(0)) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "supply can't be less or equal than 0")
	}
	return nil
}

// MsgTransferBrandedTokenOwnership
type MsgTransferBrandedTokenOwnership struct {
	Name        string         `json:"name"`
	FromAddress sdk.AccAddress `json:"from_address"`
	NewOwner    sdk.AccAddress `json:"new_owner"`
}

var _ sdk.Msg = &MsgTransferBrandedTokenOwnership{}

func NewMsgTransferBrandedTokenOwnership(name string, previousOwner sdk.AccAddress, newOwner sdk.AccAddress) MsgTransferBrandedTokenOwnership {
	return MsgTransferBrandedTokenOwnership{
		Name:        name,
		FromAddress: previousOwner,
		NewOwner:    newOwner,
	}
}

func (msg MsgTransferBrandedTokenOwnership) Route() string { return RouterKey }
func (msg MsgTransferBrandedTokenOwnership) Type() string {
	return MsgTransferBrandedTokenOwnershipConst
}
func (msg MsgTransferBrandedTokenOwnership) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "previous_owner can't be empty")
	}
	if msg.NewOwner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "new_owner can't be empty")
	}
	if len(msg.Name) <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name can't be empty")
	}
	return nil
}
func (msg MsgTransferBrandedTokenOwnership) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.FromAddress)}
}
func (msg MsgTransferBrandedTokenOwnership) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgMintBrandedToken
type MsgMintBrandedToken struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	Name        string         `json:"name"`
	Amount      sdk.Int        `json:"amount"`
}

var _ sdk.Msg = &MsgMintBrandedToken{}

func NewMsgMintBrandedToken(owner sdk.AccAddress, name string, amount sdk.Int) MsgMintBrandedToken {
	return MsgMintBrandedToken{
		FromAddress: owner,
		Name:        name,
		Amount:      amount,
	}
}

func (msg MsgMintBrandedToken) Route() string { return RouterKey }
func (msg MsgMintBrandedToken) Type() string  { return MsgMintBrandedTokenConst }
func (msg MsgMintBrandedToken) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner can't be empty")
	}
	if len(msg.Name) <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name can't be empty")
	}
	if msg.Amount.LTE(sdk.NewInt(0)) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "amount can't be empty")
	}
	return nil
}
func (msg MsgMintBrandedToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgMintBrandedToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgBurnBrandedToken
type MsgBurnBrandedToken struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	Name        string         `json:"name"`
	Amount      sdk.Int        `json:"amount"`
}

var _ sdk.Msg = &MsgBurnBrandedToken{}

func NewMsgBurnBrandedToken(owner sdk.AccAddress, name string, amount sdk.Int) MsgBurnBrandedToken {
	return MsgBurnBrandedToken{
		FromAddress: owner,
		Name:        name,
		Amount:      amount,
	}
}

func (msg MsgBurnBrandedToken) Route() string { return RouterKey }
func (msg MsgBurnBrandedToken) Type() string  { return MsgBurnBrandedTokenConst }
func (msg MsgBurnBrandedToken) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner can't be empty")
	}
	if len(msg.Name) <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name can't be empty")
	}
	if msg.Amount.LTE(sdk.NewInt(0)) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "amount can't be empty")
	}
	return nil
}
func (msg MsgBurnBrandedToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgBurnBrandedToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}
