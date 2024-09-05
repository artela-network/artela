package v0

import "github.com/artela-network/artela/x/aspect/store"

func init() {
	store.RegisterAccountStore(0, NewAccountStore)
	store.RegisterAspectMetaStore(0, NewAspectMetaStore)
	store.RegisterAspectStateStore(0, NewStateStore)
}
