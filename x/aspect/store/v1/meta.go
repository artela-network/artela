package v1

import (
	"encoding/json"
	"math"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	cuckoo "github.com/seiflotfy/cuckoofilter"

	"github.com/artela-network/artela/x/aspect/store"
	v0 "github.com/artela-network/artela/x/aspect/store/v0"
	"github.com/artela-network/artela/x/aspect/types"
)

var _ store.AspectMetaStore = (*metaStore)(nil)

const (
	bindingSlotSize    = 120
	bindingInfoLength  = 32
	bindingDataLength  = 8
	filterMaxSize      = 2047
	filterActualLimit  = filterMaxSize * 9 / 10 // assume load factor of filter is 90%
	filterManagedSlots = filterActualLimit / bindingSlotSize
	maxBindingSize     = filterManagedSlots * math.MaxUint8 * bindingSlotSize
)

type metaStore struct {
	BaseStore

	ext *Extension
	ctx *types.AspectStoreContext

	propertiesCache map[uint64]map[string][]byte
}

// NewAspectMetaStore creates a new instance of aspect meta Store.
func NewAspectMetaStore(ctx *types.AspectStoreContext, protocolExtension []byte) store.AspectMetaStore {
	var meter v0.GasMeter
	if ctx.ChargeGas() {
		meter = v0.NewGasMeter(ctx)
	} else {
		meter = v0.NewNoOpGasMeter(ctx)
	}

	ext := new(Extension)
	if err := ext.UnmarshalText(protocolExtension); err != nil {
		panic(err)
	}

	return &metaStore{
		BaseStore:       NewBaseStore(ctx.CosmosContext().Logger(), meter, ctx.CosmosContext().KVStore(ctx.AspectStoreKey())),
		ctx:             ctx,
		ext:             ext,
		propertiesCache: make(map[uint64]map[string][]byte),
	}
}

func (m *metaStore) GetCode(version uint64) ([]byte, error) {
	// key format {5B codePrefix}{8B version}{20B aspectID}
	key := store.NewKeyBuilder(V1AspectCodeKeyPrefix).
		AppendUint64(version).
		AppendBytes(m.ctx.AspectID.Bytes()).
		Build()
	return m.Load(key)
}

func (m *metaStore) GetVersionMeta(version uint64) (*types.VersionMeta, error) {
	// key format {5B metaPrefix}{8B version}{20B aspectID}
	key := store.NewKeyBuilder(V1AspectMetaKeyPrefix).
		AppendUint64(version).
		AppendBytes(m.ctx.AspectID.Bytes()).
		Build()
	raw, err := m.Load(key)
	if err != nil {
		return nil, err
	}

	meta := new(VersionMeta)
	if err := meta.UnmarshalText(raw); err != nil {
		return nil, err
	}

	return &types.VersionMeta{
		JoinPoint: meta.JoinPoint,
		CodeHash:  meta.CodeHash,
	}, nil
}

func (m *metaStore) GetMeta() (*types.AspectMeta, error) {
	return &types.AspectMeta{
		PayMaster: m.ext.PayMaster,
		Proof:     m.ext.Proof,
	}, nil
}

func (m *metaStore) GetLatestVersion() (uint64, error) {
	return m.ext.AspectVersion, nil
}

func (m *metaStore) getProperties(version uint64) (properties map[string][]byte, err error) {
	if _, ok := m.propertiesCache[version]; ok {
		return m.propertiesCache[version], nil
	}

	// key format {5B propertyPrefix}{8B version}{20B aspectID}
	key := store.NewKeyBuilder(V1AspectPropertiesKeyPrefix).
		AppendUint64(version).
		AppendBytes(m.ctx.AspectID.Bytes()).
		Build()
	allProps, err := m.Load(key)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err == nil {
			m.propertiesCache[version] = properties
		}
	}()

	if len(allProps) == 0 {
		return make(map[string][]byte), nil
	}

	propertiesArray := make([]types.Property, 0)
	if err := json.Unmarshal(allProps, &propertiesArray); err != nil {
		return nil, err
	}

	properties = make(map[string][]byte, len(propertiesArray))
	for _, property := range propertiesArray {
		properties[property.Key] = property.Value
	}

	return properties, nil
}

