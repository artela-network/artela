package pubsub

import (
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
)

type UnsubscribeFunc func()

type EventBus interface {
	AddTopic(name string, src <-chan coretypes.ResultEvent) error
	RemoveTopic(name string)
	Subscribe(name string) (<-chan coretypes.ResultEvent, UnsubscribeFunc, error)
	Topics() []string
}

type subChan struct {
	ch    chan<- coretypes.ResultEvent
	unsub func()
}

type memEventBus struct {
	topics          map[string]<-chan coretypes.ResultEvent
	topicsMux       *sync.RWMutex
	subscribers     map[string]map[uint64]subChan
	subscribersMux  *sync.RWMutex
	currentUniqueID uint64
}

func NewEventBus() EventBus {
	return &memEventBus{
		topics:         make(map[string]<-chan coretypes.ResultEvent),
		topicsMux:      new(sync.RWMutex),
		subscribers:    make(map[string]map[uint64]subChan),
		subscribersMux: new(sync.RWMutex),
	}
}

func (m *memEventBus) GenUniqueID() uint64 {
	return atomic.AddUint64(&m.currentUniqueID, 1)
}

func (m *memEventBus) Topics() (topics []string) {
	m.topicsMux.RLock()
	defer m.topicsMux.RUnlock()

	topics = make([]string, 0, len(m.topics))
	for topicName := range m.topics {
		topics = append(topics, topicName)
	}

	return topics
}

func (m *memEventBus) AddTopic(name string, src <-chan coretypes.ResultEvent) error {
	m.topicsMux.Lock()
	defer m.topicsMux.Unlock()

	if _, ok := m.topics[name]; ok {
		return errors.New("topic already registered")
	}

	m.topics[name] = src

	go m.publishTopic(name, src)

	return nil
}

func (m *memEventBus) RemoveTopic(name string) {
	m.topicsMux.Lock()
	defer m.topicsMux.Unlock()
	delete(m.topics, name)
}

func (m *memEventBus) Subscribe(name string) (<-chan coretypes.ResultEvent, UnsubscribeFunc, error) {
	m.topicsMux.RLock()
	_, ok := m.topics[name]
	m.topicsMux.RUnlock()

	if !ok {
		return nil, nil, errors.Errorf("topic not found: %s", name)
	}

	var closeOnce sync.Once
	ch := make(chan coretypes.ResultEvent)
	m.subscribersMux.Lock()
	defer m.subscribersMux.Unlock()

	id := m.GenUniqueID()
	if _, ok := m.subscribers[name]; !ok {
		m.subscribers[name] = make(map[uint64]subChan)
	}

	unsubscribe := func() {
		m.subscribersMux.Lock()
		defer m.subscribersMux.Unlock()
		if subs, ok := m.subscribers[name]; ok {
			delete(subs, id)
			if len(subs) == 0 {
				delete(m.subscribers, name)
			}
		}
		closeOnce.Do(func() { close(ch) })
	}
	m.subscribers[name][id] = subChan{ch: ch, unsub: unsubscribe}

	return ch, unsubscribe, nil
}

func (m *memEventBus) publishTopic(name string, src <-chan coretypes.ResultEvent) {
	for {
		msg, ok := <-src
		if !ok {
			m.closeAllSubscribers(name)
			m.topicsMux.Lock()
			delete(m.topics, name)
			m.topicsMux.Unlock()
			return
		}
		m.publishAllSubscribers(name, msg)
	}
}

func (m *memEventBus) closeAllSubscribers(name string) {
	m.subscribersMux.Lock()
	defer m.subscribersMux.Unlock()

	if subscribers, ok := m.subscribers[name]; ok {
		delete(m.subscribers, name)
		for _, sub := range subscribers {
			sub.unsub()
		}
	}
}

func (m *memEventBus) publishAllSubscribers(name string, msg coretypes.ResultEvent) {
	m.subscribersMux.RLock()
	defer m.subscribersMux.RUnlock()

	if subscribers, ok := m.subscribers[name]; ok {
		for _, sub := range subscribers {
			select {
			case sub.ch <- msg:
			default:
				// Handle message drop if necessary
			}
		}
	}
}
