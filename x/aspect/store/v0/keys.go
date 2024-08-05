package v0

// v0 keys
const (
	// AspectCodeKeyPrefix is the prefix to retrieve all AspectCodeStore
	V0AspectCodeKeyPrefix        = "AspectStore/Code/"
	V0AspectCodeVersionKeyPrefix = "AspectStore/Version/"
	V0AspectPropertyKeyPrefix    = "AspectStore/Property/"
	V0ContractBindKeyPrefix      = "AspectStore/ContractBind/"
	V0VerifierBindingKeyPrefix   = "AspectStore/VerifierBind/"
	V0AspectRefKeyPrefix         = "AspectStore/AspectRef/"
	V0AspectStateKeyPrefix       = "AspectStore/State/"

	V0AspectJoinPointRunKeyPrefix = "AspectStore/JoinPointRun/"

	V0AspectIDMapKey = "aspectId"
	V0VersionMapKey  = "version"
	V0PriorityMapKey = "priority"

	V0AspectAccountKey           = "Aspect_@Acount@_"
	V0AspectProofKey             = "Aspect_@Proof@_"
	V0AspectRunJoinPointKey      = "Aspect_@Run@JoinPoint@_"
	V0AspectPropertyAllKeyPrefix = "Aspect_@Property@AllKey@_"
	V0AspectPropertyAllKeySplit  = "^^^"
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
