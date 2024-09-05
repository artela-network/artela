package types

import "math"

const (
	// ModuleName string name of module
	ModuleName = "aspect"

	// StoreKey key for aspect storage data
	StoreKey = ModuleName
)

const (
	// AspectPropertyLimit is the maximum number of properties that can be set on an aspect
	AspectPropertyLimit = math.MaxUint8
)
