package tests

import (
	"encoding/json"
	"fmt"
	_ "github.com/cosmos/cosmos-sdk/client"
	clientkeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sandblockio/sandblockchain/app"
	"github.com/stretchr/testify/require"
	tmtypes "github.com/tendermint/tendermint/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	_ "time"
)

const (
	denom  = "sbc"
	denomStake = "stake" // TODO: change the config to use the same coin (sbc) for coin and stake
	keyFoo = "foo"
	keyBar = "bar"
	DefaultKeyPass = "12345678"

	brandedToken1 = "tetine"
	brandedToken2 = "zeuzo"
)

var (
	totalCoins = sdk.NewCoins(
		sdk.NewCoin(denom, sdk.TokensFromConsensusPower(2000000)),
		sdk.NewCoin(denomStake, sdk.TokensFromConsensusPower(2000000)),
	)

	startCoins = sdk.NewCoins(
		sdk.NewCoin(denom, sdk.TokensFromConsensusPower(2000000)),
		sdk.NewCoin(denomStake, sdk.TokensFromConsensusPower(2000000)),
	)
)

func init() {
	// Set the addresses prefix
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	config.Seal()
}

type Fixtures struct {
	BuildDir      string
	RootDir       string
	GaiadBinary   string
	GaiacliBinary string
	ChainID       string
	RPCAddr       string
	Port          string
	GaiadHome     string
	GaiacliHome   string
	P2PAddr       string
	T             *testing.T
}

// NewFixtures creates a new instance of Fixtures with many vars set
func NewFixtures(t *testing.T) *Fixtures {
	tmpDir, err := ioutil.TempDir("", "sandblock_integration_"+t.Name()+"_")
	require.NoError(t, err)

	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)

	p2pAddr, _, err := server.FreeTCPAddr()
	require.NoError(t, err)

	buildDir := os.Getenv("BUILDDIR")
	if buildDir == "" {
		buildDir, err = filepath.Abs("../build/")
		require.NoError(t, err)
	}

	return &Fixtures{
		T:             t,
		BuildDir:      buildDir,
		RootDir:       tmpDir,
		GaiadBinary:   filepath.Join(buildDir, "sbd"),
		GaiacliBinary: filepath.Join(buildDir, "sbcli"),
		GaiadHome:     filepath.Join(tmpDir, ".sbd"),
		GaiacliHome:   filepath.Join(tmpDir, ".sbcli"),
		RPCAddr:       servAddr,
		P2PAddr:       p2pAddr,
		Port:          port,
	}
}

// Cleanup is meant to be run at the end of a test to clean up an remaining test state
func (f *Fixtures) Cleanup(dirs ...string) {
	clean := append(dirs, f.RootDir)
	for _, d := range clean {
		require.NoError(f.T, os.RemoveAll(d))
	}
}

// Flags returns the flags necessary for making most CLI calls
func (f *Fixtures) Flags() string {
	return fmt.Sprintf("--home=%s --node=%s", f.GaiacliHome, f.RPCAddr)
}

// GenesisFile returns the path of the genesis file
func (f Fixtures) GenesisFile() string {
	return filepath.Join(f.GaiadHome, "config", "genesis.json")
}

// GenesisFile returns the application's genesis state
func (f Fixtures) GenesisState() simapp.GenesisState {
	cdc := codec.New()
	genDoc, err := tmtypes.GenesisDocFromFile(f.GenesisFile())
	require.NoError(f.T, err)

	var appState simapp.GenesisState
	require.NoError(f.T, cdc.UnmarshalJSON(genDoc.AppState, &appState))
	return appState
}

// InitFixtures is called at the beginning of a test  and initializes a chain
// with 1 validator.
func InitFixtures(t *testing.T) (f *Fixtures) {
	// Init the local fixtures
	f = NewFixtures(t)

	// Reset the previous state
	f.UnsafeResetAll()

	// Ensure the local CLI config is containing JSON as output mode
	f.CLIConfig("output", "json")

	// Reset the previous keys
	f.KeysDelete(keyFoo)
	f.KeysDelete(keyBar)
	f.KeysAdd(keyFoo)
	f.KeysAdd(keyBar)

	// Init the local test node with the keyFoo
	f.GDInit(keyFoo)

	// Ensure the local CLI config is all set
	f.CLIConfig("chain-id", f.ChainID)
	f.CLIConfig("broadcast-mode", "block")
	f.CLIConfig("trust-node", "true")

	// Add a genesis account with start coins
	f.AddGenesisAccount(f.KeyAddress(keyFoo), startCoins)

	f.GenTx(keyFoo)
	f.CollectGenTxs()

	return
}

/// ##############
//	KEYS MANAGEMENT
//	##############

// KeysDelete is gaiacli keys delete
func (f *Fixtures) KeysDelete(name string, flags ...string) {
	cmd := fmt.Sprintf("%s keys delete --keyring-backend test --home=%s %s", f.GaiacliBinary, f.GaiacliHome, name)
	executeWrite(f.T, addFlags(cmd, append(append(flags, "-y"), "-f")))
}

// KeysAdd is gaiacli keys add
func (f *Fixtures) KeysAdd(name string, flags ...string) {
	cmd := fmt.Sprintf("%s keys add --keyring-backend test --home=%s %s", f.GaiacliBinary, f.GaiacliHome, name)
	executeWriteCheckErr(f.T, addFlags(cmd, flags), DefaultKeyPass)
}

