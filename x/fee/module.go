package fee

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

	"github.com/artela-network/artela/x/fee/client/cli"
	"github.com/artela-network/artela/x/fee/keeper"
	"github.com/artela-network/artela/x/fee/types"
)

// TODO mark ConsensusVersion defines the current x/fee module consensus version.
const ConsensusVersion = 4

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ===============================================================
//          		      AppModuleBasic
// ===============================================================

// AppModuleBasic defines the basic application module used by the fee market module.
type AppModuleBasic struct{}

// Name returns the fee market module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec performs a no-op as the fee market module doesn't support amino.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// ConsensusVersion returns the consensus states-breaking version for the module.
func (AppModuleBasic) ConsensusVersion() uint64 {
	return ConsensusVersion
}

// DefaultGenesis returns default genesis states as raw bytes for the fee market
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis is the validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genesisState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genesisState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis states: %w", types.ModuleName, err)
	}

	return genesisState.Validate()
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(c client.Context, serveMux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(c)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns the root txs command for the fee market module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns no root query command for the fee market module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterInterfaces registers interfaces and implementations of the fee market module.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// ===============================================================
//          		        AppModule
// ===============================================================

// AppModule implements an application module for the fee market module.
type AppModule struct {
	AppModuleBasic

	keeper *keeper.Keeper

	// legacySubspace is used solely for migration of x/params managed parameters
	legacySubspace types.Subspace
}

// NewAppModule creates a new AppModule object
func NewAppModule(k *keeper.Keeper, ss types.Subspace) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		legacySubspace: ss,
	}
}

// Name returns the fee market module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants interface for registering invariants. Performs a no-op
// as the fee market module doesn't expose invariants.
func (am AppModule) RegisterInvariants(_ cosmos.InvariantRegistry) {}

// RegisterServices registers the GRPC query service and migrator service to respond to the
// module-specific GRPC queries and handle the upgrade store migration for the module.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
	types.RegisterMsgServer(cfg.MsgServer(), am.keeper)

	// TODO mark
	// m := keeper.NewMigrator(am.keeper, am.legacySubspace)
	// if err := cfg.RegisterMigration(types.ModuleName, 3, m.Migrate3to4); err != nil {
	//	panic(err)
	// }
}

// InitGenesis performs genesis initialization for the fee market module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx cosmos.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState

	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis states as raw bytes for the fee market
// module.
func (am AppModule) ExportGenesis(ctx cosmos.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// RegisterStoreDecoder registers a decoder for fee market module's types
func (am AppModule) RegisterStoreDecoder(_ cosmos.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent { //nolint
	return nil
}

// GenerateGenesisState creates a randomized GenState of the fee market module.
func (AppModule) GenerateGenesisState(_ *module.SimulationState) {
}

// WeightedOperations returns the all the fee market module operations with their respective weights.
func (am AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}

// BeginBlock returns the begin block for the fee market module.
func (am AppModule) BeginBlock(ctx cosmos.Context, req abci.RequestBeginBlock) {
	BeginBlock(ctx, am.keeper, req)
}

// EndBlock returns the end blocker for the fee market module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx cosmos.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlock(ctx, am.keeper, req)
	return []abci.ValidatorUpdate{}
}
