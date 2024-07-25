package evm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/artela-network/artela/x/evm/client/cli"
	"github.com/artela-network/artela/x/evm/keeper"
	"github.com/artela-network/artela/x/evm/txs"
	"github.com/artela-network/artela/x/evm/txs/support"
	"github.com/artela-network/artela/x/evm/types"
)

// TODO mark ConsensusVersion defines the current x/evm module consensus version.
const ConsensusVersion = 7

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ===============================================================
//          		      AppModuleBasic
// ===============================================================

// AppModuleBasic defines the basic application module used by the evm module.
type AppModuleBasic struct{}

// Name returns the evm module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the module's types with the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	txs.RegisterLegacyAminoCodec(cdc)
}

// ConsensusVersion returns the consensus states-breaking version for the module.
func (AppModuleBasic) ConsensusVersion() uint64 {
	return ConsensusVersion
}

// DefaultGenesis returns default genesis states as raw bytes for the evm
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(support.DefaultGenesisState())
}

// ValidateGenesis is the validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data support.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis states: %w", types.ModuleName, err)
	}

	return data.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the evm module.
func (b AppModuleBasic) RegisterGRPCGatewayRoutes(c client.Context, serveMux *runtime.ServeMux) {
	if err := txs.RegisterQueryHandlerClient(context.Background(), serveMux, txs.NewQueryClient(c)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns the root txs command for the evm module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns no root query command for the evm module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterInterfaces registers interfaces and implementations of the evm module.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	txs.RegisterInterfaces(registry)
}

// ===============================================================
//          		        AppModule
// ===============================================================

// AppModule implements an application module for the evm module.
type AppModule struct {
	AppModuleBasic

	keeper *keeper.Keeper
	ak     types.AccountKeeper

	// legacySubspace is used solely for migration of x/params managed parameters
	legacySubspace types.Subspace
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	txs.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	txs.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	m := keeper.NewMigrator(*am.keeper)

	if err := cfg.RegisterMigration(types.ModuleName, 5, m.Migrate5to6); err != nil {
		panic(err)
	}

	if err := cfg.RegisterMigration(types.ModuleName, 6, m.Migrate6to7); err != nil {
		panic(err)
	}
}

// NewAppModule creates a new AppModule object
func NewAppModule(k *keeper.Keeper, ak types.AccountKeeper, ss types.Subspace) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		ak:             ak,
		legacySubspace: ss,
	}
}

// Name returns the evm module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// QuerierRoute returns the evm module's querier route name.
func (AppModule) QuerierRoute() string { return types.RouterKey }

// InitGenesis performs genesis initialization for the evm module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx cosmos.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState support.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.ak, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis states as raw bytes for the evm
// module.
func (am AppModule) ExportGenesis(ctx cosmos.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper, am.ak)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// GenerateGenesisState creates a randomized GenState of the evm module.
func (AppModule) GenerateGenesisState(_ *module.SimulationState) {
	//  simulation.RandomizedGenState(simState) // TODO mark
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	// return simulation.ProposalMsgs() // TODO mark
	return nil
}

// RegisterStoreDecoder registers a decoder for evm module's types
func (am AppModule) RegisterStoreDecoder(_ cosmos.StoreDecoderRegistry) {
}

// WeightedOperations returns the all the evm module operations with their respective weights.
func (am AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}

// BeginBlock returns the begin block for the evm module.
func (am AppModule) BeginBlock(ctx cosmos.Context, req abci.RequestBeginBlock) {
	BeginBlock(ctx, am.keeper, req)
}

// EndBlock returns the end blocker for the evm module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx cosmos.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlock(ctx, am.keeper, req)
}
