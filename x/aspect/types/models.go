package types

import (
	"github.com/artela-network/aspect-core/types"
	"github.com/ethereum/go-ethereum/common"
)

// BindingFilter is the data model for holding the filter of querying aspect bindings
type BindingFilter struct {
	// filter bindings with given join point
	JoinPoint *types.PointCut
	// only return bindings of verifier aspects
	VerifierOnly bool
	// only return bindings of tx level aspects
	TxLevelOnly bool
}

func NewDefaultFilter(isCA bool) BindingFilter {
	if isCA {
		return BindingFilter{}
	}

	return BindingFilter{
		VerifierOnly: true,
	}
}

func NewJoinPointFilter(cut types.PointCut) BindingFilter {
	return BindingFilter{
		JoinPoint: &cut,
	}
}

// Binding is the data model for holding the binding of an aspect to an account
type Binding struct {
	Account   common.Address
	Version   uint64
	Priority  int8
	JoinPoint uint16
}

// VersionMeta is the data model for holding the metadata of a specific version of aspect
type VersionMeta struct {
	JoinPoint uint64
	CodeHash  common.Hash
}

// AspectMeta is the data model for holding the metadata of an aspect
type AspectMeta struct {
	PayMaster common.Address
	Proof     []byte
}

// Property is the data model for holding the properties of an aspect
type Property struct {
	Key   string `json:"Key"`
	Value []byte `json:"Value"`
}
