package config

import (
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/spf13/viper"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/strings"
	clientflags "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server/config"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	aspecttypes "github.com/artela-network/aspect-core/types"
)

const (
	DefaultAPIEnable       = false
	DefaultGRPCEnable      = false
	DefaultGRPCWebEnable   = false
	DefaultJSONRPCEnable   = false
	DefaultRosettaEnable   = false
	DefaultTelemetryEnable = false

	// DefaultGRPCAddress is the default address the gRPC server binds to.
	DefaultGRPCAddress = "0.0.0.0:9900"

	// DefaultJSONRPCAddress is the default address the JSON-RPC server binds to.
	DefaultJSONRPCAddress = "127.0.0.1:8545"

	// DefaultJSONRPCWsAddress is the default address the JSON-RPC WebSocket server binds to.
	DefaultJSONRPCWsAddress = "127.0.0.1:8546"

	// DefaultJsonRPCMetricsAddress is the default address the JSON-RPC Metrics server binds to.
	DefaultJSONRPCMetricsAddress = "127.0.0.1:6065"

	// DefaultEVMTracer is the default vm.Tracer type
	DefaultEVMTracer = ""

	// DefaultFixRevertGasRefundHeight is the default height at which to overwrite gas refund
	DefaultFixRevertGasRefundHeight = 0

	DefaultMaxTxGasWanted = 0

	DefaultGasCap uint64 = 25000000

	DefaultFilterCap int32 = 200

	DefaultFeeHistoryCap int32 = 100

	DefaultLogsCap int32 = 10000

	DefaultBlockRangeCap int32 = 10000

	DefaultEVMTimeout = 5 * time.Second

	// default 1.0 eth
	DefaultTxFeeCap float64 = 1.0

	DefaultHTTPTimeout = 30 * time.Second

	DefaultHTTPIdleTimeout = 120 * time.Second

	// DefaultAllowUnprotectedTxs value is false
	DefaultAllowUnprotectedTxs = false

	// DefaultMaxOpenConnections represents the amount of open connections (unlimited = 0)
	DefaultMaxOpenConnections = 0
)

var evmTracers = []string{"json", "markdown", "struct", "access_list"}

// Config defines the server's top level configuration. It includes the default app config
// from the SDK as well as the EVM configuration to enable the JSON-RPC APIs.
type Config struct {
	config.Config

	EVM     EVMConfig     `mapstructure:"evm"`
	JSONRPC JSONRPCConfig `mapstructure:"json-rpc"`
	TLS     TLSConfig     `mapstructure:"tls"`
	Aspect  AspectConfig  `mapstructure:"aspect"`
}

// EVMConfig defines the application configuration values for the EVM.
type EVMConfig struct {
	// Tracer defines vm.Tracer type that the EVM will use if the node is run in
	// trace mode. Default: 'json'.
	Tracer string `mapstructure:"tracer"`
	// MaxTxGasWanted defines the gas wanted for each eth txs returned in ante handler in check txs mode.
	MaxTxGasWanted uint64 `mapstructure:"max-txs-gas-wanted"`
}

// AspectConfig defines the application configuration values for Aspect.
type AspectConfig struct {
	// ApplyPoolSize defines capacity of aspect runtime instance pool for applying txs
	ApplyPoolSize int32
	// QueryPoolSize defines capacity of aspect runtime instance pool for querying txs
	QueryPoolSize int32
}

