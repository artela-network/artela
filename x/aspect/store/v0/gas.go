package v0

import (
	"github.com/artela-network/artela/x/aspect/types"
)

const (
	storageLoadCost     = 10
	storageStoreCost    = 1000
	storageSaveCodeCost = 1000
	storageUpdateCost   = 1000
)

type GasMeter interface {
	MeasureStorageUpdate(dataLen int) error
	MeasureStorageCodeSave(dataLen int) error
	MeasureStorageStore(dataLen int) error
	MeasureStorageLoad(dataLen int) error
	RemainingGas() uint64
	UpdateGas(newGas uint64)
	Consume(dataLen int, gasCostPer32Bytes uint64) error
}

type noopGasMeter struct {
	ctx types.StoreContext
}

func NewNoOpGasMeter(ctx types.StoreContext) GasMeter {
	return &noopGasMeter{
		ctx: ctx,
	}
}

func (n *noopGasMeter) UpdateGas(_ uint64) {
	return
}

func (n *noopGasMeter) MeasureStorageUpdate(_ int) error {
	return nil
}

func (n *noopGasMeter) MeasureStorageCodeSave(_ int) error {
	return nil
}

func (n *noopGasMeter) MeasureStorageStore(_ int) error {
	return nil
}

func (n *noopGasMeter) MeasureStorageLoad(_ int) error {
	return nil
}

func (n *noopGasMeter) RemainingGas() uint64 {
	return n.ctx.Gas()
}

func (n *noopGasMeter) Consume(_ int, _ uint64) error {
	return nil
}

// gasMeter is a simple gas metering implementation for aspect storage v0.
type gasMeter struct {
	ctx types.StoreContext
}

func NewGasMeter(ctx types.StoreContext) GasMeter {
	return &gasMeter{
		ctx: ctx,
	}
}

func (m *gasMeter) MeasureStorageUpdate(dataLen int) error {
	return m.Consume(dataLen, storageUpdateCost)
}

func (m *gasMeter) MeasureStorageCodeSave(dataLen int) error {
	return m.Consume(dataLen, storageSaveCodeCost)
}

func (m *gasMeter) MeasureStorageStore(dataLen int) error {
	return m.Consume(dataLen, storageStoreCost)
}

func (m *gasMeter) MeasureStorageLoad(dataLen int) error {
	return m.Consume(dataLen, storageLoadCost)
}

func (m *gasMeter) RemainingGas() uint64 {
	return m.ctx.Gas()
}

func (m *gasMeter) Consume(dataLen int, gasCostPer32Bytes uint64) error {
	gas := ((uint64(dataLen) + 32) >> 5) * gasCostPer32Bytes
	return m.ctx.ConsumeGas(gas)
}

func (m *gasMeter) UpdateGas(gas uint64) {
	m.ctx.UpdateGas(gas)
}
