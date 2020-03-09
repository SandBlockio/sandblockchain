package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type BrandedToken struct {
	sdk.Coin
	Owner sdk.AccAddress `json:"owner"`
}

func (token BrandedToken) GetName() string          { return token.Denom }
func (token BrandedToken) GetAmount() sdk.Int       { return token.Amount }
func (token BrandedToken) GetOwner() sdk.AccAddress { return token.Owner }
func (token BrandedToken) SetOwner(owner sdk.AccAddress) BrandedToken {
	token.Owner = owner
	return token
}

func NewBrandedToken() BrandedToken {
	return BrandedToken{}
}

func (token BrandedToken) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Name: %s|Owner: %s|TotalSupply: %d`, token.GetName(), token.GetOwner(), token.GetAmount()))
}