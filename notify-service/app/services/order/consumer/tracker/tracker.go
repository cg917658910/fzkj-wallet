package tracker

import (
	"sync"
	"time"

	"github.com/IBM/sarama"
)

// 单个 offset 的状态结构体
type OffsetStatus struct {
	State     string    // "pending" or "done"
	Timestamp time.Time // 创建或完成时间
}

// OffsetState 表示单个分区的 offset 跟踪器
type OffsetState struct {
	mu       sync.Mutex
	status   map[int64]*OffsetStatus
	minReady int64
	ttl      time.Duration
}

type (
	topicName   string
	partitionID int32
)
type KafkaSafeConsumer struct {
	trackers   map[topicName]map[partitionID]*OffsetState
	trackersMu sync.Mutex
	topic      string
}

func NewKafkaSafeConsumer() *KafkaSafeConsumer {
	return &KafkaSafeConsumer{
		trackers:   map[topicName]map[partitionID]*OffsetState{},
		trackersMu: sync.Mutex{},
	}
}

func (c *KafkaSafeConsumer) GetTracker(topic topicName, partition partitionID) *OffsetState {
	c.trackersMu.Lock()
	defer c.trackersMu.Unlock()
	/* trackerTopic, ok := c.trackers[topic]
	if !ok {
		trackerTopic = map[partitionID]*OffsetState{}
		c.trackers[partition] = tracker
	} */
	return nil
}

func (c *KafkaSafeConsumer) SafeCommit(session sarama.ConsumerGroupSession) {
	if session == nil {
		return
	}
	c.trackersMu.Lock()
	defer c.trackersMu.Unlock()
	// mark offset
	/* for _, tracker := range c.trackers{

	} */
}

// 构造一个新的 OffsetState
func NewOffsetState() *OffsetState {
	return &OffsetState{
		status:   make(map[int64]*OffsetStatus),
		minReady: -1,
		ttl:      10 * time.Minute,
	}
}

// Init 标记 offset
func (s *OffsetState) SetInitOffset(offset int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status[offset] = &OffsetStatus{
		State:     "pending",
		Timestamp: time.Now(),
	}
}

// Start 标记 offset 开始处理
func (s *OffsetState) Start(offset int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status[offset] = &OffsetStatus{
		State:     "pending",
		Timestamp: time.Now(),
	}
}

// Done 标记 offset 完成
func (s *OffsetState) Done(offset int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status[offset] = &OffsetStatus{
		State:     "done",
		Timestamp: time.Now(),
	}
}

// CommitOffset 找出最大连续 done 的 offset，并清理过期状态
func (s *OffsetState) CommitOffset() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := s.minReady + 1
	for {
		entry, ok := s.status[next]
		if !ok || entry.State != "done" {
			break
		}
		delete(s.status, next)
		s.minReady = next
		next++
	}
	return s.minReady + 1
}