// JSONRPCConfig defines configuration for the EVM RPC server.
type JSONRPCConfig struct {
	// API defines a list of JSON-RPC namespaces that should be enabled
	API []string `mapstructure:"api"`
	// Address defines the HTTP server to listen on
	Address string `mapstructure:"address"`
	// WsAddress defines the WebSocket server to listen on
	WsAddress string `mapstructure:"ws-address"`
	// GasCap is the global gas cap for eth-call variants.
	GasCap uint64 `mapstructure:"gas-cap"`
	// EVMTimeout is the global timeout for eth-call.
	EVMTimeout time.Duration `mapstructure:"evm-timeout"`
	// TxFeeCap is the global txs-fee cap for send txs
	TxFeeCap float64 `mapstructure:"txfee-cap"`
	// FilterCap is the global cap for total number of filters that can be created.
	FilterCap int32 `mapstructure:"filter-cap"`
	// FeeHistoryCap is the global cap for total number of blocks that can be fetched
	FeeHistoryCap int32 `mapstructure:"feehistory-cap"`
	// Enable defines if the EVM RPC server should be enabled.
	Enable bool `mapstructure:"enable"`
	// LogsCap defines the max number of results can be returned from single `eth_getLogs` query.
	LogsCap int32 `mapstructure:"logs-cap"`
	// BlockRangeCap defines the max block range allowed for `eth_getLogs` query.
	BlockRangeCap int32 `mapstructure:"block-range-cap"`
	// HTTPTimeout is the read/write timeout of http json-rpc server.
	HTTPTimeout time.Duration `mapstructure:"http-timeout"`
	// HTTPIdleTimeout is the idle timeout of http json-rpc server.
	HTTPIdleTimeout time.Duration `mapstructure:"http-idle-timeout"`
	// AllowUnprotectedTxs restricts unprotected (non EIP155 signed) transactions to be submitted via
	// the node's RPC when global parameter is disabled.
	AllowUnprotectedTxs bool `mapstructure:"allow-unprotected-txs"`
	// MaxOpenConnections sets the maximum number of simultaneous connections
	// for the server listener.
	MaxOpenConnections int `mapstructure:"max-open-connections"`
	// EnableIndexer defines if enable the custom indexer service.
	EnableIndexer bool `mapstructure:"enable-indexer"`
	// MetricsAddress defines the metrics server to listen on
	MetricsAddress string `mapstructure:"metrics-address"`
	// FixRevertGasRefundHeight defines the upgrade height for fix of revert gas refund logic when txs reverted
	FixRevertGasRefundHeight int64 `mapstructure:"fix-revert-gas-refund-height"`
}

// TLSConfig defines the certificate and matching private key for the server.
type TLSConfig struct {
	// CertificatePath the file path for the certificate .pem file
	CertificatePath string `mapstructure:"certificate-path"`
	// KeyPath the file path for the key .pem file
	KeyPath string `mapstructure:"key-path"`
}

// AppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func AppConfig(denom string) (string, interface{}) {
	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := config.DefaultConfig()

	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In artela, we set the min gas prices to 0.
	if denom != "" {
		srvCfg.MinGasPrices = "0" + denom
	}

	customAppConfig := Config{
		Config:  *srvCfg,
		EVM:     *DefaultEVMConfig(),
		JSONRPC: *DefaultJSONRPCConfig(),
		TLS:     *DefaultTLSConfig(),
	}

	customAppTemplate := config.DefaultConfigTemplate + DefaultConfigTemplate

	return customAppTemplate, customAppConfig
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	return &Config{
		Config:  *DefaultServerConfig(),
		EVM:     *DefaultEVMConfig(),
		JSONRPC: *DefaultJSONRPCConfig(),
		TLS:     *DefaultTLSConfig(),
		Aspect:  *DefaultAspectConfig(),
	}
}

// DefaultConfig returns server's default configuration.
func DefaultServerConfig() *config.Config {
	return &config.Config{
		BaseConfig: config.BaseConfig{
			MinGasPrices:        "",
			InterBlockCache:     true,
			Pruning:             pruningtypes.PruningOptionCustom,
			PruningKeepRecent:   "2",
			PruningInterval:     "10",
			MinRetainBlocks:     0,
			IndexEvents:         make([]string, 0),
			IAVLCacheSize:       781250,
			IAVLDisableFastNode: false,
			IAVLLazyLoading:     false,
			AppDBBackend:        "",
		},
		Telemetry: telemetry.Config{
			Enabled:      false,
			GlobalLabels: [][]string{},
		},
		API: config.APIConfig{
			Enable:             false,
			Swagger:            false,
			Address:            config.DefaultAPIAddress,
			MaxOpenConnections: 1000,
			RPCReadTimeout:     10,
			RPCMaxBodyBytes:    1000000,
		},
		GRPC: config.GRPCConfig{
			Enable:         true,
			Address:        config.DefaultGRPCAddress,
			MaxRecvMsgSize: config.DefaultGRPCMaxRecvMsgSize,
			MaxSendMsgSize: config.DefaultGRPCMaxSendMsgSize,
		},
		Rosetta: config.RosettaConfig{
			Enable:              false,
			Address:             ":8080",
			Blockchain:          "app",
			Network:             "network",
			Retries:             3,
			Offline:             false,
			EnableFeeSuggestion: false,
			GasToSuggest:        clientflags.DefaultGasLimit,
			DenomToSuggest:      "uatom",
		},
		GRPCWeb: config.GRPCWebConfig{
			Enable:  true,
			Address: config.DefaultGRPCWebAddress,
		},
		StateSync: config.StateSyncConfig{
			SnapshotInterval:   2000, // creating snapshot every 2000 blocks
			SnapshotKeepRecent: 5,
		},
		Store: config.StoreConfig{
			Streamers: []string{},
		},
		Streamers: config.StreamersConfig{
			File: config.FileStreamerConfig{
				Keys:            []string{"*"},
				WriteDir:        "",
				OutputMetadata:  true,
				StopNodeOnError: true,
				// NOTICE: The default config doesn't protect the streamer data integrity
				// in face of system crash.
				Fsync: false,
			},
		},
		Mempool: config.MempoolConfig{
			MaxTxs: 5_000,
		},
	}
}