// KeysAddRecover prepares gaiacli keys add --recover
func (f *Fixtures) KeysAddRecover(name, mnemonic string, flags ...string) (exitSuccess bool, stdout, stderr string) {
	cmd := fmt.Sprintf("%s keys add --keyring-backend test --home=%s --recover %s", f.GaiacliBinary, f.GaiacliHome, name)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), DefaultKeyPass, mnemonic)
}

// KeysAddRecoverHDPath prepares gaiacli keys add --recover --account --index
func (f *Fixtures) KeysAddRecoverHDPath(name, mnemonic string, account uint32, index uint32, flags ...string) {
	cmd := fmt.Sprintf("%s keys add --keyring-backend test --home=%s --recover %s --account %d --index %d", f.GaiacliBinary, f.GaiacliHome, name, account, index)
	executeWriteCheckErr(f.T, addFlags(cmd, flags), DefaultKeyPass, mnemonic)
}

// KeysShow is gaiacli keys show
func (f *Fixtures) KeysShow(name string, flags ...string) keys.KeyOutput {
	cmd := fmt.Sprintf("%s keys show --keyring-backend test --home=%s %s", f.GaiacliBinary, f.GaiacliHome, name)
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var ko keys.KeyOutput
	err := clientkeys.UnmarshalJSON([]byte(out), &ko)
	require.NoError(f.T, err)
	return ko
}

// KeyAddress returns the SDK account address from the key
func (f *Fixtures) KeyAddress(name string) sdk.AccAddress {
	ko := f.KeysShow(name)
	accAddr, err := sdk.AccAddressFromBech32(ko.Address)
	require.NoError(f.T, err)
	return accAddr
}

/// ##############
//	GLOBAL SBD COMMANDS
//	##############

func (f *Fixtures) UnsafeResetAll(flags ...string) {
	cmd := fmt.Sprintf("%s --home=%s unsafe-reset-all", f.GaiadBinary, f.GaiadHome)
	executeWrite(f.T, addFlags(cmd, flags))
	err := os.RemoveAll(filepath.Join(f.GaiadHome, "config", "gentx"))
	require.NoError(f.T, err)
}

// GDInit is gaiad init
// NOTE: GDInit sets the ChainID for the Fixtures instance
func (f *Fixtures) GDInit(moniker string, flags ...string) {
	cmd := fmt.Sprintf("%s init -o --home=%s %s", f.GaiadBinary, f.GaiadHome, moniker)
	_, stderr := tests.ExecuteT(f.T, addFlags(cmd, flags), DefaultKeyPass)

	var chainID string
	var initRes map[string]json.RawMessage

	err := json.Unmarshal([]byte(stderr), &initRes)
	require.NoError(f.T, err)

	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(f.T, err)

	f.ChainID = chainID
}

// AddGenesisAccount is gaiad add-genesis-account
func (f *Fixtures) AddGenesisAccount(address sdk.AccAddress, coins sdk.Coins, flags ...string) {
	cmd := fmt.Sprintf("%s add-genesis-account %s %s --home=%s", f.GaiadBinary, address, coins, f.GaiadHome)
	executeWriteCheckErr(f.T, addFlags(cmd, flags))
}

// GenTx is gaiad gentx
func (f *Fixtures) GenTx(name string, flags ...string) {
	cmd := fmt.Sprintf("%s gentx --keyring-backend test --name=%s --home=%s --home-client=%s", f.GaiadBinary, name, f.GaiadHome, f.GaiacliHome)
	executeWriteCheckErr(f.T, addFlags(cmd, flags), DefaultKeyPass)
}

// CollectGenTxs is gaiad collect-gentxs
func (f *Fixtures) CollectGenTxs(flags ...string) {
	cmd := fmt.Sprintf("%s collect-gentxs --home=%s", f.GaiadBinary, f.GaiadHome)
	executeWriteCheckErr(f.T, addFlags(cmd, flags), DefaultKeyPass)
}

// GDStart runs gaiad start with the appropriate flags and returns a process
func (f *Fixtures) GDStart(flags ...string) *tests.Process {
	cmd := fmt.Sprintf("%s start --home=%s --rpc.laddr=%v --p2p.laddr=%v", f.GaiadBinary, f.GaiadHome, f.RPCAddr, f.P2PAddr)
	proc := tests.GoExecuteTWithStdout(f.T, addFlags(cmd, flags))
	tests.WaitForTMStart(f.Port)
	tests.WaitForNextNBlocksTM(1, f.Port)
	return proc
}

// GDTendermint returns the results of gaiad tendermint [query]
func (f *Fixtures) GDTendermint(query string) string {
	cmd := fmt.Sprintf("%s tendermint %s --home=%s", f.GaiadBinary, query, f.GaiadHome)
	success, stdout, stderr := executeWriteRetStdStreams(f.T, cmd)
	require.Empty(f.T, stderr)
	require.True(f.T, success)
	return strings.TrimSpace(stdout)
}

// ValidateGenesis runs gaiad validate-genesis
func (f *Fixtures) ValidateGenesis() {
	cmd := fmt.Sprintf("%s validate-genesis --home=%s", f.GaiadBinary, f.GaiadHome)
	executeWriteCheckErr(f.T, cmd)
}

/// ##############
//	GLOBAL SBCLI COMMANDS
//	##############

func (f *Fixtures) CLIConfig(key, value string, flags ...string) {
	cmd := fmt.Sprintf("%s config --home=%s %s %s", f.GaiacliBinary, f.GaiacliHome, key, value)
	executeWriteCheckErr(f.T, addFlags(cmd, flags))
}
