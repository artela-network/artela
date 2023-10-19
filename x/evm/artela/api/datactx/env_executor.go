package datactx

import (
	"github.com/artela-network/artela/x/evm/artela/types"
	"github.com/artela-network/artela/x/evm/txs/support"
	artelatypes "github.com/artela-network/artelasdk/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/params"
)

type EnvBaseInfo struct {
	getEthTxContext func() *types.EthTxContext
}

func NewEnvBaseInfo(getEthTxContext func() *types.EthTxContext) *EnvBaseInfo {
	return &EnvBaseInfo{getEthTxContext: getEthTxContext}
}
func (c EnvBaseInfo) Execute(sdkContext sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}
	txCtx := c.getEthTxContext()
	if txCtx == nil || txCtx.EvmCfg() == nil {
		return nil
	}
	fee := txCtx.EvmCfg().BaseFee
	content := &artelatypes.EnvContent{BaseFee: fee.Uint64()}
	response := artelatypes.NewContextQueryResponse(true, "success")
	response.SetData(content)
	return response
}

type EnvEvmParams struct {
	getEthTxContext func() *types.EthTxContext
}

func NewEnvEvmParams(getEthTxContext func() *types.EthTxContext) *EnvEvmParams {
	return &EnvEvmParams{getEthTxContext: getEthTxContext}
}
func (c EnvEvmParams) Execute(sdkContext sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}
	txCtx := c.getEthTxContext()

	if txCtx == nil || txCtx.EvmCfg() == nil {
		return nil
	}
	chainConfig := ConvertEvmParams(txCtx.EvmCfg().Params)
	response := artelatypes.NewContextQueryResponse(true, "success")
	response.SetData(chainConfig)
	return response
}

func ConvertEvmParams(ctx support.Params) *artelatypes.EvmParams {
	return &artelatypes.EvmParams{
		EvmDenom:            ctx.EvmDenom,
		EnableCreate:        ctx.EnableCreate,
		EnableCall:          ctx.EnableCall,
		ExtraEips:           ctx.ExtraEIPs,
		AllowUnprotectedTxs: ctx.AllowUnprotectedTxs,
	}
}

type EnvChainConfig struct {
	getEthTxContext func() *types.EthTxContext
}

func NewEnvChainConfig(getEthTxContext func() *types.EthTxContext) *EnvChainConfig {
	return &EnvChainConfig{getEthTxContext: getEthTxContext}
}

func (c EnvChainConfig) Execute(sdkContext sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil {
		return nil
	}
	txCtx := c.getEthTxContext()

	if txCtx == nil || txCtx.EvmCfg() == nil || txCtx.EvmCfg().ChainConfig == nil {
		return nil
	}

	chainConfig := ConvertChainConfig(txCtx.EvmCfg().ChainConfig)
	response := artelatypes.NewContextQueryResponse(true, "success")
	response.SetData(chainConfig)
	return response
}
func ConvertChainConfig(param *params.ChainConfig) *artelatypes.ChainConfig {
	if param == nil {
		return nil
	}
	clique := &artelatypes.Clique{}
	if param.Clique != nil {
		clique.Epoch = param.Clique.Epoch
		clique.Period = param.Clique.Period
	}
	return &artelatypes.ChainConfig{
		ChainId:                       param.ChainID.String(),
		HomesteadBlock:                param.HomesteadBlock.String(),
		DaoForkBlock:                  param.DAOForkBlock.String(),
		DaoForkSupport:                param.DAOForkSupport,
		Eip150Block:                   param.EIP150Block.String(),
		Eip155Block:                   param.EIP155Block.String(),
		Eip158Block:                   param.EIP158Block.String(),
		ByzantiumBlock:                param.ByzantiumBlock.String(),
		ConstantinopleBlock:           param.ConstantinopleBlock.String(),
		PetersburgBlock:               param.PetersburgBlock.String(),
		IstanbulBlock:                 param.IstanbulBlock.String(),
		MuirGlacierBlock:              param.MuirGlacierBlock.String(),
		BerlinBlock:                   param.BerlinBlock.String(),
		LondonBlock:                   param.LondonBlock.String(),
		ArrowGlacierBlock:             param.ArrowGlacierBlock.String(),
		GrayGlacierBlock:              param.GrayGlacierBlock.String(),
		MergeNetsplitBlock:            param.MergeNetsplitBlock.String(),
		ShanghaiTime:                  *param.ShanghaiTime,
		CancunTime:                    *param.CancunTime,
		PragueTime:                    *param.PragueTime,
		TerminalTotalDifficulty:       param.TerminalTotalDifficulty.String(),
		TerminalTotalDifficultyPassed: param.TerminalTotalDifficultyPassed,
		Clique:                        clique,
	}
}

type EnvConsParams struct {
}

func NewEnvConsParams() *EnvConsParams {
	return &EnvConsParams{}
}

// getAspectContext data
func (c EnvConsParams) Execute(sdkCtx sdk.Context, ctx *artelatypes.RunnerContext, keys []string) *artelatypes.ContextQueryResponse {
	if ctx == nil || ctx.ContractAddr == nil || ctx.AspectId == nil || sdkCtx.ConsensusParams() == nil {
		return nil
	}
	consParams := ConvertConsParams(sdkCtx.ConsensusParams())
	response := artelatypes.NewContextQueryResponse(true, "success")
	response.SetData(consParams)
	return response
}

func ConvertConsParams(param *tmproto.ConsensusParams) *artelatypes.ConsParams {
	if param == nil {
		return nil
	}
	blockParams := &artelatypes.BlockParams{}
	if param.Block != nil {
		blockParams.MaxGas = param.Block.MaxGas
		blockParams.MaxBytes = param.Block.MaxBytes
	}
	evidenceParams := &artelatypes.EvidenceParams{}
	if param.Evidence != nil {
		evidenceParams.MaxAgeNumBlocks = param.Evidence.MaxAgeNumBlocks
		evidenceParams.MaxAgeDuration = param.Evidence.MaxAgeDuration.Microseconds()
		evidenceParams.MaxBytes = param.Evidence.MaxBytes
	}
	validatorParams := &artelatypes.ValidatorParams{}
	if param.Validator != nil {
		validatorParams.PubKeyTypes = param.Validator.PubKeyTypes
	}
	versionParams := &artelatypes.VersionParams{}
	if param.Version != nil {
		versionParams.AppVersion = param.Version.App
	}
	return &artelatypes.ConsParams{
		Block:     blockParams,
		Evidence:  evidenceParams,
		Validator: validatorParams,
		Version:   versionParams,
	}
}