// DefaultEVMConfig returns the default EVM configuration
func DefaultEVMConfig() *EVMConfig {
	return &EVMConfig{
		Tracer:         DefaultEVMTracer,
		MaxTxGasWanted: DefaultMaxTxGasWanted,
	}
}

// Validate returns an error if the tracer type is invalid.
func (c EVMConfig) Validate() error {
	if c.Tracer != "" && !strings.StringInSlice(c.Tracer, evmTracers) {
		return fmt.Errorf("invalid tracer type %s, available types: %v", c.Tracer, evmTracers)
	}

	return nil
}

// DefaultAspectConfig returns the default Aspect configuration
func DefaultAspectConfig() *AspectConfig {
	return &AspectConfig{
		ApplyPoolSize: aspecttypes.DefaultAspectPoolSize,
		QueryPoolSize: aspecttypes.DefaultAspectPoolSize,
	}
}

// Validate returns an error if the tracer type is invalid.
func (a AspectConfig) Validate() error {
	if a.ApplyPoolSize < 0 {
		return errors.New("aspect apply-pool-size cannot be negative")
	}

	if a.QueryPoolSize < 0 {
		return errors.New("aspect query-pool-size cannot be negative")
	}

	return nil
}

// GetDefaultAPINamespaces returns the default list of JSON-RPC namespaces that should be enabled
func GetDefaultAPINamespaces() []string {
	return []string{"eth", "net", "web3"}
}

// GetAPINamespaces returns the all the available JSON-RPC API namespaces.
func GetAPINamespaces() []string {
	return []string{"web3", "eth", "personal", "net", "txpool", "debug", "miner"}
}

// DefaultJSONRPCConfig returns an EVM config with the JSON-RPC API enabled by default
func DefaultJSONRPCConfig() *JSONRPCConfig {
	return &JSONRPCConfig{
		Enable:                   true,
		API:                      GetDefaultAPINamespaces(),
		Address:                  DefaultJSONRPCAddress,
		WsAddress:                DefaultJSONRPCWsAddress,
		GasCap:                   DefaultGasCap,
		EVMTimeout:               DefaultEVMTimeout,
		TxFeeCap:                 DefaultTxFeeCap,
		FilterCap:                DefaultFilterCap,
		FeeHistoryCap:            DefaultFeeHistoryCap,
		BlockRangeCap:            DefaultBlockRangeCap,
		LogsCap:                  DefaultLogsCap,
		HTTPTimeout:              DefaultHTTPTimeout,
		HTTPIdleTimeout:          DefaultHTTPIdleTimeout,
		AllowUnprotectedTxs:      DefaultAllowUnprotectedTxs,
		MaxOpenConnections:       DefaultMaxOpenConnections,
		EnableIndexer:            false,
		MetricsAddress:           DefaultJSONRPCMetricsAddress,
		FixRevertGasRefundHeight: DefaultFixRevertGasRefundHeight,
	}
}

// Validate returns an error if the JSON-RPC configuration fields are invalid.
func (c JSONRPCConfig) Validate() error {
	if c.Enable && len(c.API) == 0 {
		return errors.New("cannot enable JSON-RPC without defining any API namespace")
	}

	if c.FilterCap < 0 {
		return errors.New("JSON-RPC filter-cap cannot be negative")
	}

	if c.FeeHistoryCap <= 0 {
		return errors.New("JSON-RPC feehistory-cap cannot be negative or 0")
	}

	if c.TxFeeCap < 0 {
		return errors.New("JSON-RPC txs fee cap cannot be negative")
	}

	if c.EVMTimeout < 0 {
		return errors.New("JSON-RPC EVM timeout duration cannot be negative")
	}

	if c.LogsCap < 0 {
		return errors.New("JSON-RPC logs cap cannot be negative")
	}

	if c.BlockRangeCap < 0 {
		return errors.New("JSON-RPC block range cap cannot be negative")
	}

	if c.HTTPTimeout < 0 {
		return errors.New("JSON-RPC HTTP timeout duration cannot be negative")
	}

	if c.HTTPIdleTimeout < 0 {
		return errors.New("JSON-RPC HTTP idle timeout duration cannot be negative")
	}

	// check for duplicates
	seenAPIs := make(map[string]bool)
	for _, api := range c.API {
		if seenAPIs[api] {
			return fmt.Errorf("repeated API namespace '%s'", api)
		}

		seenAPIs[api] = true
	}

	return nil
}

