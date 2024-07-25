package rpc

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/eth/gasprice"

	"github.com/artela-network/artela/ethereum/server/config"
)

const (
	// clientIdentifier is the identifier string for the client.
	clientIdentifier = "artela"

	// gpoDefault is the default gpo starting point.
	gpoDefault = 1000000000
)

// DefaultConfig returns the default JSON-RPC config.
func DefaultConfig() *Config {
	gpoConfig := ethconfig.FullNodeGPO
	gpoConfig.Default = big.NewInt(gpoDefault)
	return &Config{
		GPO:           &gpoConfig,
		RPCGasCap:     ethconfig.Defaults.RPCGasCap,
		RPCTxFeeCap:   ethconfig.Defaults.RPCTxFeeCap,
		RPCEVMTimeout: ethconfig.Defaults.RPCEVMTimeout,
	}
}

// Config represents the configurable parameters for Polaris.
type Config struct {
	// AppCfg preserve the server config
	AppCfg *config.Config

	// Gas Price Oracle config.
	GPO *gasprice.Config

	// RPCGasCap is the global gas cap for eth-call variants.
	RPCGasCap uint64 `toml:""`

	// RPCEVMTimeout is the global timeout for eth-call.
	RPCEVMTimeout time.Duration `toml:""`

	// RPCTxFeeCap is the global transaction fee(price * gaslimit) cap for
	// send-transaction variants. The unit is ether.
	RPCTxFeeCap float64 `toml:""`
}

// LoadConfigFromFilePath reads in a Polaris config file from the fileystem.
func LoadConfigFromFilePath(filename string) (*Config, error) {
	var config Config

	// Read the TOML file
	bytes, err := os.ReadFile(filename) //#nosec: G304 // required.
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filename, err)
	}

	// Unmarshal the TOML data into a struct
	if err = toml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("error parsing TOML data: %w", err)
	}

	return &config, nil
}
