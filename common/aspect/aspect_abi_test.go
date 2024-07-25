package aspect

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/holiman/uint256"
	jsoniter "github.com/json-iterator/go"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/artela-network/artela/x/evm/artela/types"
)

func TestAbi(t *testing.T) {
	bytes, errs := rlp.EncodeToBytes("owner")
	if errs != nil {
		return
	}
	var owner string
	err := rlp.DecodeBytes(bytes, &owner)
	if err != nil {
		return
	}

	fmt.Println(owner)

	exp := abi.ABI{
		Methods: methods,
	}
	marshal, _ := jsoniter.Marshal(exp)
	fmt.Println(string(marshal))

	abiIndex := make(map[string]abi.Method)
	for _, expM := range exp.Methods {
		abiIndex[hex.EncodeToString(expM.ID)] = expM
	}

	for name, expM := range exp.Methods {
		gotM, exist := exp.Methods[name]
		if !exist {
			t.Errorf("Missing expected method %v", name)
		}
		if !reflect.DeepEqual(gotM, expM) {
			t.Errorf("\nGot abi method: \n%v\ndoes not match expected method\n%v", gotM, expM)
		}
		sig := gotM.Sig
		keccak256 := crypto.Keccak256([]byte(gotM.Sig))
		fmt.Println(sig, hex.EncodeToString(gotM.ID), hex.EncodeToString(keccak256[:4]))
	}
}

type Dummy struct {
	Key   []byte `json:"Key"`
	Value []byte `json:"Value"`
}

func TestPack(t *testing.T) {
	exp := abi.ABI{
		Methods: methods,
	}

	code, _ := hex.DecodeString("324234132131")

	ke1, _ := hex.DecodeString("444444444444")
	ke2, _ := hex.DecodeString("111111111111")
	value1, _ := hex.DecodeString("2222222222222")
	value2, _ := hex.DecodeString("3333333333333")

	arrin := []struct {
		Key   []byte
		Value []byte
	}{
		{ke1, value1},
		{ke2, value2},
	}
	fixedArrStrPack, _ := exp.Pack("Deploy", code, arrin)
	fmt.Println(hex.EncodeToString(fixedArrStrPack))
	out := make(map[string]interface{})
	// err := exp.UnpackIntoMap(out, "Deploy", fixedArrStrPack)
	err := exp.Methods["Deploy"].Inputs.UnpackIntoMap(out, fixedArrStrPack[4:])
	if err != nil {
		fmt.Println(err)
		return
	}

	properties := out["properties"].([]struct {
		Key   []byte `json:"Key"`
		Value []byte `json:"Value"`
	})
	for i := range properties {
		s := properties[i]
		fmt.Println(hex.EncodeToString(s.Key), hex.EncodeToString(s.Value))
	}
	fmt.Println(properties)
}

func TestAspectsOfPack(t *testing.T) {
	ma := make(map[string]interface{}, 2)
	ma[types.AspectIDMapKey] = "1111"
	ma[types.VersionMapKey] = 0
	ma[types.PriorityMapKey] = 1
	mb := make(map[string]interface{}, 2)
	mb[types.AspectIDMapKey] = "2222"
	mb[types.VersionMapKey] = 1
	mb[types.PriorityMapKey] = -10
	mc := make(map[string]interface{}, 2)
	mc[types.AspectIDMapKey] = "3333"
	mc[types.VersionMapKey] = 2
	mc[types.PriorityMapKey] = 3

	queue := pq.NewWith(types.ByMapKeyPriority) // empty
	queue.Enqueue(ma)                           // {a 1}
	queue.Enqueue(mb)                           // {c 3}, {a 1}
	queue.Enqueue(mc)

	outputs := make([]types.AspectMeta, 0)
	iterator := queue.Iterator()

	for iterator.Next() {
		m := iterator.Value().(map[string]interface{})
		e := types.AspectMeta{
			Id:       common.HexToAddress(m[types.AspectIDMapKey].(string)),
			Priority: int64(m[types.PriorityMapKey].(int)),
			Version:  uint256.NewInt(m[types.VersionMapKey].(uint64)),
		}
		outputs = append(outputs, e)
	}
	pack, err := methods["AspectsOf"].Outputs.Pack(outputs)
	if err != nil {
		fmt.Println("pack error", err)
	}
	maps := make(map[string]interface{}, 0)
	err2 := methods["AspectsOf"].Outputs.UnpackIntoMap(maps, pack)
	if err2 != nil {
		fmt.Println("pack error", err2)
	}
	fmt.Println("unpack==", maps)
	aspects := maps["aspectBoundInfo"].([]struct {
		AspectId common.Address `json:"AspectId"`
		Version  uint64         `json:"Version"`
		Priority int8           `json:"Priority"`
	})
	fmt.Println(aspects)
}

func TestContractOfPack(t *testing.T) {
	treeset := treeset.NewWithStringComparator()
	treeset.Add("aaaaaaa")
	treeset.Add("bbbbbbbb")
	treeset.Add("ccccccc")

	addressAry := make([]common.Address, 0)
	iterator := treeset.Iterator()
	for iterator.Next() {
		contract := iterator.Value()
		if contract != nil {
			contractAddr := common.HexToAddress(contract.(string))
			addressAry = append(addressAry, contractAddr)
		}
	}
	ret, packErr := methods["ContractsOf"].Outputs.Pack(addressAry)
	if packErr != nil {
		fmt.Println("pack error", packErr)
	}

	maps := make(map[string]interface{}, 0)
	err2 := methods["ContractsOf"].Outputs.UnpackIntoMap(maps, ret)
	if err2 != nil {
		fmt.Println("unpack error", err2)
	}
	fmt.Println("unpack==", maps)
	aspects := maps["contracts"].([]common.Address)
	fmt.Println(aspects)
	// mock response
}