// DefaultTLSConfig returns the default TLS configuration
func DefaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		CertificatePath: "",
		KeyPath:         "",
	}
}

// Validate returns an error if the TLS certificate and key file extensions are invalid.
func (c TLSConfig) Validate() error {
	certExt := path.Ext(c.CertificatePath)

	if c.CertificatePath != "" && certExt != ".pem" {
		return fmt.Errorf("invalid extension %s for certificate path %s, expected '.pem'", certExt, c.CertificatePath)
	}

	keyExt := path.Ext(c.KeyPath)

	if c.KeyPath != "" && keyExt != ".pem" {
		return fmt.Errorf("invalid extension %s for key path %s, expected '.pem'", keyExt, c.KeyPath)
	}

	return nil
}

// GetConfig returns a fully parsed Config object.
func GetConfig(v *viper.Viper) (Config, error) {
	cfg, err := config.GetConfig(v)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Config: cfg,
		EVM: EVMConfig{
			Tracer:         v.GetString("evm.tracer"),
			MaxTxGasWanted: v.GetUint64("evm.max-txs-gas-wanted"),
		},
		JSONRPC: JSONRPCConfig{
			Enable:                   v.GetBool("json-rpc.enable"),
			API:                      v.GetStringSlice("json-rpc.api"),
			Address:                  v.GetString("json-rpc.address"),
			WsAddress:                v.GetString("json-rpc.ws-address"),
			GasCap:                   v.GetUint64("json-rpc.gas-cap"),
			FilterCap:                v.GetInt32("json-rpc.filter-cap"),
			FeeHistoryCap:            v.GetInt32("json-rpc.feehistory-cap"),
			TxFeeCap:                 v.GetFloat64("json-rpc.txfee-cap"),
			EVMTimeout:               v.GetDuration("json-rpc.evm-timeout"),
			LogsCap:                  v.GetInt32("json-rpc.logs-cap"),
			BlockRangeCap:            v.GetInt32("json-rpc.block-range-cap"),
			HTTPTimeout:              v.GetDuration("json-rpc.http-timeout"),
			HTTPIdleTimeout:          v.GetDuration("json-rpc.http-idle-timeout"),
			MaxOpenConnections:       v.GetInt("json-rpc.max-open-connections"),
			EnableIndexer:            v.GetBool("json-rpc.enable-indexer"),
			MetricsAddress:           v.GetString("json-rpc.metrics-address"),
			FixRevertGasRefundHeight: v.GetInt64("json-rpc.fix-revert-gas-refund-height"),
			AllowUnprotectedTxs:      v.GetBool("json-rpc.allow-unprotected-txs"),
		},
		TLS: TLSConfig{
			CertificatePath: v.GetString("tls.certificate-path"),
			KeyPath:         v.GetString("tls.key-path"),
		},
		Aspect: AspectConfig{
			ApplyPoolSize: v.GetInt32("aspect.apply-pool-size"),
			QueryPoolSize: v.GetInt32("aspect.query-pool-size"),
		},
	}, nil
}

// ParseConfig retrieves the default environment configuration for the
// application.
func ParseConfig(v *viper.Viper) (*Config, error) {
	conf := DefaultConfig()
	err := v.Unmarshal(conf)

	return conf, err
}

// ValidateBasic returns an error any of the application configuration fields are invalid
func (c Config) ValidateBasic() error {
	if err := c.EVM.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid evm config value: %s", err.Error())
	}

	if err := c.JSONRPC.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid json-rpc config value: %s", err.Error())
	}

	if err := c.TLS.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid tls config value: %s", err.Error())
	}

	if err := c.Aspect.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid aspect config value: %s", err.Error())
	}

	return c.Config.ValidateBasic()
}