func (m *metaStore) GetProperty(version uint64, propKey string) ([]byte, error) {
	if version == 0 {
		// for non-exist version property, return nil
		return nil, nil
	}

	allProps, err := m.getProperties(version)
	if err != nil {
		return nil, err
	}

	return allProps[propKey], nil
}

func (m *metaStore) BumpVersion() (ver uint64, err error) {
	key := store.NewKeyBuilder(store.AspectProtocolInfoKeyPrefix).AppendBytes(m.ctx.AspectID.Bytes()).Build()
	raw, err := m.Load(key)
	if err != nil {
		return 0, err
	}

	m.ext.AspectVersion += 1
	defer func() {
		if err != nil {
			// if failed, rollback the version
			m.ext.AspectVersion -= 1
		}
	}()

	marshaled, err := m.ext.MarshalText()
	if err != nil {
		return 0, err
	}

	extensionOffset := store.ProtocolInfoLen + store.ProtocolVersionLen
	result := make([]byte, 0, len(raw))
	result = append(result, raw[:extensionOffset]...)
	result = append(result, marshaled...)

	return m.ext.AspectVersion, m.Store(key, result)
}

func (m *metaStore) StoreVersionMeta(version uint64, meta *types.VersionMeta) error {
	// key format {5B metaPrefix}{8B version}{20B aspectID}
	key := store.NewKeyBuilder(V1AspectMetaKeyPrefix).
		AppendUint64(version).
		AppendBytes(m.ctx.AspectID.Bytes()).
		Build()

	marshaled, err := VersionMeta(*meta).MarshalText()
	if err != nil {
		return err
	}

	return m.Store(key, marshaled)
}

func (m *metaStore) StoreMeta(meta *types.AspectMeta) (err error) {
	oldPayMaster := m.ext.PayMaster
	oldProof := m.ext.Proof

	m.ext.PayMaster = meta.PayMaster
	m.ext.Proof = meta.Proof

	defer func() {
		// rollback if failed
		if err != nil {
			m.ext.PayMaster = oldPayMaster
			m.ext.Proof = oldProof
		}
	}()

	marshaled, err := m.ext.MarshalText()
	if err != nil {
		return err
	}

	key := store.NewKeyBuilder(store.AspectProtocolInfoKeyPrefix).AppendBytes(m.ctx.AspectID.Bytes()).Build()
	raw, err := m.Load(key)
	if err != nil {
		return err
	}

	extensionOffset := store.ProtocolInfoLen + store.ProtocolVersionLen
	result := make([]byte, 0, len(raw))
	result = append(result, raw[:extensionOffset]...)
	result = append(result, marshaled...)

	return m.Store(key, result)
}

func (m *metaStore) StoreCode(version uint64, code []byte) error {
	// key format {5B codePrefix}{8B version}{20B aspectID}
	key := store.NewKeyBuilder(V1AspectCodeKeyPrefix).
		AppendUint64(version).
		AppendBytes(m.ctx.AspectID.Bytes()).
		Build()
	return m.Store(key, code)
}

func (m *metaStore) StoreProperties(version uint64, properties []types.Property) (err error) {
	// for version > 1 (upgrade case), we need to load all previous properties and append the new ones
	if version > 1 {
		oldVersionProperties, err := m.getProperties(version - 1)
		if err != nil {
			return err
		}

		keySet := make(map[string]int, len(properties))
		for i, prop := range properties {
			keySet[prop.Key] = i
		}

		for key, value := range oldVersionProperties {
			if _, ok := keySet[key]; ok {
				// for exiting one, ignore the old value
				continue
			}

			// for non-existing one, we need to append it
			properties = append(properties, types.Property{
				Key:   key,
				Value: value,
			})
		}
	}

	// check limits
	if len(properties) > types.AspectPropertyLimit {
		return store.ErrTooManyProperties
	}

	// sort the slice lexicographically
	sort.Slice(properties, func(i, j int) bool {
		return properties[i].Key < properties[j].Key
	})

	// key format {5B propertyPrefix}{8B version}{20B aspectID}
	key := store.NewKeyBuilder(V1AspectPropertiesKeyPrefix).
		AppendUint64(version).
		AppendBytes(m.ctx.AspectID.Bytes()).
		Build()

	bytes, err := json.Marshal(properties)
	if err != nil {
		return store.ErrSerdeFail
	}

	defer func() {
		if err == nil {
			m.propertiesCache[version] = make(map[string][]byte, len(properties))
			for _, prop := range properties {
				m.propertiesCache[version][prop.Key] = prop.Value
			}
		}
	}()

	return m.Store(key, bytes)
}

