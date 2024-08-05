package types

import "github.com/ethereum/go-ethereum/common"

// Binding is the data model for holding the binding of an aspect to an account
type Binding struct {
	Account  common.Address
	Version  uint64
	Priority int8
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
