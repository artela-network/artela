package process

import (
	"github.com/artela-network/artela/x/evm/process/support"
	"github.com/ethereum/go-ethereum/common"
	ethereum "github.com/ethereum/go-ethereum/core/types"
)

// AccessList is an EIP-2930 access list that represents the slice of
// the protobuf AccessTuples.
type AccessList []support.AccessTuple

// NewAccessList creates a new protobuf-compatible AccessList from an ethereum
// core AccessList type
func NewAccessList(ethAccessList *ethereum.AccessList) AccessList {
	if ethAccessList == nil {
		return nil
	}

	al := AccessList{}
	for _, tuple := range *ethAccessList {
		storageKeys := make([]string, len(tuple.StorageKeys))

		for i := range tuple.StorageKeys {
			storageKeys[i] = tuple.StorageKeys[i].String()
		}

		al = append(al, support.AccessTuple{
			Address:     tuple.Address.String(),
			StorageKeys: storageKeys,
		})
	}

	return al
}

// ToEthAccessList is a utility function to convert the protobuf compatible
// AccessList to eth core AccessList from go-ethereum
func (al AccessList) ToEthAccessList() *ethereum.AccessList {
	var ethAccessList ethereum.AccessList

	for _, tuple := range al {
		storageKeys := make([]common.Hash, len(tuple.StorageKeys))

		for i := range tuple.StorageKeys {
			storageKeys[i] = common.HexToHash(tuple.StorageKeys[i])
		}

		ethAccessList = append(ethAccessList, ethereum.AccessTuple{
			Address:     common.HexToAddress(tuple.Address),
			StorageKeys: storageKeys,
		})
	}

	return &ethAccessList
}
