package contract

import (
	"bytes"
	errorsmod "cosmossdk.io/errors"
	"github.com/artela-network/artela/x/evm/artela/types"
	evmtxs "github.com/artela-network/artela/x/evm/txs"
	evmtypes "github.com/artela-network/artela/x/evm/types"
	"github.com/artela-network/evm/vm"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
	"math/big"
	"strings"
)

type AspectNativeContract struct {
	aspectService    *AspectService
	applyMessageFunc func(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*evmtxs.MsgEthereumTxResponse, error)
}

func NewAspectNativeContract(storeKey storetypes.StoreKey,
	getCtxByHeight func(height int64, prove bool) (sdk.Context, error),
	applyMessageFunc func(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*evmtxs.MsgEthereumTxResponse, error)) *AspectNativeContract {

	return &AspectNativeContract{
		aspectService:    NewAspectService(storeKey, getCtxByHeight),
		applyMessageFunc: applyMessageFunc,
	}
}

func (k *AspectNativeContract) Query(ctx sdk.Context, tx *ethtypes.Transaction) (*evmtxs.MsgEthereumTxResponse, error) {
	if tx == nil {
		return nil, nil
	}
	method, params, err := types.ParseInput(tx.Data())
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(method.Name) {
	case "versionof":
		{
			aspectId := params["aspectId"].(common.Address)
			return k.version(ctx, tx, method, aspectId)
		}

	case "aspectsof":
		{
			contract := params["contract"].(common.Address)
			return k.aspectsOf(ctx, tx, method, contract)
		}

	case "contractsof":
		{
			aspectId := params["aspectId"].(common.Address)
			return k.contractsOf(ctx, tx, method, aspectId)
		}
	case "entrypoint":
		{
			aspectId := params["aspectId"].(common.Address)
			data := params["optArgs"].([]byte)
			return k.entrypoint(ctx, tx, method, aspectId, data)
		}
	}

	return nil, nil
}

func (k *AspectNativeContract) ApplyTx(ctx sdk.Context, tx *ethtypes.Transaction, msg core.Message) (*evmtxs.MsgEthereumTxResponse, error) {
	method, parameters, err := types.ParseMethod(tx)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(method.Name) {
	case "deploy":
		{
			code := parameters["code"].([]byte)
			properties := parameters["properties"].([]struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			})
			var propertyAry []types.Property
			for i := range properties { //nolint:typecheck
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

			proof := parameters["proof"].([]byte)
			if bytes.Equal(account.Bytes(), msg.From.Bytes()) {
				accountProperty := types.Property{
					Key:   types.AspectAccountKey,
					Value: account.Hex(),
				}
				propertyAry = append(propertyAry, accountProperty)
			} else {
				return nil, errors.New("Account verification failed during aspect deploy")
			}
			if len(proof) > 0 {
				proofProperty := types.Property{
					Key:   types.AspectProofKey,
					Value: common.Bytes2Hex(proof),
				}
				propertyAry = append(propertyAry, proofProperty)
			} else {
				return nil, errors.New("No proof provided during aspect deploy")
			}

			aspectId := crypto.CreateAddress(sender.Address(), msg.Nonce)
			return k.deploy(ctx, tx, aspectId, code, propertyAry)
		}
	case "upgrade":
		{
			aspectId := parameters["aspectId"].(common.Address)
			code := parameters["code"].([]byte)
			properties := parameters["properties"].([]struct {
				Key   string `json:"key"`
				Value string `json:"value"`
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
			aspectOwner, checkErr := k.checkAspectOwner(ctx, aspectId, sender.Address())
			if checkErr != nil {
				return nil, checkErr
			}
			if !aspectOwner {
				return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "failed to check if the sender is the owner, unable to upgrade, sender: %s , aspectId: %s", sender.Address().String(), aspectId.String())
			}
			return k.deploy(ctx, tx, aspectId, code, propertyAry)
		}
	case "bind":
		{
			aspectId := parameters["aspectId"].(common.Address)
			aspectVersion := parameters["aspectVersion"].(*big.Int)
			contract := parameters["contract"].(common.Address)
			priority := parameters["priority"].(int8)
			versionU256, _ := uint256.FromBig(aspectVersion)
			sender := vm.AccountRef(msg.From)
			owner := k.checkContractOwner(ctx, &contract, tx.Nonce()+1, sender.Address())
			if !owner {
				return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "check sender isOwner fail, sender: %s , contract: %s", sender.Address().String(), contract.String())
			}
			binding, err := k.checkContractBinding(ctx, aspectId, contract)
			if err != nil {
				return nil, err
			}
			if !binding {
				return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "check contract binding fail, aspectId: %s , contract: %s", aspectId.String(), contract.String())
			}

			return k.bind(ctx, tx, aspectId, contract, versionU256, priority)
		}

	case "unbind":
		{
			aspectId := parameters["aspectId"].(common.Address)
			contract := parameters["contract"].(common.Address)
			sender := vm.AccountRef(msg.From)
			owner := k.checkContractOwner(ctx, &contract, tx.Nonce()+1, sender.Address())
			aspectOwner, err := k.checkAspectOwner(ctx, aspectId, sender.Address())
			if err != nil {
				return nil, err
			}
			if !owner || !aspectOwner {
				return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "failed to check if the sender is the owner, unable to unbind, sender: %s , contract: %s", sender.Address().String(), contract.String())
			}
			return k.unbind(ctx, tx, aspectId, contract)

		}
	case "changeversion":
		{
			aspectId := parameters["aspectId"].(common.Address)
			contract := parameters["contract"].(common.Address)
			version := parameters["version"].(uint64)
			sender := vm.AccountRef(msg.From)
			aspectOwner, err := k.checkAspectOwner(ctx, aspectId, sender.Address())
			if err != nil {
				return nil, err
			}
			if !aspectOwner {
				return nil, errorsmod.Wrapf(evmtypes.ErrCallContract, "failed to check if the sender is the owner, unable to changeversion, sender: %s , contract: %s", sender.Address().String(), contract.String())
			}
			return k.changeVersion(ctx, tx, aspectId, contract, version)
		}
	case "versionof":
		{
			aspectId := parameters["aspectId"].(common.Address)
			return k.version(ctx, tx, method, aspectId)
		}

	case "aspectsof":
		{
			contract := parameters["contract"].(common.Address)
			return k.aspectsOf(ctx, tx, method, contract)
		}

	case "contractsof":
		{
			aspectId := parameters["aspectId"].(common.Address)
			return k.contractsOf(ctx, tx, method, aspectId)
		}
	case "entrypoint":
		{
			aspectId := parameters["aspectId"].(common.Address)
			data := parameters["optArgs"].([]byte)
			return k.entrypoint(ctx, tx, method, aspectId, data)
		}
	}
	return nil, nil
}
