package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path"
	"testing"
)

func TestCLIConfig(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)
	node := fmt.Sprintf("%s:%s", f.RPCAddr, f.Port)

	// Set available configuration options
	f.CLIConfig("broadcast-mode", "block")
	f.CLIConfig("node", node)
	f.CLIConfig("output", "text")
	f.CLIConfig("trust-node", "true")
	f.CLIConfig("chain-id", f.ChainID)
	f.CLIConfig("trace", "false")
	f.CLIConfig("indent", "true")

	config, err := ioutil.ReadFile(path.Join(f.GaiacliHome, "config", "config.toml"))
	require.NoError(t, err)
	expectedConfig := fmt.Sprintf(`broadcast-mode = "block"
chain-id = "%s"
indent = true
node = "%s"
output = "text"
trace = false
trust-node = true
`, f.ChainID, node)
	require.Equal(t, expectedConfig, string(config))

	f.Cleanup()
}