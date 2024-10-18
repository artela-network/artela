package precompiled

import (
	"github.com/artela-network/artela-evm/vm"
	"github.com/ethereum/go-ethereum/common"

	"github.com/artela-network/artela/x/evm/precompile/erc20"
)

var Contracts = map[common.Address]vm.PrecompiledContract{
	common.HexToAddress("0x0000000000000000000000000000000000001234"): erc20.GlobalERC20Contract,
}

func RegisterPrecompiles() {
	for k, v := range Contracts {
		vm.PrecompiledContractsHomestead[k] = v
	}
	for k, v := range Contracts {
		vm.PrecompiledContractsByzantium[k] = v
	}
	for k, v := range Contracts {
		vm.PrecompiledContractsIstanbul[k] = v
	}
	for k, v := range Contracts {
		vm.PrecompiledContractsBerlin[k] = v
	}

	for k := range Contracts {
		vm.PrecompiledAddressesHomestead = append(vm.PrecompiledAddressesHomestead, k)
	}
	for k := range Contracts {
		vm.PrecompiledAddressesByzantium = append(vm.PrecompiledAddressesByzantium, k)
	}
	for k := range Contracts {
		vm.PrecompiledAddressesIstanbul = append(vm.PrecompiledAddressesIstanbul, k)
	}
	for k := range Contracts {
		vm.PrecompiledAddressesBerlin = append(vm.PrecompiledAddressesBerlin, k)
	}
}
