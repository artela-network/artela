package datactx

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/proto"

	"github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/aspect-core/context"
	artelatypes "github.com/artela-network/aspect-core/types"
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
	loaders[context.EnvExtraEIPs] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.IntArrayData{Data: ethTxCtx.EvmCfg().Params.ExtraEIPs}
	}
	loaders[context.EnvEnableCreate] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.BoolData{Data: &ethTxCtx.EvmCfg().Params.EnableCreate}
	}
	loaders[context.EnvEnableCall] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.BoolData{Data: &ethTxCtx.EvmCfg().Params.EnableCall}
	}
	loaders[context.EnvAllowUnprotectedTxs] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		return &artelatypes.BoolData{Data: &ethTxCtx.EvmCfg().Params.AllowUnprotectedTxs}
	}
	loaders[context.EnvChainChainId] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		chainId := c.keeper.ChainID().Uint64()
		return &artelatypes.UintData{Data: &chainId}
	}
	loaders[context.EnvChainHomesteadBlock] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		homesteadBlock := ethTxCtx.EvmCfg().ChainConfig.HomesteadBlock.Uint64()
		return &artelatypes.UintData{Data: &homesteadBlock}
	}
	loaders[context.EnvChainDaoForkBlock] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		forkBlock := ethTxCtx.EvmCfg().ChainConfig.DAOForkBlock.Uint64()
		return &artelatypes.UintData{Data: &forkBlock}
	}
	loaders[context.EnvChainDaoForkSupport] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		forkBlock := ethTxCtx.EvmCfg().ChainConfig.DAOForkSupport
		return &artelatypes.BoolData{Data: &forkBlock}
	}
	loaders[context.EnvChainEip150Block] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		eip150Block := ethTxCtx.EvmCfg().ChainConfig.EIP150Block.Uint64()
		return &artelatypes.UintData{Data: &eip150Block}
	}
	loaders[context.EnvChainEip155Block] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		eip155block := ethTxCtx.EvmCfg().ChainConfig.EIP155Block.Uint64()
		return &artelatypes.UintData{Data: &eip155block}
	}
	loaders[context.EnvChainEip158Block] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		eip158block := ethTxCtx.EvmCfg().ChainConfig.EIP158Block.Uint64()
		return &artelatypes.UintData{Data: &eip158block}
	}
	loaders[context.EnvChainByzantiumBlock] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		byzantium := ethTxCtx.EvmCfg().ChainConfig.ByzantiumBlock.Uint64()

		return &artelatypes.UintData{Data: &byzantium}
	}
	loaders[context.EnvChainConstantinopleBlock] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		constantinople := ethTxCtx.EvmCfg().ChainConfig.ConstantinopleBlock.Uint64()
		return &artelatypes.UintData{Data: &constantinople}
	}
	loaders[context.EnvChainPetersburgBlock] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		petersburgBlock := ethTxCtx.EvmCfg().ChainConfig.PetersburgBlock.Uint64()
		return &artelatypes.UintData{Data: &petersburgBlock}
	}
	loaders[context.EnvChainIstanbulBlock] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		istanbul := ethTxCtx.EvmCfg().ChainConfig.IstanbulBlock.Uint64()
		return &artelatypes.UintData{Data: &istanbul}
	}
	loaders[context.EnvChainMuirGlacierBlock] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		muirGlacier := ethTxCtx.EvmCfg().ChainConfig.MuirGlacierBlock.Uint64()
		return &artelatypes.UintData{Data: &muirGlacier}
	}
	loaders[context.EnvChainBerlinBlock] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		berlinBlock := ethTxCtx.EvmCfg().ChainConfig.BerlinBlock.Uint64()
		return &artelatypes.UintData{Data: &berlinBlock}
	}
	loaders[context.EnvChainLondonBlock] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		londonBlock := ethTxCtx.EvmCfg().ChainConfig.LondonBlock.Uint64()
		return &artelatypes.UintData{Data: &londonBlock}
	}
	loaders[context.EnvChainArrowGlacierBlock] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		arrowGlacierBlock := ethTxCtx.EvmCfg().ChainConfig.ArrowGlacierBlock.Uint64()
		return &artelatypes.UintData{Data: &arrowGlacierBlock}
	}
	loaders[context.EnvChainGrayGlacierBlock] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		grayGlacierBlock := ethTxCtx.EvmCfg().ChainConfig.GrayGlacierBlock.Uint64()
		return &artelatypes.UintData{Data: &grayGlacierBlock}
	}
	loaders[context.EnvChainMergeNetSplitBlock] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		mergeNetsplitBlock := ethTxCtx.EvmCfg().ChainConfig.MergeNetsplitBlock.Uint64()
		return &artelatypes.UintData{Data: &mergeNetsplitBlock}
	}
	loaders[context.EnvChainShanghaiTime] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		shanghaiTime := ethTxCtx.EvmCfg().ChainConfig.ShanghaiTime
		if shanghaiTime != nil {
			return &artelatypes.UintData{Data: shanghaiTime}
		}
		return &artelatypes.UintData{Data: &defUint64}
	}
	loaders[context.EnvChainCancunTime] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		cancunTime := ethTxCtx.EvmCfg().ChainConfig.CancunTime
		if cancunTime != nil {
			return &artelatypes.UintData{Data: cancunTime}
		} else {
			return &artelatypes.UintData{Data: &defUint64}
		}
	}
	loaders[context.EnvChainPragueTime] = func(_ sdk.Context, ethTxCtx *types.EthTxContext) proto.Message {
		pragueTime := ethTxCtx.EvmCfg().ChainConfig.PragueTime
		if pragueTime != nil {
			return &artelatypes.UintData{Data: pragueTime}
		}
		return &artelatypes.UintData{Data: &defUint64}
	}
	loaders[context.EnvConsensusParamsBlockMaxGas] = func(sdkCtx sdk.Context, _ *types.EthTxContext) proto.Message {
		return &artelatypes.IntData{Data: &sdkCtx.ConsensusParams().Block.MaxGas}
	}
	loaders[context.EnvConsensusParamsBlockMaxBytes] = func(sdkCtx sdk.Context, _ *types.EthTxContext) proto.Message {
		return &artelatypes.IntData{Data: &sdkCtx.ConsensusParams().Block.MaxBytes}
	}
	loaders[context.EnvConsensusParamsValidatorPubKeyTypes] = func(sdkCtx sdk.Context, _ *types.EthTxContext) proto.Message {
		return &artelatypes.StringArrayData{Data: sdkCtx.ConsensusParams().Validator.PubKeyTypes}
	}
	loaders[context.EnvConsensusParamsEvidenceMaxBytes] = func(sdkCtx sdk.Context, _ *types.EthTxContext) proto.Message {
		return &artelatypes.IntData{Data: &sdkCtx.ConsensusParams().Evidence.MaxBytes}
	}
	loaders[context.EnvConsensusParamsEvidenceMaxAgeNumBlocks] = func(sdkCtx sdk.Context, _ *types.EthTxContext) proto.Message {
		return &artelatypes.IntData{Data: &sdkCtx.ConsensusParams().Evidence.MaxAgeNumBlocks}
	}
	loaders[context.EnvConsensusParamsEvidenceMaxAgeDuration] = func(sdkCtx sdk.Context, _ *types.EthTxContext) proto.Message {
		data := sdkCtx.ConsensusParams().Evidence.MaxAgeDuration.Milliseconds() / 1000
		return &artelatypes.IntData{Data: &data}
	}
	loaders[context.EnvConsensusParamsAppVersion] = func(sdkCtx sdk.Context, _ *types.EthTxContext) proto.Message {
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
