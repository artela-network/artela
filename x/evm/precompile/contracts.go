package precompiled

import (
	"github.com/artela-network/artela-evm/vm"
	"github.com/ethereum/go-ethereum/common"
)

func RegisterPrecompiles(address common.Address, p vm.PrecompiledContract) {
	vm.PrecompiledContractsHomestead[address] = p
	vm.PrecompiledContractsByzantium[address] = p
	vm.PrecompiledContractsIstanbul[address] = p
	vm.PrecompiledContractsBerlin[address] = p

	vm.PrecompiledAddressesHomestead = append(vm.PrecompiledAddressesHomestead, address)
	vm.PrecompiledAddressesByzantium = append(vm.PrecompiledAddressesByzantium, address)
	vm.PrecompiledAddressesIstanbul = append(vm.PrecompiledAddressesIstanbul, address)
	vm.PrecompiledAddressesBerlin = append(vm.PrecompiledAddressesBerlin, address)
}
