package txs

// EIP-2930 was part of the Berlin upgrade and is seen as a step towards improving the
// efficiency and flexibility of the Ethereum network. It allows for more sophisticated
// transaction handling and can mitigate some of the effects of increased gas costs
// introduced by other changes in the protocol.
//
// Access lists in EIP-2930 provide several benefits:
// 1. **Gas Cost Predictability**: By specifying the addresses and keys that will be accessed,
// the sender can pre-calculate the gas costs, making them more predictable.
// 2. **Increased Efficiency**: Clients and miners can txs transactions more efficiently,
// as they know in advance which addresses and keys will be accessed.
// 3. **Compatibility with Future Upgrades**: Access lists can facilitate smoother upgrades,
// such as the transition to Ethereum 2.0, by allowing transactions to explicitly states their dependencies.

import (
	"github.com/ethereum/go-ethereum/common"
	ethereum "github.com/ethereum/go-ethereum/core/types"

	"github.com/artela-network/artela/x/evm/txs/support"
)

// AccessList is EIP-2930 access list
type AccessList []support.AccessTuple

// NewAccessList creates a new protobuf-compatible AccessList from an ethereum type
func NewAccessList(ethAccessList *ethereum.AccessList) AccessList {
	if ethAccessList == nil {
		return nil
	}

	accessList := AccessList{}
	for _, tuple := range *ethAccessList {
		storageKeys := make([]string, len(tuple.StorageKeys))

		for i := range tuple.StorageKeys {
			storageKeys[i] = tuple.StorageKeys[i].String()
		}

		accessList = append(accessList, support.AccessTuple{
			Address:     tuple.Address.String(),
			StorageKeys: storageKeys,
		})
	}

	return accessList
}

// ToEthAccessList convert the protobuf-compatible AccessList to an ethereum AccessList
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