func (m *metaStore) StoreBinding(account common.Address, version uint64, joinPoint uint64, priority int8) (err error) {
	// bindingKey format {5B codePrefix}{8B version}{20B aspectID}
	bindingKey := store.NewKeyBuilder(V1AspectBindingKeyPrefix).
		AppendBytes(m.ctx.AspectID.Bytes())

	// load first slot
	bindingSlotKey := bindingKey.AppendByte(V1AspectBindingDataKeyPrefix)
	firstSlotKey := bindingSlotKey.AppendUint64(0).Build()
	firstSlot, err := m.Load(firstSlotKey)
	if err != nil {
		return err
	}

	// marshal binding info
	newBinding := Binding{
		Account:   account,
		Version:   version,
		Priority:  priority,
		JoinPoint: uint16(joinPoint),
	}
	newBindingBytes, err := newBinding.MarshalText()
	if err != nil {
		return err
	}

	var length DataLength
	if err := length.UnmarshalText(firstSlot); err != nil {
		return err
	}

	// check max binding amount, an aspect can bind with at most 459000 accounts
	if length >= maxBindingSize {
		return store.ErrBindingLimitExceeded
	}

	filterKey := bindingKey.AppendByte(V1AspectBindingFilterKeyPrefix)
	if length == 0 {
		// first time this aspect is bound
		// use cuckoo filter instead of bloom, since cuckoo also support delete,
		// we use a separate slot to store the filters,
		// each filter will manage 28 slots, which should be enough for a long time.

		// first time this aspect is bound
		filter := NewLoggedFilter(m.ctx.Logger(), cuckoo.NewFilter(filterMaxSize))
		// insert {aspectId}:{filterManagedSlotOffset 0-27} into cuckoo filter
		// each filter will manage 28 binding slots
		filter.Insert(store.NewKeyBuilder(account.Bytes()).AppendUint8(0).Build())
		filterData := filter.Encode()
		// save filter
		if err := m.Store(filterKey.AppendUint8(0).Build(), filterData); err != nil {
			return err
		}
		// save binding data
		length := DataLength(1)
		lengthBytes, err := length.MarshalText()
		if err != nil {
			return err
		}

		// store binding data
		if err := m.Store(firstSlotKey, append(lengthBytes, newBindingBytes...)); err != nil {
			return err
		}

		return nil
	}

	lastSlot := uint64(length / bindingSlotSize)
	lastFilterSlot := uint8(lastSlot / filterManagedSlots)

	var lastFilter Filter
	// first use filter to check whether the account is already bound
	for i := uint8(0); i <= lastFilterSlot; i++ {
		key := filterKey.AppendUint8(i).Build()
		filterData, err := m.Load(key)
		if err != nil {
			return err
		}
		if len(filterData) == 0 {
			// reached the end of the filter
			m.ctx.Logger().Debug("filter not found", "key", key)
			break
		}

		cuckooFilter, err := cuckoo.Decode(filterData)
		if err != nil {
			m.ctx.Logger().Error("failed to decode filter", "err", err)
			return err
		}

		filter := NewLoggedFilter(m.ctx.Logger(), cuckooFilter)

		// cache the last filter
		if i == lastFilterSlot {
			lastFilter = filter
		}

		accountKey := store.NewKeyBuilder(account.Bytes())

		// if total slot is less than managed, just test all slots, otherwise test managed
		slotsToTest := uint8(filterManagedSlots)
		if filterManagedSlots > lastSlot {
			slotsToTest = uint8(lastSlot + 1)
		}

		// check aspect is already bound
		for j := uint8(0); j < slotsToTest; j++ {
			if !filter.Lookup(accountKey.AppendUint8(j).Build()) {
				// filter test fail, continue searching
				m.ctx.Logger().Debug("filter test failed", "account", account.Hex(), "slot", j)
				continue
			}

			// filter test succeeded, double check with the actual data
			var bindingDataInSlot []byte
			if i == 0 && j == 0 {
				// since we have already loaded the first slot, no need to load it again
				bindingDataInSlot = firstSlot[bindingDataLength:]
			} else {
				slotIndex := uint64(i)*filterManagedSlots + uint64(j)
				bindingDataInSlot, err = m.Load(bindingSlotKey.AppendUint64(slotIndex).Build())
				if err != nil {
					return err
				}
			}
			for k := 0; k < len(bindingDataInSlot); k += bindingInfoLength {
				data := bindingDataInSlot[k : k+bindingInfoLength]
				var binding Binding
				if err := binding.UnmarshalText(data); err != nil {
					return err
				}
				if binding.Account == account {
					return store.ErrAlreadyBound
				}
			}
		}
	}

	// check if the last slot is full
	if length%bindingSlotSize == 0 {
		// last slot is full, check if we need create a new filter
		if lastSlot%filterManagedSlots == 0 {
			// create a new filter
			filter := NewLoggedFilter(m.ctx.Logger(), cuckoo.NewFilter(filterMaxSize))
			// insert {aspectId}:{filterManagedSlotOffset 0-27} into cuckoo filter
			// each filter will manage 28 binding slots
			filter.Insert(store.NewKeyBuilder(account.Bytes()).AppendUint8(0).Build())
			filterData := filter.Encode()
			// save filter
			if err := m.Store(filterKey.AppendUint8(lastFilterSlot).Build(), filterData); err != nil {
				return err
			}
		} else {
			// update filter data
			lastFilterKey := filterKey.AppendUint8(lastFilterSlot).Build()
			if lastFilter == nil {
				// if last filter is not loaded for some reason, load and decode
				filterData, err := m.Load(lastFilterKey)
				if err != nil {
					return err
				}

				cuckooFilter, err := cuckoo.Decode(filterData)
				if err != nil {
					m.ctx.Logger().Error("failed to decode filter", "err", err)
					return err
				}
				lastFilter = NewLoggedFilter(m.ctx.Logger(), cuckooFilter)
			}
			lastFilter.Insert(store.NewKeyBuilder(account.Bytes()).AppendUint8(uint8(lastSlot % filterManagedSlots)).Build())
			filterData := lastFilter.Encode()
			if err := m.Store(lastFilterKey, filterData); err != nil {
				return err
			}
		}
		// create a new slot
		newSlotKey := bindingSlotKey.AppendUint64(lastSlot).Build()
		if err := m.Store(newSlotKey, newBindingBytes); err != nil {
			return err
		}
	} else {
		// slot is not full, update the slot and filter directly
		// update filter data
		lastFilterKey := filterKey.AppendUint8(lastFilterSlot).Build()
		if lastFilter == nil {
			// if last filter is not loaded for some reason, load and decode
			filterData, err := m.Load(lastFilterKey)
			if err != nil {
				return err
			}

			cuckooFilter, err := cuckoo.Decode(filterData)
			if err != nil {
				return err
			}
			lastFilter = NewLoggedFilter(m.ctx.Logger(), cuckooFilter)
		}
		lastFilter.Insert(store.NewKeyBuilder(account.Bytes()).AppendUint8(uint8(lastSlot % filterManagedSlots)).Build())
		filterData := lastFilter.Encode()
		if err := m.Store(lastFilterKey, filterData); err != nil {
			return err
		}

		// update the slot data
		if lastSlot == 0 {
			// since we have already loaded the first slot, no need to load it again
			// just update the first slot data
			firstSlot = append(firstSlot, newBindingBytes...)
		} else {
			// otherwise we need to load the last slot data
			lastSlotKey := bindingSlotKey.AppendUint64(lastSlot).Build()
			lastSlotData, err := m.Load(lastSlotKey)
			if err != nil {
				return err
			}

			// append the new binding to the last slot and save it
			lastSlotData = append(lastSlotData, newBindingBytes...)
			if err := m.Store(lastSlotKey, lastSlotData); err != nil {
				return err
			}
		}
	}

	// update length
	length++
	lengthBytes, err := length.MarshalText()
	if err != nil {
		return err
	}

	// overwrite the length in first slot
	copy(firstSlot, lengthBytes)
	// update first slot
	return m.Store(firstSlotKey, firstSlot)
}

