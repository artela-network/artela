package v1

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/status-im/keycard-go/hexutils"
)

type Filter interface {
	Lookup(data []byte) bool
	Insert(data []byte) bool
	Delete(data []byte) bool
	Count() uint
	Encode() []byte
}

type LoggedFilter struct {
	inner  Filter
	logger log.Logger
}

func NewLoggedFilter(logger log.Logger, filter Filter) Filter {
	return &LoggedFilter{
		inner:  filter,
		logger: logger,
	}
}

func (lf *LoggedFilter) Lookup(data []byte) bool {
	res := lf.inner.Lookup(data)
	lf.logger.Debug("lookup in filter", "data", hexutils.BytesToHex(data), "result", res)
	return res
}

func (lf *LoggedFilter) Insert(data []byte) bool {
	res := lf.inner.Insert(data)
	lf.logger.Debug("insert into filter", "data", hexutils.BytesToHex(data), "result", res)
	return res
}

func (lf *LoggedFilter) Delete(data []byte) bool {
	res := lf.inner.Delete(data)
	lf.logger.Debug("delete from filter", "data", hexutils.BytesToHex(data), "result", res)
	return res
}

func (lf *LoggedFilter) Count() uint {
	res := lf.inner.Count()
	lf.logger.Debug("count in filter", "result", res)
	return res
}

func (lf *LoggedFilter) Encode() []byte {
	res := lf.inner.Encode()
	lf.logger.Debug("encode filter", "result", hexutils.BytesToHex(res))
	return res
}
