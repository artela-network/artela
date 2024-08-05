package v0

import (
	"github.com/artela-network/artela/x/aspect/types"
)

const (
	storageLoadCost     = 50
	storageStoreCost    = 20000
	storageSaveCodeCost = 1000
	storageUpdateCost   = 5000
)

type GasMeter interface {
	measureStorageUpdate(dataLen int) error
	measureStorageCodeSave(dataLen int) error
	measureStorageStore(dataLen int) error
	measureStorageLoad(dataLen int) error
	remainingGas() uint64
	consume(dataLen int, gasCostPer32Bytes uint64) error
}

type noopGasMeter struct {
	ctx types.StoreContext
}

func newNoOpGasMeter(ctx types.StoreContext) GasMeter {
	return &noopGasMeter{
		ctx: ctx,
	}
}

func (n *noopGasMeter) measureStorageUpdate(_ int) error {
	return nil
}

func (n *noopGasMeter) measureStorageCodeSave(_ int) error {
	return nil
}

func (n *noopGasMeter) measureStorageStore(_ int) error {
	return nil
}

func (n *noopGasMeter) measureStorageLoad(_ int) error {
	return nil
}

func (n *noopGasMeter) remainingGas() uint64 {
	return n.ctx.Gas()
}

func (n *noopGasMeter) consume(_ int, _ uint64) error {
	return nil
}

// gasMeter is a simple gas metering implementation for aspect storage v0.
type gasMeter struct {
	ctx types.StoreContext
}

func newGasMeter(ctx types.StoreContext) GasMeter {
	return &gasMeter{
		ctx: ctx,
	}
}

func (m *gasMeter) measureStorageUpdate(dataLen int) error {
	return m.consume(dataLen, storageUpdateCost)
}

func (m *gasMeter) measureStorageCodeSave(dataLen int) error {
	return m.consume(dataLen, storageSaveCodeCost)
}

func (m *gasMeter) measureStorageStore(dataLen int) error {
	return m.consume(dataLen, storageStoreCost)
}

func (m *gasMeter) measureStorageLoad(dataLen int) error {
	return m.consume(dataLen, storageLoadCost)
}

func (m *gasMeter) remainingGas() uint64 {
	return m.ctx.Gas()
}

func (m *gasMeter) consume(dataLen int, gasCostPer32Bytes uint64) error {
	gas := ((uint64(dataLen) + 32) >> 5) * gasCostPer32Bytes
	return m.ctx.ConsumeGas(gas)
}
