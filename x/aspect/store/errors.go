package store

import "errors"

var (
	ErrInvalidProtocolInfo     = errors.New("invalid protocol info")
	ErrUnknownProtocolVersion  = errors.New("unknown protocol version")
	ErrCodeNotFound            = errors.New("code not found")
	ErrInvalidStorageKey       = errors.New("invalid storage key")
	ErrTooManyProperties       = errors.New("aspect property limit exceeds")
	ErrInvalidBinding          = errors.New("invalid binding")
	ErrStorageCorrupted        = errors.New("storage corrupted")
	ErrNoJoinPoint             = errors.New("cannot bind with no-joinpoint aspect")
	ErrBindingLimitExceeded    = errors.New("binding limit exceeded")
	ErrAlreadyBound            = errors.New("aspect already bound")
	ErrNotBound                = errors.New("aspect not bound")
	ErrInvalidStoreContext     = errors.New("invalid store context")
	ErrPropertyReserved        = errors.New("property key reserved")
	ErrInvalidExtension        = errors.New("invalid extension")
	ErrInvalidVersionMeta      = errors.New("invalid version meta")
	ErrSerdeFail               = errors.New("serialize or deserialize fail")
	ErrBoundNonVerifierWithEOA = errors.New("binding non-verifier aspect with EOA")
	ErrInvalidJoinPoint        = errors.New("invalid join point")
)
