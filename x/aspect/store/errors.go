package store

import "errors"

var (
	ErrInvalidProtocolInfo    = errors.New("invalid protocol info")
	ErrUnknownProtocolVersion = errors.New("unknown protocol version")
	ErrCodeNotFound           = errors.New("code not found")
	ErrAspectNotFound         = errors.New("aspect not found")
	ErrInvalidStorageKey      = errors.New("invalid storage key")
	ErrTooManyProperties      = errors.New("aspect property limit exceeds")
	ErrInvalidBinding         = errors.New("invalid binding")
	ErrStorageCorrupted       = errors.New("storage corrupted")
	ErrBindingLimitExceeded   = errors.New("binding limit exceeded")
	ErrAlreadyBound           = errors.New("aspect already bound")
	ErrInvalidStoreContext    = errors.New("invalid store context")
	ErrPropertyReserved       = errors.New("property key reserved")
)