func (m *metaStore) LoadAspectBoundAccounts() ([]types.Binding, error) {
	// key format {5B codePrefix}{8B version}{20B aspectID}
	key := store.NewKeyBuilder(V1AspectBindingKeyPrefix).
		AppendBytes(m.ctx.AspectID.Bytes()).AppendByte(V1AspectBindingDataKeyPrefix)

	firstSlot, err := m.Load(key.AppendUint64(0).Build())
	if err != nil {
		return nil, err
	}

	var length DataLength
	if err := length.UnmarshalText(firstSlot); err != nil {
		return nil, err
	}

	// binding format for first slot {8B Length}{256B Bloom}{32B Binding}{32B Binding}...
	// each slot will save maximum 120 binding info, which is 3840 bytes
	bindingData := firstSlot[bindingDataLength:]
	bindings := make([]types.Binding, 0, length)
	for i := uint64(0); i < uint64(length); i += bindingSlotSize {
		for j := 0; j < len(bindingData); j += bindingInfoLength {
			data := bindingData[j : j+bindingInfoLength]
			var binding Binding
			if err := binding.UnmarshalText(data); err != nil {
				return nil, err
			}
			bindings = append(bindings, types.Binding(binding))
			if uint64(len(bindings)) == uint64(length) {
				// EOF
				return bindings, nil
			}
		}

		// load next slot
		bindingData, err = m.Load(key.AppendUint64((i / bindingSlotSize) + 1).Build())
		if err != nil {
			return nil, err
		}
	}

	return bindings, nil
}

