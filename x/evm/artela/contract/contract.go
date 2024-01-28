package contract

import (
	"bytes"
	"math/big"
	"strings"

	"github.com/artela-network/artela/x/evm/states"

	errorsmod "cosmossdk.io/errors"
	"github.com/artela-network/artela-evm/vm"
	"github.com/cometbft/cometbft/libs/log"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"

	"github.com/artela-network/artela/x/evm/artela/types"
	evmtxs "github.com/artela-network/artela/x/evm/txs"
	evmtypes "github.com/artela-network/artela/x/evm/types"
)

type AspectNativeContract struct {
	aspectService *AspectService
	evmState      *states.StateDB
	evm           *vm.EVM
}

func NewAspectNativeContract(storeKey storetypes.StoreKey,
	evm *vm.EVM,
	getBlockHeight func() int64,
	evmState *states.StateDB,
	logger log.Logger,
) *AspectNativeContract {
	return &AspectNativeContract{
		aspectService: NewAspectService(storeKey, getBlockHeight, logger),
		evm:           evm,
		evmState:      evmState,
	}
}

func (k *AspectNativeContract) ApplyMessage(ctx sdk.Context, msg *core.Message, commit bool) (*evmtxs.MsgEthereumTxResponse, error) {
	var writeCacheFunc func()
	ctx, writeCacheFunc = ctx.CacheContext()
	applyMsg, err := k.applyMsg(ctx, msg, commit)
	if err == nil && commit {
		writeCacheFunc()
	}
	return applyMsg, err
}

func (k *AspectNativeContract) applyMsg(ctx sdk.Context, msg *core.Message, commit bool) (*evmtxs.MsgEthereumTxResponse, error) {
	method, parameters, err := types.ParseMethod(msg.Data)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(method.Name) {
	case "deploy":
		{
			code := parameters["code"].([]byte)
			properties := parameters["properties"].([]struct {
				Key   string `json:"key"`
				Value []byte `json:"value"`
			})
			var propertyAry []types.Property
			for i := range properties {
				s := properties[i]
				if types.AspectProofKey == s.Key || types.AspectAccountKey == s.Key {
					// Block query of account and Proof
					return nil, errors.New("Cannot use Aspect-defined keys")
				}
				propertyAry = append(propertyAry, types.Property{
					Key:   s.Key,
					Value: s.Value,
				})
			}
			sender := vm.AccountRef(msg.From)
			account := parameters["account"].(common.Address)
			if bytes.Equal(account.Bytes(), msg.From.Bytes()) {
				accountProperty := types.Property{
					Key:   types.AspectAccountKey,
					Value: account.Bytes(),
				}
				propertyAry = append(propertyAry, accountProperty)
			} else {
				return nil, errors.New("Account verification failed during aspect deploy")
			}
			joinPoints := parameters["joinPoints"].(*big.Int)

			if len(code) == 0 {
				return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "Code verification failed during aspect deploy")
			}
			if joinPoints == nil {
				return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "JoinPoints verification failed during aspect deploy")
			}

			aspectId := crypto.CreateAddress(sender.Address(), msg.Nonce)
			return k.deploy(ctx, aspectId, code, propertyAry, joinPoints)
		}
	case "upgrade":
		{
			aspectId := parameters["aspectId"].(common.Address)
			code := parameters["code"].([]byte)
			properties := parameters["properties"].([]struct {
				Key   string `json:"key"`
				Value []byte `json:"value"`
			})
			var propertyAry []types.Property
			for i := range properties {
				s := properties[i]
				if types.AspectProofKey == s.Key || types.AspectAccountKey == s.Key {
					// Block query of account and Proof
					return nil, errors.New("Cannot use Aspect-defined keys")
				}
				propertyAry = append(propertyAry, types.Property{
					Key:   s.Key,
					Value: s.Value,
				})
			}
			sender := vm.AccountRef(msg.From)
			aspectOwner, checkErr := k.checkAspectOwner(ctx, aspectId, sender.Address(), commit)
			if checkErr != nil {
				return nil, checkErr
			}
			if !aspectOwner {
				return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "failed to check if the sender is the owner, unable to upgrade, sender: %s , aspectId: %s", sender.Address().String(), aspectId.String())
			}
			joinPoints := parameters["joinPoints"].(*big.Int)
			return k.deploy(ctx, aspectId, code, propertyAry, joinPoints)
		}
	case "bind":
		{
			aspectId := parameters["aspectId"].(common.Address)
			aspectVersion := parameters["aspectVersion"].(*big.Int)
			account := parameters["contract"].(common.Address)
			priority := parameters["priority"].(int8)
			versionU256, _ := uint256.FromBig(aspectVersion)
			sender := vm.AccountRef(msg.From)
			isContract := len(k.evmState.GetCode(account)) > 0
			if isContract {
				cOwner := k.checkContractOwner(ctx, &account, msg.Nonce+1, sender.Address())
				// Bind with contract account, need to verify contract ownerships first
				// owner, _ := k.checkAspectOwner(ctx, aspectId, sender.Address(), commit)
				if !(cOwner) {
					return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "check sender isOwner fail, sender: %s , contract: %s", sender.Address().String(), account.String())
				}
			} else if account != sender.Address() {
				// For EoA account binding, only the account itself can issue the bind request
				return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "unauthorized EoA account aspect binding")
			}

			return k.bind(ctx, aspectId, account, versionU256, priority, isContract)
		}

	case "unbind":
		{
			aspectId := parameters["aspectId"].(common.Address)
			account := parameters["contract"].(common.Address)
			sender := vm.AccountRef(msg.From)
			isContract := len(k.evmState.GetCode(account)) > 0
			if isContract {
				cOwner := k.checkContractOwner(ctx, &account, msg.Nonce+1, sender.Address())
				// Bind with contract account, need to verify contract ownerships first
				// owner, _ := k.checkAspectOwner(ctx, aspectId, sender.Address(), commit)
				if !(cOwner) {
					return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "check sender isOwner fail, sender: %s , contract: %s", sender.Address().String(), account.String())
				}
			} else if account != sender.Address() {
				// For EoA account binding, only the account itself can issue the bind request
				return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "unauthorized EoA account aspect unbinding")
			}

			return k.unbind(ctx, aspectId, account, isContract)

		}
	case "changeversion":
		{
			aspectId := parameters["aspectId"].(common.Address)
			contract := parameters["contract"].(common.Address)
			version := parameters["version"].(uint64)
			sender := vm.AccountRef(msg.From)
			aspectOwner, err := k.checkAspectOwner(ctx, aspectId, sender.Address(), commit)
			if err != nil {
				return nil, err
			}
			if !aspectOwner {
				return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "failed to check if the sender is the owner, unable to changeversion, sender: %s , contract: %s", sender.Address().String(), contract.String())
			}
			return k.changeVersion(ctx, aspectId, contract, version)
		}
	case "versionof":
		{
			aspectId := parameters["aspectId"].(common.Address)
			return k.version(ctx, method, aspectId)
		}

	case "aspectsof":
		{
			account := parameters["contract"].(common.Address)
			isContract := len(k.evmState.GetCode(account)) > 0

			return k.aspectsOf(ctx, method, account, isContract)
		}

	case "boundaddressesof":
		{
			aspectId := parameters["aspectId"].(common.Address)
			return k.boundAddressesOf(ctx, method, aspectId)
		}
	case "entrypoint":
		{
			aspectId := parameters["aspectId"].(common.Address)
			data := parameters["optArgs"].([]byte)
			return k.entrypoint(ctx, msg, method, aspectId, data, commit)
		}
	}
	return nil, nil
}
