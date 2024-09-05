package rpc

import (
	"fmt"
	"runtime"

	"github.com/cosmos/cosmos-sdk/version"
)

// ClientVersion returns the current client version.
func (b *BackendImpl) ClientVersion() string {
	return fmt.Sprintf(
		"%s/%s/%s/%s",
		version.Name,
		version.Version,
		runtime.GOOS+"-"+runtime.GOARCH,
		runtime.Version(),
	)
}
