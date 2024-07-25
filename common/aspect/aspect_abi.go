package aspect

import (
	"encoding/hex"

	errorsmod "cosmossdk.io/errors"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artela-network/aspect-core/types"
)

var (
	Uint256, _ = abi.NewType("uint256", "", nil)
	Uint64, _  = abi.NewType("uint64", "", nil)
	Uint32, _  = abi.NewType("uint32", "", nil)
	Uint16, _  = abi.NewType("uint16", "", nil)
	Uint8, _   = abi.NewType("uint8", "", nil)

	String, _     = abi.NewType("string", "", nil)
	Bool, _       = abi.NewType("bool", "", nil)
	Bytes, _      = abi.NewType("bytes", "", nil)
	Bytes32, _    = abi.NewType("bytes32", "", nil)
	Address, _    = abi.NewType("address", "", nil)
	Uint64Arr, _  = abi.NewType("uint64[]", "", nil)
	AddressArr, _ = abi.NewType("address[]", "", nil)
	Int8, _       = abi.NewType("int8", "", nil)
	// Special types for testing
	Uint32Arr2, _       = abi.NewType("uint32[2]", "", nil)
	Uint64Arr2, _       = abi.NewType("uint64[2]", "", nil)
	Uint256Arr, _       = abi.NewType("uint256[]", "", nil)
	Uint256Arr2, _      = abi.NewType("uint256[2]", "", nil)
	Uint256Arr3, _      = abi.NewType("uint256[3]", "", nil)
	Uint256ArrNested, _ = abi.NewType("uint256[2][2]", "", nil)
	Uint8ArrNested, _   = abi.NewType("uint8[][2]", "", nil)
	Uint8SliceNested, _ = abi.NewType("uint8[][]", "", nil)
	KvPair, _           = abi.NewType("tuple", "struct Overloader.F", []abi.ArgumentMarshaling{
		{Name: "key", Type: "bytes"},
		{Name: "value", Type: "bytes"},
	})
	KvPairArr, _ = abi.NewType("tuple[]", "struct Overloader.F", []abi.ArgumentMarshaling{
		{Name: "key", Type: "string"},
		{Name: "value", Type: "bytes"},
	})
	AspectBoundInfoArr, _ = abi.NewType("tuple[]", "struct Overloader.F", []abi.ArgumentMarshaling{
		{Name: "aspectId", Type: "address"},
		{Name: "version", Type: "uint64"},
		{Name: "priority", Type: "int8"},
	})
)

var methods = map[string]abi.Method{
	"deploy": abi.NewMethod("deploy", "deploy", abi.Function, "", false, false, []abi.Argument{
		{Name: "code", Type: Bytes, Indexed: false},
		{Name: "initdata", Type: Bytes, Indexed: false},
		{Name: "properties", Type: KvPairArr, Indexed: false},
		{Name: "account", Type: Address, Indexed: false},
		{Name: "proof", Type: Bytes, Indexed: false},
		{Name: "joinPoints", Type: Uint256, Indexed: false},
	}, nil),
	"upgrade": abi.NewMethod("upgrade", "upgrade", abi.Function, "", false, false, []abi.Argument{
		{Name: "aspectId", Type: Address, Indexed: false},
		{Name: "code", Type: Bytes, Indexed: false},
		{Name: "properties", Type: KvPairArr, Indexed: false},
		{Name: "joinPoints", Type: Uint256, Indexed: false},
	}, nil),
	"bind": abi.NewMethod("bind", "bind", abi.Function, "", false, false, []abi.Argument{
		{Name: "aspectId", Type: Address, Indexed: false},
		{Name: "aspectVersion", Type: Uint256, Indexed: false},
		{Name: "contract", Type: Address, Indexed: false},
		{Name: "priority", Type: Int8, Indexed: false},
	}, nil),
	"unbind": abi.NewMethod("unbind", "unbind", abi.Function, "", false, false, []abi.Argument{
		{Name: "aspectId", Type: Address, Indexed: false},
		{Name: "contract", Type: Address, Indexed: false},
	}, nil),
	"changeVersion": abi.NewMethod("changeVersion", "changeVersion", abi.Function, "", false, false, []abi.Argument{
		{Name: "aspectId", Type: Address, Indexed: false},
		{Name: "contract", Type: Address, Indexed: false},
		{Name: "version", Type: Uint64, Indexed: false},
	}, nil),
	"versionOf": abi.NewMethod("versionOf", "versionOf", abi.Function, "", false, false, []abi.Argument{
		{Name: "aspectId", Type: Address, Indexed: false},
	}, []abi.Argument{
		{Name: "version", Type: Uint64, Indexed: false},
	}),
	"aspectsOf": abi.NewMethod("aspectsOf", "aspectsOf", abi.Function, "", false, false, []abi.Argument{
		{Name: "contract", Type: Address, Indexed: false},
	}, []abi.Argument{
		{Name: "aspectBoundInfo", Type: AspectBoundInfoArr, Indexed: false},
	}),
	"boundAddressesOf": abi.NewMethod("boundAddressesOf", "boundAddressesOf", abi.Function, "", false, false, []abi.Argument{
		{Name: "aspectId", Type: Address, Indexed: false},
	}, []abi.Argument{
		{Name: "account", Type: AddressArr, Indexed: false},
	}),
	"entrypoint": abi.NewMethod("entrypoint", "entrypoint", abi.Function, "", false, false, []abi.Argument{
		{Name: "aspectId", Type: Address, Indexed: false},
		{Name: "optArgs", Type: Bytes, Indexed: false},
	}, []abi.Argument{
		{Name: "resultMap", Type: Bytes, Indexed: false},
	}),
}

var methodsLookup = AbiMap()

var AbiMap = func() map[string]string {
	abiIndex := make(map[string]string)
	for name, expM := range methods {
		abiIndex[hex.EncodeToString(expM.ID)] = name
	}
	return abiIndex
}

func GetMethodName(callData []byte) (string, error) {
	methodID := hex.EncodeToString(callData[:4])
	methodName, ok := methodsLookup[methodID]
	if !ok {
		return "", errorsmod.Wrapf(errortypes.ErrInvalidRequest, "method with id %s not found", methodID)
	}
	return methodName, nil
}

func ParseMethod(callData []byte) (*abi.Method, map[string]interface{}, error) {
	methodName, err := GetMethodName(callData)
	if err != nil {
		return nil, nil, err
	}

	method, ok := methods[methodName]
	if !ok {
		return nil, nil, errorsmod.Wrapf(errortypes.ErrInvalidRequest, "method %s is not valid", methodName)
	}

	argsMap := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(argsMap, callData[4:]); err != nil {
		return nil, nil, err
	}

	return &method, argsMap, nil
}

func ParseInput(tx []byte) (*abi.Method, map[string]interface{}, error) {
	methodId := hex.EncodeToString(tx[:4])
	methodName, exist := AbiMap()[methodId]
	if !exist {
		return nil, nil, errorsmod.Wrapf(errortypes.ErrInvalidRequest, "method with id %s not found", methodId)
	}
	method := methods[methodName]
	argsMap := make(map[string]interface{})
	inputs := method.Inputs
	err := inputs.UnpackIntoMap(argsMap, tx[4:])
	if err != nil {
		return nil, nil, err
	}

	return &method, argsMap, nil
}

func IsAspectDeploy(to *common.Address, callData []byte) bool {
	if !types.IsAspectContractAddr(to) {
		return false
	}

	methodName, err := GetMethodName(callData)
	return err == nil && methodName == "deploy"
}
