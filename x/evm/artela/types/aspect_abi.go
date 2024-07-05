package types

import (
	"encoding/hex"

	errorsmod "cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

// nolint
var methods = map[string]abi.Method{
	"deploy":           abi.NewMethod("deploy", "deploy", abi.Function, "", false, false, []abi.Argument{{"code", Bytes, false}, {"initdata", Bytes, false}, {"properties", KvPairArr, false}, {"account", Address, false}, {"proof", Bytes, false}, {"joinPoints", Uint256, false}}, nil),
	"upgrade":          abi.NewMethod("upgrade", "upgrade", abi.Function, "", false, false, []abi.Argument{{"aspectId", Address, false}, {"code", Bytes, false}, {"properties", KvPairArr, false}, {"joinPoints", Uint256, false}}, nil),
	"bind":             abi.NewMethod("bind", "bind", abi.Function, "", false, false, []abi.Argument{{"aspectId", Address, false}, {"aspectVersion", Uint256, false}, {"contract", Address, false}, {"priority", Int8, false}}, nil),
	"unbind":           abi.NewMethod("unbind", "unbind", abi.Function, "", false, false, []abi.Argument{{"aspectId", Address, false}, {"contract", Address, false}}, nil),
	"changeVersion":    abi.NewMethod("changeVersion", "changeVersion", abi.Function, "", false, false, []abi.Argument{{"aspectId", Address, false}, {"contract", Address, false}, {"version", Uint64, false}}, nil),
	"versionOf":        abi.NewMethod("versionOf", "versionOf", abi.Function, "", false, false, []abi.Argument{{"aspectId", Address, false}}, []abi.Argument{{"version", Uint64, false}}),
	"aspectsOf":        abi.NewMethod("aspectsOf", "aspectsOf", abi.Function, "", false, false, []abi.Argument{{"contract", Address, false}}, []abi.Argument{{"aspectBoundInfo", AspectBoundInfoArr, false}}),
	"boundAddressesOf": abi.NewMethod("boundAddressesOf", "boundAddressesOf", abi.Function, "", false, false, []abi.Argument{{"aspectId", Address, false}}, []abi.Argument{{"account", AddressArr, false}}),
	"entrypoint":       abi.NewMethod("entrypoint", "entrypoint", abi.Function, "", false, false, []abi.Argument{{"aspectId", Address, false}, {"optArgs", Bytes, false}}, []abi.Argument{{"resultMap", Bytes, false}}),
}

// nolint
var AspectOwnableMethod = map[string]abi.Method{
	"isOwner": abi.NewMethod("isOwner", "isOwner", abi.Function, "", false, false, []abi.Argument{{"sender", Address, false}}, []abi.Argument{{"result", Bool, false}}),
}

var methodsLookup = AbiMap()

var AbiMap = func() map[string]string {
	abiIndex := make(map[string]string)
	for name, expM := range methods {
		abiIndex[hex.EncodeToString(expM.ID)] = name
	}
	return abiIndex
}

func ParseMethod(callData []byte) (*abi.Method, map[string]interface{}, error) {
	methodId := hex.EncodeToString(callData[:4])
	methodName, ok := methodsLookup[methodId]
	if !ok {
		return nil, nil, errorsmod.Wrapf(nil, "missing expected method %s", methodId)
	}

	method, ok := methods[methodName]
	if !ok {
		return nil, nil, errorsmod.Wrapf(nil, "method %s does not exist", methodName)
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
		return nil, nil, errorsmod.Wrapf(nil, "Missing expected method %v", methodId)
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
