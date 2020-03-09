package surprise

import (
	"github.com/sandblockio/sandblockchain/x/surprise/internal/keeper"
	"github.com/sandblockio/sandblockchain/x/surprise/internal/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	QuerierRoute      = types.QuerierRoute
)

var (
	// functions aliases
	NewKeeper                           = keeper.NewKeeper
	NewQuerier                          = keeper.NewQuerier
	RegisterCodec                       = types.RegisterCodec
	NewGenesisState                     = types.NewGenesisState
	DefaultGenesisState                 = types.DefaultGenesisState
	ValidateGenesis                     = types.ValidateGenesis
	NewMsgCreateBrandedToken            = types.NewMsgCreateBrandedToken
	NewMsgTransferBrandedTokenOwnership = types.NewMsgTransferBrandedTokenOwnership
	NewMsgMintBrandedToken              = types.NewMsgMintBrandedToken
	NewMsgBurnBrandedToken              = types.NewMsgBurnBrandedToken

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper = keeper.Keeper
	GenesisState = types.GenesisState
	Params = types.Params

	MsgCreateBrandedToken = types.MsgCreateBrandedToken
	MsgTransferBrandedTokenOwnership = types.MsgTransferBrandedTokenOwnership
	MsgMintBrandedToken = types.MsgMintBrandedToken
	MsgBurnBrandedToken = types.MsgBurnBrandedToken
)
