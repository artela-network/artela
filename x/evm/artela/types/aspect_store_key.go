package types

import (
	"bytes"
	"encoding/binary"
	"strings"

	artela "github.com/artela-network/aspect-core/types"
	"github.com/emirpasic/gods/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var _ binary.ByteOrder

const (
	// AspectCodeKeyPrefix is the prefix to retrieve all AspectCodeStore
	AspectCodeKeyPrefix        = "AspectStore/Code/"
	AspectCodeVersionKeyPrefix = "AspectStore/Version/"
	AspectPropertyKeyPrefix    = "AspectStore/Property/"
	ContractBindKeyPrefix      = "AspectStore/ContractBind/"
	VerifierBindingKeyPrefix   = "AspectStore/VerifierBind/"
	AspectRefKeyPrefix         = "AspectStore/AspectRef/"
	AspectBlockKeyPrefix       = "AspectStore/Block/"
	AspectStateKeyPrefix       = "AspectStore/State/"

	AspectIdMapKey = "aspectId"
	VersionMapKey  = "version"
	PriorityMapKey = "priority"

	AspectAccountKey = "Aspect_@Acount@_"
	AspectProofKey   = "Aspect_@Proof@_"
)

var (
	PathSeparator    = []byte("/")
	PathSeparatorLen = len(PathSeparator)
)

func AspectArrayKey(keys ...[]byte) []byte {
	var key []byte
	for _, b := range keys {
		key = append(key, b...)
		key = append(key, PathSeparator...)
	}
	return key
}

// AspectCodeStoreKey returns the store key to retrieve a AspectCodeStore from the index fields
func AspectPropertyKey(
	aspectId []byte,
	propertyKey []byte,
) []byte {
	key := make([]byte, 0, len(aspectId)+PathSeparatorLen*2+len(propertyKey))

	key = append(key, aspectId...)
	key = append(key, PathSeparator...)
	key = append(key, propertyKey...)
	key = append(key, PathSeparator...)

	return key
}

func AspectVersionKey(
	aspectId []byte,
	version []byte,
) []byte {
	key := make([]byte, 0, len(aspectId)+PathSeparatorLen*2+len(version))

	key = append(key, aspectId...)
	key = append(key, PathSeparator...)
	key = append(key, version...)
	key = append(key, PathSeparator...)

	return key
}

func AspectIdKey(
	aspectId []byte,
) []byte {
	key := make([]byte, 0, len(aspectId)+PathSeparatorLen)
	key = append(key, aspectId...)
	key = append(key, PathSeparator...)

	return key
}

func AspectBlockKey() []byte {
	var key []byte
	key = append(key, []byte("AspectBlock")...)
	key = append(key, PathSeparator...)
	return key
}

func AccountKey(
	account []byte,
) []byte {
	key := make([]byte, 0, len(account)+PathSeparatorLen)
	key = append(key, account...)
	key = append(key, PathSeparator...)
	return key
}

type AspectMeta struct {
	Id       common.Address `json:"id"`
	Version  *uint256.Int   `json:"version"`
	Priority int64          `json:"priority"`
}
type Property struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type BoundAspectCode struct {
	AspectId common.Address `json:"aspectId"`
	Version  *uint256.Int   `json:"version"`
	Priority int64          `json:"priority"`
	Code     []byte         `json:"code"`
}

func ByMapKeyPriority(a, b interface{}) int {
	priorityA, ok := a.(map[string]interface{})[PriorityMapKey]
	if !ok {
		priorityA = 0
	}
	priorityB, okb := b.(map[string]interface{})[PriorityMapKey]
	if !okb {
		priorityB = 1
	}
	return utils.IntComparator(priorityA, priorityB) // "-" descending order
}

func NewBindingPriorityComparator(x []*AspectMeta) func(i, j int) bool {
	return func(i, j int) bool {
		return x[i].Priority < x[j].Priority && (bytes.Compare(x[i].Id.Bytes(), x[j].Id.Bytes()) < 0)
	}
}

func NewBindingAspectPriorityComparator(x []*artela.AspectCode) func(i, j int) bool {
	return func(i, j int) bool {
		return (x[i].Priority < x[j].Priority) && (strings.Compare(x[i].AspectId, x[j].AspectId) < 0)
	}
}