func (m *metaStore) RemoveBinding(account common.Address) (err error) {
	// bindingKey format {5B codePrefix}{8B version}{20B aspectID}
	bindingKey := store.NewKeyBuilder(V1AspectBindingKeyPrefix).
		AppendBytes(m.ctx.AspectID.Bytes())

	// load first slot
	bindingSlotKey := bindingKey.AppendByte(V1AspectBindingDataKeyPrefix)
	firstSlotKey := bindingSlotKey.AppendUint64(0).Build()
	firstSlot, err := m.Load(firstSlotKey)
	if err != nil {
		return err
	}

	var length DataLength
	if err := length.UnmarshalText(firstSlot); err != nil {
		return err
	}

	if length == 0 {
		// if not bound, just pass
		return nil
	}

	lastSlot := uint64(length / bindingSlotSize)
	if length%bindingSlotSize == 0 {
		lastSlot = lastSlot - 1
	}
	lastFilterSlot := uint8(lastSlot / filterManagedSlots)

	var (
		lastFilter        *cuckoo.Filter
		bindingFilter     *cuckoo.Filter
		bindingSlot       *uint64
		bindingSlotData   []byte
		bindingSlotOffset *int
	)
	// first use filter to check whether the account is already bound
	filterKey := bindingKey.AppendByte(V1AspectBindingFilterKeyPrefix)
	for i := uint8(0); i <= lastFilterSlot; i++ {
		key := filterKey.AppendUint8(i).Build()
		filterData, err := m.Load(key)
		if err != nil {
			return err
		}
		if len(filterData) == 0 {
			// reached last filter
			break
		}

		filter, err := cuckoo.Decode(filterData)
		if err != nil {
			return err
		}

		// cache the last filter
		if i == lastFilterSlot {
			lastFilter = filter
		}

		accountKey := store.NewKeyBuilder(account.Bytes())
		for j := uint8(0); j < filterManagedSlots; j++ {
			if !filter.Lookup(accountKey.AppendUint8(j).Build()) {
				// filter test fail, continue searching
				continue
			}

			// filter test succeeded, double check with the actual data
			var bindingDataInSlot []byte
			slotIndex := uint64(0)
			if i == 0 && j == 0 {
				// since we have already loaded the first slot, no need to load it again
				bindingDataInSlot = firstSlot[bindingDataLength:]
			} else {
				slotIndex = uint64(i)*filterManagedSlots + uint64(j)
				bindingDataInSlot, err = m.Load(bindingSlotKey.AppendUint64(slotIndex).Build())
				if err != nil {
					return err
				}
			}
			// search in the slot
			for k := 0; k < len(bindingDataInSlot); k += bindingInfoLength {
				data := bindingDataInSlot[k : k+bindingInfoLength]
				var binding Binding
				if err := binding.UnmarshalText(data); err != nil {
					return err
				}
				if binding.Account == account {
					// found binding, record position and break out the loops
					bindingSlotData = bindingDataInSlot
					bindingSlotOffset = &k
					bindingSlot = &slotIndex
					bindingFilter = filter
					goto BindingFound
				}
			}
		}
	}

BindingFound:
	if bindingSlotData == nil || bindingSlot == nil || bindingFilter == nil || bindingSlotOffset == nil {
		// if not bound, just pass
		return nil
	}

	// remove the binding from the slot
	if *bindingSlot == lastSlot {
		// only first slot used, or binding in last slot, we can just remove the binding from the slot
		if len(bindingSlotData) > bindingInfoLength {
			copy(bindingSlotData[*bindingSlotOffset:], bindingSlotData[*bindingSlotOffset+bindingInfoLength:])
		}
		bindingSlotData = bindingSlotData[:len(bindingSlotData)-bindingInfoLength]
		if *bindingSlot > 0 {
			// if not the first slot we just save the data
			if err := m.Store(bindingSlotKey.AppendUint64(*bindingSlot).Build(), bindingSlotData); err != nil {
				return err
			}
		} else {
			// for the first slot we just need update the first slot data
			firstSlot = append(firstSlot[:bindingDataLength], bindingSlotData...)
		}

		// remove the account from the filter, and update filter
		filterSlot := uint8(*bindingSlot / filterManagedSlots)
		bindingFilter.Delete(store.NewKeyBuilder(account.Bytes()).AppendUint8(filterSlot).Build())
		// if there is nothing in the filter, delete it
		var updatedFilter []byte
		if bindingFilter.Count() > 0 {
			// otherwise we update it
			updatedFilter = bindingFilter.Encode()
		}

		if err := m.Store(filterKey.AppendUint8(filterSlot).Build(), updatedFilter); err != nil {
			return err
		}
	} else if lastSlot > *bindingSlot {
		// move the last binding to the removed position
		lastSlotKey := bindingSlotKey.AppendUint64(lastSlot).Build()
		lastSlotData, err := m.Load(lastSlotKey)
		if err != nil {
			return err
		}
		lastBindingBytes := lastSlotData[len(lastSlotData)-bindingInfoLength:]

		// replace delete binding with the last binding
		copy(bindingSlotData[*bindingSlotOffset:], lastBindingBytes)
		lastSlotData = lastSlotData[:len(lastSlotData)-bindingInfoLength]

		// update last slot
		if err := m.Store(lastSlotKey, lastSlotData); err != nil {
			return err
		}

		// update the updated slot
		if *bindingSlot == 0 {
			// for the first slot we just need to copy over the data
			copy(firstSlot[bindingDataLength:], bindingSlotData)
		} else {
			// for other slots we need to save it first
			if err := m.Store(bindingSlotKey.AppendUint64(*bindingSlot).Build(), bindingSlotData); err != nil {
				return err
			}
		}

		// load the filter if needed
		if lastFilter == nil {
			// if last filter is not loaded for some reason, load and decode
			filterData, err := m.Load(filterKey.AppendUint8(lastFilterSlot).Build())
			if err != nil {
				return err
			}
			lastFilter, err = cuckoo.Decode(filterData)
			if err != nil {
				return err
			}
		}

		lastBinding := new(Binding)
		if err := lastBinding.UnmarshalText(lastBindingBytes); err != nil {
			return err
		}

		// remove the binding from the filters
		if lastSlot/filterManagedSlots == *bindingSlot/filterManagedSlots {
			bindingFilterOffsetKey := uint8(*bindingSlot % filterManagedSlots)
			lastFilterOffsetKey := uint8(lastSlot % filterManagedSlots)
			// need to remove both last binding and the unbound one from their old position
			// and add last one to the new position
			lastFilter.Delete(store.NewKeyBuilder(account.Bytes()).AppendUint8(bindingFilterOffsetKey).Build())
			lastFilter.Delete(store.NewKeyBuilder(lastBinding.Account.Bytes()).AppendUint8(lastFilterOffsetKey).Build())
			lastFilter.Insert(store.NewKeyBuilder(lastBinding.Account.Bytes()).AppendUint8(bindingFilterOffsetKey).Build())

			// update last filter
			if err := m.Store(filterKey.AppendUint8(lastFilterSlot).Build(), lastFilter.Encode()); err != nil {
				return err
			}
		} else {
			// if the last slot and the binding slot are in different filters, we need to update both filters
			// remove the moved binding account from the last slot filter, and update filter
			lastFilterOffsetKey := uint8(lastSlot % filterManagedSlots)
			lastFilter.Delete(store.NewKeyBuilder(lastBinding.Account.Bytes()).AppendUint8(lastFilterOffsetKey).Build())
			var updatedFilter []byte
			if lastFilter.Count() > 0 {
				updatedFilter = lastFilter.Encode()
			}
			if err := m.Store(filterKey.AppendUint8(lastFilterSlot).Build(), updatedFilter); err != nil {
				return err
			}

			// remove the account from the binding slot filter, and update filter
			bindingFilterOffset := uint8(*bindingSlot % filterManagedSlots)
			bindingFilter.Delete(store.NewKeyBuilder(account.Bytes()).AppendUint8(bindingFilterOffset).Build())
			bindingFilter.Insert(store.NewKeyBuilder(lastBinding.Account.Bytes()).AppendUint8(bindingFilterOffset).Build())
			if err := m.Store(filterKey.AppendUint8(uint8(*bindingSlot/filterManagedSlots)).Build(), bindingFilter.Encode()); err != nil {
				return err
			}
		}
	} else {
		// should not happen
		return store.ErrStorageCorrupted
	}

	// update length
	length--
	lengthBytes, err := length.MarshalText()
	if err != nil {
		return err
	}

	// overwrite the length in first slot
	copy(firstSlot, lengthBytes)
	// update first slot
	return m.Store(firstSlotKey, firstSlot)
}

func (m *metaStore) MigrateFrom(old store.AspectMetaStore) error {
	//TODO implement me
	panic("implement me")
}

func (m *metaStore) Used() (bool, error) {
	return m.ext.AspectVersion > 0, nil
}

func (m *metaStore) Init() error {
	versionBytes, err := protocolVersion.MarshalText()
	if err != nil {
		return err
	}

	info := &store.AspectInfo{
		MetaVersion:  protocolVersion,
		StateVersion: protocolVersion,
	}

	infoBytes, err := info.MarshalText()
	if err != nil {
		return err
	}

	extension := &Extension{
		AspectVersion: 0,
		PayMaster:     common.Address{},
		Proof:         nil,
	}
	extBytes, err := extension.MarshalText()
	if err != nil {
		return err
	}

	key := store.NewKeyBuilder(store.AspectProtocolInfoKeyPrefix).AppendBytes(m.ctx.AspectID.Bytes()).Build()
	result := make([]byte, 0, len(versionBytes)+len(infoBytes)+len(extBytes))
	result = append(result, versionBytes...)
	result = append(result, infoBytes...)
	result = append(result, extBytes...)

	return m.Store(key, result)
}
