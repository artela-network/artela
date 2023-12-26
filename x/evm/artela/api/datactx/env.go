package datactx

import (
	"github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/aspect-core/context"
	artelatypes "github.com/artela-network/aspect-core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/proto"
)

type EnvContextFieldLoader func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message

type EnvContext struct {
	envContentLoaders map[string]EnvContextFieldLoader
	ctx               *types.AspectRuntimeContext
	keeper            EVMKeeper
}

func NewEnvContext(ctx *types.AspectRuntimeContext, keeper EVMKeeper) *EnvContext {
	envCtx := &EnvContext{
		envContentLoaders: make(map[string]EnvContextFieldLoader),
		ctx:               ctx,
		keeper:            keeper,
	}
	envCtx.registerLoaders()
	return envCtx
}

func (c *EnvContext) registerLoaders() {
	loaders := c.envContentLoaders
	loaders[context.EnvExtraEIPs] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.IntArrayData{Data: ethTxCtx.EvmCfg().Params.ExtraEIPs}
	}
	loaders[context.EnvEnableCreate] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.BoolData{Data: ethTxCtx.EvmCfg().Params.EnableCreate}
	}
	loaders[context.EnvEnableCall] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.BoolData{Data: ethTxCtx.EvmCfg().Params.EnableCall}
	}
	loaders[context.EnvAllowUnprotectedTxs] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.BoolData{Data: ethTxCtx.EvmCfg().Params.AllowUnprotectedTxs}
	}
	loaders[context.EnvChainChainId] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: c.keeper.ChainID().Uint64()}
	}
	loaders[context.EnvChainHomesteadBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.HomesteadBlock.Uint64()}
	}
	loaders[context.EnvChainDaoForkBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.DAOForkBlock.Uint64()}
	}
	loaders[context.EnvChainDaoForkSupport] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.BoolData{Data: ethTxCtx.EvmCfg().ChainConfig.DAOForkSupport}
	}
	loaders[context.EnvChainEip150Block] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.EIP150Block.Uint64()}
	}
	loaders[context.EnvChainEip155Block] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.EIP155Block.Uint64()}
	}
	loaders[context.EnvChainEip158Block] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.EIP158Block.Uint64()}
	}
	loaders[context.EnvChainByzantiumBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.ByzantiumBlock.Uint64()}
	}
	loaders[context.EnvChainConstantinopleBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.ConstantinopleBlock.Uint64()}
	}
	loaders[context.EnvChainPetersburgBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.PetersburgBlock.Uint64()}
	}
	loaders[context.EnvChainIstanbulBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.IstanbulBlock.Uint64()}
	}
	loaders[context.EnvChainMuirGlacierBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.MuirGlacierBlock.Uint64()}
	}
	loaders[context.EnvChainBerlinBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.BerlinBlock.Uint64()}
	}
	loaders[context.EnvChainLondonBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.LondonBlock.Uint64()}
	}
	loaders[context.EnvChainArrowGlacierBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.ArrowGlacierBlock.Uint64()}
	}
	loaders[context.EnvChainGrayGlacierBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.GrayGlacierBlock.Uint64()}
	}
	loaders[context.EnvChainMergeNetSplitBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: ethTxCtx.EvmCfg().ChainConfig.MergeNetsplitBlock.Uint64()}
	}
	loaders[context.EnvChainShanghaiTime] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg().ChainConfig.ShanghaiTime != nil {
			return &artelatypes.UintData{Data: *ethTxCtx.EvmCfg().ChainConfig.ShanghaiTime}
		} else {
			return &artelatypes.UintData{Data: 0}
		}
	}
	loaders[context.EnvChainCancunTime] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg().ChainConfig.CancunTime != nil {
			return &artelatypes.UintData{Data: *ethTxCtx.EvmCfg().ChainConfig.CancunTime}
		} else {
			return &artelatypes.UintData{Data: 0}
		}
	}
	loaders[context.EnvChainPragueTime] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg().ChainConfig.PragueTime != nil {
			return &artelatypes.UintData{Data: *ethTxCtx.EvmCfg().ChainConfig.PragueTime}
		} else {
			return &artelatypes.UintData{Data: 0}
		}
	}
	loaders[context.EnvConsensusParamsBlockMaxGas] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.IntData{Data: sdkCtx.ConsensusParams().Block.MaxGas}
	}
	loaders[context.EnvConsensusParamsBlockMaxBytes] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.IntData{Data: sdkCtx.ConsensusParams().Block.MaxBytes}
	}
	loaders[context.EnvConsensusParamsValidatorPubKeyTypes] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.StringArrayData{Data: sdkCtx.ConsensusParams().Validator.PubKeyTypes}
	}
	loaders[context.EnvConsensusParamsEvidenceMaxBytes] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.IntData{Data: sdkCtx.ConsensusParams().Evidence.MaxBytes}
	}
	loaders[context.EnvConsensusParamsEvidenceMaxAgeNumBlocks] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.IntData{Data: sdkCtx.ConsensusParams().Evidence.MaxAgeNumBlocks}
	}
	loaders[context.EnvConsensusParamsEvidenceMaxAgeDuration] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.IntData{Data: sdkCtx.ConsensusParams().Evidence.MaxAgeDuration.Milliseconds() / 1000}
	}
	loaders[context.EnvConsensusParamsAppVersion] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.UintData{Data: sdkCtx.ConsensusParams().Version.App}
	}
}

func (c *EnvContext) ValueLoader(key string) ContextLoader {
	return func(_ *artelatypes.RunnerContext) ([]byte, error) {
		if c.ctx.EthTxContext() == nil || c.ctx.EthTxContext().EvmCfg() == nil {
			return nil, nil
		}

		return proto.Marshal(c.envContentLoaders[key](c.ctx.CosmosContext(), c.ctx.EthTxContext()))
	}
}
