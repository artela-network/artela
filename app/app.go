package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authmodule "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankmodule "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitymodule "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusmodule "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisismodule "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrmodule "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencemodule "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilmodule "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govmodule "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	mintmodule "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsmodule "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingmodule "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingmodule "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgrademodule "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	transfermodule "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	solomachine "github.com/cosmos/ibc-go/v7/modules/light-clients/06-solomachine"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/spf13/cast"

	"github.com/artela-network/artela/app/upgrades/v047rc7"
	"github.com/artela-network/artela/app/upgrades/v048rc8"
	evmmodule "github.com/artela-network/artela/x/evm"
	evmmodulekeeper "github.com/artela-network/artela/x/evm/keeper"
	evmmoduletypes "github.com/artela-network/artela/x/evm/types"
	feemodule "github.com/artela-network/artela/x/fee"
	feemodulekeeper "github.com/artela-network/artela/x/fee/keeper"
	feemoduletypes "github.com/artela-network/artela/x/fee/types"

	// this line is used by starport scaffolding # stargate/app/moduleImport

	"github.com/artela-network/artela/app/ante"
	ethante "github.com/artela-network/artela/app/ante/evm"
	appparams "github.com/artela-network/artela/app/params"
	"github.com/artela-network/artela/app/post"
	"github.com/artela-network/artela/common"
	"github.com/artela-network/artela/docs"
	srvflags "github.com/artela-network/artela/ethereum/server/flags"
	artela "github.com/artela-network/artela/ethereum/types"
	aspecttypes "github.com/artela-network/aspect-core/types"

	// do not remove this, this will register the native evm tracers
	_ "github.com/artela-network/artela-evm/tracers/native"
)

const (
	AccountAddressPrefix = "artela"
	Name                 = "artela"
	ExeName              = "artelad"
)

// this line is used by starport scaffolding # stargate/wasm/app/enabledProposals

func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler
	// this line is used by starport scaffolding # stargate/app/govProposalHandlers

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
		// this line is used by starport scaffolding # stargate/app/govProposalHandler
	)

	return govProposalHandlers
}

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutilmodule.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(getGovProposalHandlers()),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		groupmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibctm.AppModuleBasic{},
		solomachine.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		ica.AppModuleBasic{},
		vesting.AppModuleBasic{},
		consensus.AppModuleBasic{},
		evmmodule.AppModuleBasic{},
		feemodule.AppModuleBasic{},
		// this line is used by starport scaffolding # stargate/app/moduleBasic
	)

	// module account permissions
	maccPerms = map[string][]string{
		authmodule.FeeCollectorName:     nil,
		distrmodule.ModuleName:          nil,
		icatypes.ModuleName:             nil,
		mintmodule.ModuleName:           {authmodule.Minter},
		stakingmodule.BondedPoolName:    {authmodule.Burner, authmodule.Staking},
		stakingmodule.NotBondedPoolName: {authmodule.Burner, authmodule.Staking},
		govmodule.ModuleName:            {authmodule.Burner},
		transfermodule.ModuleName:       {authmodule.Minter, authmodule.Burner},
		// this line is used by starport scaffolding # stargate/app/maccPerms
		evmmoduletypes.ModuleName: {authmodule.Minter, authmodule.Burner},
	}
)

var (
	_ runtime.AppI            = (*Artela)(nil)
	_ servertypes.Application = (*Artela)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+ExeName)

	// manually update the power reduction
	cosmos.DefaultPowerReduction = artela.PowerReduction
}

// Artela extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type Artela struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry
	txConfig          client.TxConfig

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	AuthzKeeper           authzkeeper.Keeper
	BankKeeper            bankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             govkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper        evidencekeeper.Keeper
	TransferKeeper        ibctransferkeeper.Keeper
	ICAHostKeeper         icahostkeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper  capabilitykeeper.ScopedKeeper

	EvmKeeper *evmmodulekeeper.Keeper

	FeeKeeper *feemodulekeeper.Keeper
	// this line is used by starport scaffolding # stargate/app/keeperDeclaration

	// mm is the module manager
	mm *module.Manager

	// sm is the simulation manager
	sm           *module.SimulationManager
	configurator module.Configurator
}

