package v1

import "github.com/artela-network/artela/x/aspect/store"

func init() {
	store.RegisterAccountStore(1, NewAccountStore)
	store.RegisterAspectMetaStore(1, NewAspectMetaStore)
	store.RegisterAspectStateStore(1, NewStateStore)
}
