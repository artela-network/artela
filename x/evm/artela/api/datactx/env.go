package datactx

import (
	"strings"

	"github.com/artela-network/aspect-core/context"
	artelatypes "github.com/artela-network/aspect-core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/proto"

	"github.com/artela-network/artela/x/evm/artela/types"
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
	defUint64 := uint64(0)
	loaders[context.EnvExtraEIPs] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.IntArrayData{Data: ethTxCtx.EvmCfg().Params.ExtraEIPs}
	}
	loaders[context.EnvEnableCreate] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.BoolData{}
		}
		return &artelatypes.BoolData{Data: &ethTxCtx.EvmCfg().Params.EnableCreate}
	}
	loaders[context.EnvEnableCall] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.BoolData{}
		}
		return &artelatypes.BoolData{Data: &ethTxCtx.EvmCfg().Params.EnableCall}
	}
	loaders[context.EnvAllowUnprotectedTxs] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.BoolData{}
		}
		return &artelatypes.BoolData{Data: &ethTxCtx.EvmCfg().Params.AllowUnprotectedTxs}
	}
	loaders[context.EnvChainChainId] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		chainId := c.keeper.ChainID().Uint64()
		return &artelatypes.UintData{Data: &chainId}
	}
	loaders[context.EnvChainHomesteadBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		homesteadBlock := ethTxCtx.EvmCfg().ChainConfig.HomesteadBlock.Uint64()
		return &artelatypes.UintData{Data: &homesteadBlock}
	}
	loaders[context.EnvChainDaoForkBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		forkBlock := ethTxCtx.EvmCfg().ChainConfig.DAOForkBlock.Uint64()
		return &artelatypes.UintData{Data: &forkBlock}
	}
	loaders[context.EnvChainDaoForkSupport] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.BoolData{}
		}
		forkBlock := ethTxCtx.EvmCfg().ChainConfig.DAOForkSupport
		return &artelatypes.BoolData{Data: &forkBlock}
	}
	loaders[context.EnvChainEip150Block] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		eip150Block := ethTxCtx.EvmCfg().ChainConfig.EIP150Block.Uint64()
		return &artelatypes.UintData{Data: &eip150Block}
	}
	loaders[context.EnvChainEip155Block] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		eip155block := ethTxCtx.EvmCfg().ChainConfig.EIP155Block.Uint64()
		return &artelatypes.UintData{Data: &eip155block}
	}
	loaders[context.EnvChainEip158Block] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		eip158block := ethTxCtx.EvmCfg().ChainConfig.EIP158Block.Uint64()
		return &artelatypes.UintData{Data: &eip158block}
	}
	loaders[context.EnvChainByzantiumBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		byzantium := ethTxCtx.EvmCfg().ChainConfig.ByzantiumBlock.Uint64()

		return &artelatypes.UintData{Data: &byzantium}
	}
	loaders[context.EnvChainConstantinopleBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		constantinople := ethTxCtx.EvmCfg().ChainConfig.ConstantinopleBlock.Uint64()
		return &artelatypes.UintData{Data: &constantinople}
	}
	loaders[context.EnvChainPetersburgBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		petersburgBlock := ethTxCtx.EvmCfg().ChainConfig.PetersburgBlock.Uint64()
		return &artelatypes.UintData{Data: &petersburgBlock}
	}
	loaders[context.EnvChainIstanbulBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		istanbul := ethTxCtx.EvmCfg().ChainConfig.IstanbulBlock.Uint64()
		return &artelatypes.UintData{Data: &istanbul}
	}
	loaders[context.EnvChainMuirGlacierBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		muirGlacier := ethTxCtx.EvmCfg().ChainConfig.MuirGlacierBlock.Uint64()
		return &artelatypes.UintData{Data: &muirGlacier}
	}
	loaders[context.EnvChainBerlinBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		berlinBlock := ethTxCtx.EvmCfg().ChainConfig.BerlinBlock.Uint64()
		return &artelatypes.UintData{Data: &berlinBlock}
	}
	loaders[context.EnvChainLondonBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		londonBlock := ethTxCtx.EvmCfg().ChainConfig.LondonBlock.Uint64()
		return &artelatypes.UintData{Data: &londonBlock}
	}
	loaders[context.EnvChainArrowGlacierBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		arrowGlacierBlock := ethTxCtx.EvmCfg().ChainConfig.ArrowGlacierBlock.Uint64()
		return &artelatypes.UintData{Data: &arrowGlacierBlock}
	}
	loaders[context.EnvChainGrayGlacierBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		grayGlacierBlock := ethTxCtx.EvmCfg().ChainConfig.GrayGlacierBlock.Uint64()
		return &artelatypes.UintData{Data: &grayGlacierBlock}
	}
	loaders[context.EnvChainMergeNetSplitBlock] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		mergeNetsplitBlock := ethTxCtx.EvmCfg().ChainConfig.MergeNetsplitBlock.Uint64()
		return &artelatypes.UintData{Data: &mergeNetsplitBlock}
	}
	loaders[context.EnvChainShanghaiTime] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		shanghaiTime := ethTxCtx.EvmCfg().ChainConfig.ShanghaiTime
		if shanghaiTime != nil {
			return &artelatypes.UintData{Data: shanghaiTime}
		} else {
			return &artelatypes.UintData{Data: &defUint64}
		}

	}
	loaders[context.EnvChainCancunTime] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		cancunTime := ethTxCtx.EvmCfg().ChainConfig.CancunTime
		if cancunTime != nil {
			return &artelatypes.UintData{Data: cancunTime}
		} else {
			return &artelatypes.UintData{Data: &defUint64}
		}

	}
	loaders[context.EnvChainPragueTime] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {

		if ethTxCtx.EvmCfg() == nil {
			return &artelatypes.UintData{}
		}
		pragueTime := ethTxCtx.EvmCfg().ChainConfig.PragueTime
		if pragueTime != nil {
			return &artelatypes.UintData{Data: pragueTime}
		} else {
			return &artelatypes.UintData{Data: &defUint64}

		}

	}
	loaders[context.EnvConsensusParamsBlockMaxGas] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if sdkCtx.ConsensusParams() == nil {
			return &artelatypes.IntData{}
		}
		return &artelatypes.IntData{Data: &sdkCtx.ConsensusParams().Block.MaxGas}
	}
	loaders[context.EnvConsensusParamsBlockMaxBytes] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if sdkCtx.ConsensusParams() == nil {
			return &artelatypes.IntData{}
		}
		return &artelatypes.IntData{Data: &sdkCtx.ConsensusParams().Block.MaxBytes}
	}
	loaders[context.EnvConsensusParamsValidatorPubKeyTypes] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.StringArrayData{Data: sdkCtx.ConsensusParams().Validator.PubKeyTypes}
	}
	loaders[context.EnvConsensusParamsEvidenceMaxBytes] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if sdkCtx.ConsensusParams() == nil {
			return &artelatypes.IntData{}
		}
		return &artelatypes.IntData{Data: &sdkCtx.ConsensusParams().Evidence.MaxBytes}
	}
	loaders[context.EnvConsensusParamsEvidenceMaxAgeNumBlocks] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if sdkCtx.ConsensusParams() == nil {
			return &artelatypes.IntData{}
		}
		return &artelatypes.IntData{Data: &sdkCtx.ConsensusParams().Evidence.MaxAgeNumBlocks}
	}
	loaders[context.EnvConsensusParamsEvidenceMaxAgeDuration] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if sdkCtx.ConsensusParams() == nil {
			return &artelatypes.IntData{}
		}
		data := sdkCtx.ConsensusParams().Evidence.MaxAgeDuration.Milliseconds() / 1000
		return &artelatypes.IntData{Data: &data}
	}
	loaders[context.EnvConsensusParamsAppVersion] = func(sdkCtx sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		if sdkCtx.ConsensusParams() == nil || sdkCtx.ConsensusParams().Version == nil {
			return &artelatypes.IntData{}
		}
		return &artelatypes.UintData{Data: &sdkCtx.ConsensusParams().Version.App}
	}
}

func (c *EnvContext) ValueLoader(key string) ContextLoader {
	return func(_ *artelatypes.RunnerContext) ([]byte, error) {
		if c.ctx.EthTxContext() == nil || c.ctx.EthTxContext().EvmCfg() == nil {
			return nil, nil
		}

		if strings.HasPrefix(key, "env.consensusParams.") && c.ctx.CosmosContext().ConsensusParams() == nil {
			// when it is in an eth call, we are not able to return env.consensusParams.* values
			return nil, nil

		}

		return proto.Marshal(c.envContentLoaders[key](c.ctx.CosmosContext(), c.ctx.EthTxContext()))
	}
}