// New returns a reference to an initialized blockchain app
func NewArtela(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig appparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *Artela {
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	bApp := baseapp.NewBaseApp(
		Name,
		logger,
		db,
		encodingConfig.TxConfig.TxDecoder(),
		baseAppOptions...,
	)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	keys := cosmos.NewKVStoreKeys(
		authmodule.StoreKey, authz.ModuleName, bankmodule.StoreKey, stakingmodule.StoreKey,
		crisismodule.StoreKey, mintmodule.StoreKey, distrmodule.StoreKey, slashingmodule.StoreKey,
		govmodule.StoreKey, paramsmodule.StoreKey, ibcexported.StoreKey, upgrademodule.StoreKey,
		feegrant.StoreKey, evidencemodule.StoreKey, transfermodule.StoreKey, icahosttypes.StoreKey,
		capabilitymodule.StoreKey, group.StoreKey, icacontrollertypes.StoreKey, consensusmodule.StoreKey,
		evmmoduletypes.StoreKey,
		feemoduletypes.StoreKey,
		// this line is used by starport scaffolding # stargate/app/storeKey
	)
	tkeys := cosmos.NewTransientStoreKeys(paramsmodule.TStoreKey, evmmoduletypes.TransientKey, feemoduletypes.TransientKey)
	memKeys := cosmos.NewMemoryStoreKeys(capabilitymodule.MemStoreKey)

	app := &Artela{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		txConfig:          encodingConfig.TxConfig,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(
		appCodec,
		cdc,
		keys[paramsmodule.StoreKey],
		tkeys[paramsmodule.TStoreKey],
	)

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, keys[upgrademodule.StoreKey], authmodule.NewModuleAddress(govmodule.ModuleName).String())
	bApp.SetParamStore(&app.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitymodule.StoreKey],
		memKeys[capabilitymodule.MemStoreKey],
	)

	// set the runner cache capacity of aspect-runtime
	aspecttypes.InitRuntimePool(context.Background(), common.WrapLogger(app.Logger()), cast.ToInt32(appOpts.Get(srvflags.ApplyPoolSize)), cast.ToInt32(appOpts.Get(srvflags.QueryPoolSize)))

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedICAControllerKeeper := app.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(transfermodule.ModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	// this line is used by starport scaffolding # stargate/app/scopedKeeper

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authmodule.StoreKey],
		artela.ProtoAccount,
		maccPerms,
		cosmos.Bech32PrefixAccAddr,
		authmodule.NewModuleAddress(govmodule.ModuleName).String(),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authz.ModuleName],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		keys[bankmodule.StoreKey],
		app.AccountKeeper,
		app.BlockedModuleAccountAddrs(),
		authmodule.NewModuleAddress(govmodule.ModuleName).String(),
	)

	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingmodule.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		authmodule.NewModuleAddress(govmodule.ModuleName).String(),
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)

	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		keys[mintmodule.StoreKey],
		app.StakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		authmodule.FeeCollectorName,
		authmodule.NewModuleAddress(govmodule.ModuleName).String(),
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrmodule.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		authmodule.FeeCollectorName,
		authmodule.NewModuleAddress(govmodule.ModuleName).String(),
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		cdc,
		keys[slashingmodule.StoreKey],
		app.StakingKeeper,
		authmodule.NewModuleAddress(govmodule.ModuleName).String(),
	)

	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		keys[crisismodule.StoreKey],
		invCheckPeriod,
		app.BankKeeper,
		authmodule.FeeCollectorName,
		authmodule.NewModuleAddress(govmodule.ModuleName).String(),
	)

	groupConfig := group.DefaultConfig()
	/*
		Example of setting group params:
		groupConfig.MaxMetadataLen = 1000
	*/
	app.GroupKeeper = groupkeeper.NewKeeper(
		keys[group.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		groupConfig,
	)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgrademodule.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
		authmodule.NewModuleAddress(govmodule.ModuleName).String(),
	)

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibcexported.StoreKey],
		app.GetSubspace(ibcexported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[transfermodule.StoreKey],
		app.GetSubspace(transfermodule.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)

	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, keys[icahosttypes.StoreKey],
		app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)
	icaControllerKeeper := icacontrollerkeeper.NewKeeper(
		appCodec, keys[icacontrollertypes.StoreKey],
		app.GetSubspace(icacontrollertypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper, // may be replaced with middleware such as ics29 fee
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper, app.MsgServiceRouter(),
	)
	icaModule := ica.NewAppModule(&icaControllerKeeper, &app.ICAHostKeeper)
	icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		keys[evidencemodule.StoreKey],
		app.StakingKeeper,
		app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	govConfig := govmodule.DefaultConfig()
	govKeeper := govkeeper.NewKeeper(
		appCodec,
		keys[govmodule.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.MsgServiceRouter(),
		govConfig,
		authmodule.NewModuleAddress(govmodule.ModuleName).String(),
	)

	govRouter := govv1beta1.NewRouter()
	govRouter.
		AddRoute(govmodule.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(upgrademodule.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))
	govKeeper.SetLegacyRouter(govRouter)

	app.GovKeeper = *govKeeper.SetHooks(
		govmodule.NewMultiGovHooks(
		// register the governance hooks
		),
	)

	app.FeeKeeper = feemodulekeeper.NewKeeper(
		appCodec, authmodule.NewModuleAddress(govmodule.ModuleName),
		keys[feemoduletypes.StoreKey],
		tkeys[feemoduletypes.TransientKey],
		app.GetSubspace(feemoduletypes.ModuleName),
	)
	feeModule := feemodule.NewAppModule(app.FeeKeeper, app.GetSubspace(feemoduletypes.ModuleName))

	app.EvmKeeper = evmmodulekeeper.NewKeeper(
		appCodec, keys[evmmoduletypes.StoreKey], tkeys[evmmoduletypes.TransientKey], authmodule.NewModuleAddress(govmodule.ModuleName),
		app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.FeeKeeper,
		"1", app.GetSubspace(evmmoduletypes.ModuleName), bApp, logger,
	)
	evmModule := evmmodule.NewAppModule(app.EvmKeeper, app.AccountKeeper, app.GetSubspace(evmmoduletypes.ModuleName))

	// this line is used by starport scaffolding # stargate/app/keeperDefinition

	/**** IBC Routing ****/

	// Sealing prevents other modules from creating scoped sub-keepers
	app.CapabilityKeeper.Seal()

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		AddRoute(transfermodule.ModuleName, transferIBCModule)
	// this line is used by starport scaffolding # ibc/app/router
	app.IBCKeeper.SetRouter(ibcRouter)

	/**** Module Hooks ****/

	// register hooks after all modules have been initialized

	app.StakingKeeper.SetHooks(
		stakingmodule.NewMultiStakingHooks(
			// insert staking hooks receivers here
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks(),
		),
	)

	/**** Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper,
			app.StakingKeeper,
			app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authmodule.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(bankmodule.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		groupmodule.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govmodule.ModuleName)),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(mintmodule.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingmodule.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrmodule.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingmodule.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,
		icaModule,
		evmModule,
		feeModule,
		// this line is used by starport scaffolding # stargate/app/appModule

		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisismodule.ModuleName)), // always be last to make sure that it checks for all invariants and not only part of them
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		// upgrades should be run first
		upgrademodule.ModuleName,
		capabilitymodule.ModuleName,
		feemoduletypes.ModuleName,
		evmmoduletypes.ModuleName,
		mintmodule.ModuleName,
		distrmodule.ModuleName,
		slashingmodule.ModuleName,
		evidencemodule.ModuleName,
		stakingmodule.ModuleName,
		authmodule.ModuleName,
		bankmodule.ModuleName,
		govmodule.ModuleName,
		crisismodule.ModuleName,
		transfermodule.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		genutilmodule.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramsmodule.ModuleName,
		vestingtypes.ModuleName,
		consensusmodule.ModuleName,
		// this line is used by starport scaffolding # stargate/app/beginBlockers
	)

	app.mm.SetOrderEndBlockers(
		crisismodule.ModuleName,
		govmodule.ModuleName,
		stakingmodule.ModuleName,
		evmmoduletypes.ModuleName,
		feemoduletypes.ModuleName,
		transfermodule.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		capabilitymodule.ModuleName,
		authmodule.ModuleName,
		bankmodule.ModuleName,
		distrmodule.ModuleName,
		slashingmodule.ModuleName,
		mintmodule.ModuleName,
		genutilmodule.ModuleName,
		evidencemodule.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramsmodule.ModuleName,
		upgrademodule.ModuleName,
		vestingtypes.ModuleName,
		consensusmodule.ModuleName,
		// this line is used by starport scaffolding # stargate/app/endBlockers
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	genesisModuleOrder := []string{
		capabilitymodule.ModuleName,
		authmodule.ModuleName,
		bankmodule.ModuleName,
		distrmodule.ModuleName,
		stakingmodule.ModuleName,
		slashingmodule.ModuleName,
		govmodule.ModuleName,
		mintmodule.ModuleName,
		crisismodule.ModuleName,
		evmmoduletypes.ModuleName,
		feemoduletypes.ModuleName,
		genutilmodule.ModuleName,
		transfermodule.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		evidencemodule.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramsmodule.ModuleName,
		upgrademodule.ModuleName,
		vestingtypes.ModuleName,
		consensusmodule.ModuleName,
		// this line is used by starport scaffolding # stargate/app/initGenesis
	}
	app.mm.SetOrderInitGenesis(genesisModuleOrder...)
	app.mm.SetOrderExportGenesis(genesisModuleOrder...)

	// Uncomment if you want to set a custom migration order here.
	// app.mm.SetOrderMigrations(custom order)

	app.mm.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))
	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// create the simulation manager and define the order of the modules for deterministic simulations
	overrideModules := map[string]module.AppModuleSimulation{
		authmodule.ModuleName: auth.NewAppModule(app.appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authmodule.ModuleName)),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.mm.Modules, overrideModules)
	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	maxGasWanted := cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted))
	app.setAnteHandler(encodingConfig.TxConfig, maxGasWanted)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	// init Aspect
	app.setPostHandler()

	// // aspect add ProposalHandler
	// aspectProposalHandler := handle.NewArtelaProposalHandler(bApp.GetMemPool(), bApp)
	// bApp.SetPrepareProposal(aspectProposalHandler.PrepareProposalHandler())
	// bApp.SetProcessProposal(aspectProposalHandler.ProcessProposalHandler())
	// setupUpgradeHandlers should be called before `LoadLatestVersion()`
	// because StoreLoad is sealed after that
	// app upgrade
	app.setupUpgradeHandlers()

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	// this line is used by starport scaffolding # stargate/app/beforeInitReturn

	return app
}

// TODO mark
func (app *Artela) setAnteHandler(txConfig client.TxConfig, maxGasWanted uint64) {
	options := ante.AnteDecorators{
		Cdc:                    app.appCodec,
		AccountKeeper:          app.AccountKeeper,
		BankKeeper:             app.BankKeeper,
		ExtensionOptionChecker: artela.HasDynamicFeeExtensionOption,
		EvmKeeper:              app.EvmKeeper,
		FeegrantKeeper:         app.FeeGrantKeeper,
		DistributionKeeper:     app.DistrKeeper,
		FeeKeeper:              app.FeeKeeper,
		SignModeHandler:        txConfig.SignModeHandler(),
		SigGasConsumer:         ante.SigVerificationGasConsumer,
		MaxTxGasWanted:         maxGasWanted,
		TxFeeChecker:           ethante.NewDynamicFeeChecker(app.EvmKeeper),

		// TODO StakingKeeper:          app.StakingKeeper,
		IBCKeeper: app.IBCKeeper,
	}

	if err := options.Validate(); err != nil {
		panic(err)
	}

	app.SetAnteHandler(ante.NewAnteHandler(app.BaseApp, options))
}

func (app *Artela) setPostHandler() {
	options := post.PostDecorators{
		EvmKeeper: app.EvmKeeper,
	}

	if err := options.Validate(); err != nil {
		panic(err)
	}

	app.SetPostHandler(post.NewPostHandler(app.BaseApp, options))
}

// Name returns the name of the App
func (app *Artela) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *Artela) BeginBlocker(ctx cosmos.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *Artela) EndBlocker(ctx cosmos.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *Artela) InitChainer(ctx cosmos.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// Configurator get app configurator
func (app *Artela) Configurator() module.Configurator {
	return app.configurator
}

// LoadHeight loads a particular height
func (app *Artela) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *Artela) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authmodule.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedModuleAccountAddrs returns all the app's blocked module account
// addresses.
func (app *Artela) BlockedModuleAccountAddrs() map[string]bool {
	modAccAddrs := app.ModuleAccountAddrs()
	delete(modAccAddrs, authmodule.NewModuleAddress(govmodule.ModuleName).String())

	return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Artela) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns an app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Artela) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns an InterfaceRegistry
func (app *Artela) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// TxConfig returns SimApp's TxConfig
func (app *Artela) TxConfig() client.TxConfig {
	return app.txConfig
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Artela) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Artela) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *Artela) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Artela) GetSubspace(moduleName string) paramsmodule.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *Artela) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new txs routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register app's OpenAPI routes.
	docs.RegisterOpenAPIService(Name, apiSvr.Router)

	app.EvmKeeper.SetClientContext(clientCtx)
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *Artela) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *Artela) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// RegisterNodeService implements the Application.RegisterNodeService method.
func (app *Artela) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authmodule.ModuleName)
	paramsKeeper.Subspace(bankmodule.ModuleName)
	paramsKeeper.Subspace(stakingmodule.ModuleName)
	paramsKeeper.Subspace(mintmodule.ModuleName)
	paramsKeeper.Subspace(distrmodule.ModuleName)
	paramsKeeper.Subspace(slashingmodule.ModuleName)
	paramsKeeper.Subspace(govmodule.ModuleName).WithKeyTable(govv1.ParamKeyTable()) //nolint:staticcheck
	paramsKeeper.Subspace(crisismodule.ModuleName)
	paramsKeeper.Subspace(transfermodule.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(evmmoduletypes.ModuleName)
	paramsKeeper.Subspace(feemoduletypes.ModuleName)
	// this line is used by starport scaffolding # stargate/app/paramSubspace

	return paramsKeeper
}

// SimulationManager returns the app SimulationManager
func (app *Artela) SimulationManager() *module.SimulationManager {
	return app.sm
}

// ModuleManager returns the app ModuleManager
func (app *Artela) ModuleManager() *module.Manager {
	return app.mm
}

func (app *Artela) setupUpgradeHandlers() {
	// v0.4.7-rc7 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v047rc7.UpgradeName,
		v047rc7.CreateUpgradeHandler(
			app.mm, app.configurator, app.BankKeeper, app.AccountKeeper,
		),
	)

	// v0.4.8-rc8 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v048rc8.UpgradeName,
		v048rc8.CreateUpgradeHandler(
			app.mm, app.configurator,
		),
	)

	// When a planned update height is reached, the old binary will panic
	// writing on disk the height and name of the update that triggered it
	// This will read that value, and execute the preparations for the upgrade.
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Errorf("failed to read upgrade info from disk: %w", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	var storeUpgrades *storetypes.StoreUpgrades

	switch upgradeInfo.Name {
	case v047rc7.UpgradeName:
		// no store upgrades
	case v048rc8.UpgradeName:
		// no store upgrades in v048rc8
	default:
		// no-op
	}

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgrademodule.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
