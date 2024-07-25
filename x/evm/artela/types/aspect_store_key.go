package types

import (
	"bytes"
	"encoding/binary"
	"math"
	"strings"

	"github.com/emirpasic/gods/utils"
	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/common"

	artela "github.com/artela-network/aspect-core/types"
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
	AspectStateKeyPrefix       = "AspectStore/State/"

	AspectJoinPointRunKeyPrefix = "AspectStore/JoinPointRun/"

	AspectIDMapKey = "aspectId"
	VersionMapKey  = "version"
	PriorityMapKey = "priority"

	AspectAccountKey           = "Aspect_@Acount@_"
	AspectProofKey             = "Aspect_@Proof@_"
	AspectRunJoinPointKey      = "Aspect_@Run@JoinPoint@_"
	AspectPropertyAllKeyPrefix = "Aspect_@Property@AllKey@_"
	AspectPropertyAllKeySplit  = "^^^"
	AspectPropertyLimit        = math.MaxUint8
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
	aspectID []byte,
	propertyKey []byte,
) []byte {
	key := make([]byte, 0, len(aspectID)+PathSeparatorLen*2+len(propertyKey))

	key = append(key, aspectID...)
	key = append(key, PathSeparator...)
	key = append(key, propertyKey...)
	key = append(key, PathSeparator...)

	return key
}

func AspectVersionKey(
	aspectID []byte,
	version []byte,
) []byte {
	key := make([]byte, 0, len(aspectID)+PathSeparatorLen*2+len(version))

	key = append(key, aspectID...)
	key = append(key, PathSeparator...)
	key = append(key, version...)
	key = append(key, PathSeparator...)

	return key
}

func AspectIDKey(
	aspectID []byte,
) []byte {
	key := make([]byte, 0, len(aspectID)+PathSeparatorLen)
	key = append(key, aspectID...)
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

type AspectInfo struct {
	AspectId common.Address `json:"AspectId"`
	Version  uint64         `json:"Version"`
	Priority int8           `json:"Priority"`
}

type AspectMeta struct {
	Id       common.Address `json:"id"`
	Version  *uint256.Int   `json:"version"`
	Priority int64          `json:"priority"`
}
type Property struct {
	Key   string `json:"Key"`
	Value []byte `json:"Value"`
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
